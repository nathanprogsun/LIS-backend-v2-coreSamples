package service

import (
	"context"
	"coresamples/common"
	"coresamples/config"
	"coresamples/dbutils"
	"coresamples/ent"
	"coresamples/ent/serviceship"
	"coresamples/ent/serviceshipbillingplan"
	"coresamples/external"
	pb "coresamples/proto"
	"coresamples/publisher"
	_ "embed"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

var WatcherRecoverCount = 0

const (
	MaxRecoverCount = 10
)

type IServiceshipService interface {
	SubscriptionAllowed(clinicName string, ctx context.Context) bool
	CreateServiceship(tag string, svcType string, fees []float32, billingCycles []int32, intervals []string, ctx context.Context) error
	GetServiceshipConfig(svcshipID int, ctx context.Context) (interface{}, serviceship.Type, error)
	Subscribe(accountID int64, accountType string, subscriberName string, email string, billingPlanID int, keepActive bool, platform string, paymentToken string, customerToken string, paymentMethodID int64, ctx context.Context) error
	GetAccountSubscriptions(accountID int64, accountType string, includeOutdated bool, ctx context.Context) ([]*ent.AccountSubscription, error)
	GetChargingSubscription(subscriptionID int, ctx context.Context) (*external.SubscriptionInfo, error)
	//AddPermissionsToServiceship(name string, permissions map[string][]string, callerAccountID int32) error
	GetServiceships(ctx context.Context) ([]*ent.Serviceship, error)
	GetServiceshipsByType(svcType string, ctx context.Context) ([]*ent.Serviceship, error)
	AddBillingPlan(svcshipID int, fee float32, billingCycle int32, interval string, ctx context.Context) error
	CreateBillingPlanSet(svcshipID int, fee []float32, billingCycles []int32, interval []string, effective int64, ctx context.Context) error
	GetLatestBillingPlanSet(svcshipID int, ctx context.Context) ([]*ent.ServiceshipBillingPlan, error)
	GetBillingPlanByID(billingPlanID int, ctx context.Context) (*ent.ServiceshipBillingPlan, error)
	GetBillingPlanBySubscriptionID(subscriptionID int, ctx context.Context) (*ent.ServiceshipBillingPlan, error)
	GetPaymentMethods(accountID int64, accountType string, ctx context.Context) ([]*external.PaymentMethod, error)
	CreatePaymentMethod(accountID int64, accountType string, platform string, paymentToken string, customerToken string, cardType string, expiryDate string, lastFour string, ctx context.Context) (int64, error)
	UpdateSubscriptionPaymentMethod(subscriptionID int, paymentID int64, ctx context.Context) error
	DeletePaymentMethod(paymentMethodID int64, accountID int64, accountType string, ctx context.Context) error
	PauseAutoRenew(subscriptionID int, ctx context.Context) error
	ResumeAutoRenew(subscriptionID int, platform string, paymentToken string, customerToken string, paymentMethodID int64, ctx context.Context) error
	UpdateSubscriptionBillingPlan(subscriptionID int, billingPlanID int, ctx context.Context) error
	UpdateSubscriptionEmail(email string, subscriptionID int, ctx context.Context) error
	UpdateAllSubscriptionEmail(email string, accountID int64, accountType string, ctx context.Context) error
	GetSubscriptionTransactionInfos(subscriptionID int, ctx context.Context) (*external.Transactions, error)
}

type ServiceshipService struct {
	Service
	rbacService    IRBACService
	allowedClinics []string
}

func newServiceshipService(dbClient *ent.Client, redisClient *common.RedisClient, allowedClinics []string) IServiceshipService {
	s := &ServiceshipService{
		Service:        InitService(dbClient, redisClient),
		allowedClinics: allowedClinics,
	}
	ctx := context.Background()
	go s.watchSubscriptions(ctx)
	return s
}

func (s *ServiceshipService) CreateServiceship(tag string, svcType string, fees []float32, billingCycles []int32, intervals []string, ctx context.Context) error {
	switch svcType {
	case "membership":
		if _, ok := config.GetServiceshipConfig().Membership[tag]; !ok {
			return errors.New("no corresponding service found")
		}
		return dbutils.CreateServiceship(tag, svcType, fees, billingCycles, intervals, s.dbClient, ctx)
	default:
		return errors.New("no corresponding service type found")
	}
}

func (s *ServiceshipService) GetServiceshipConfig(svcshipID int, ctx context.Context) (interface{}, serviceship.Type, error) {
	svc, err := dbutils.GetServiceshipByID(svcshipID, s.dbClient, ctx)
	if err != nil {
		return nil, "", err
	}
	switch svc.Type {
	case serviceship.TypeMembership:
		return config.GetServiceshipConfig().Membership[svc.Tag], svc.Type, nil
	}
	return nil, "", errors.New("wrong type")
}

func (s *ServiceshipService) SubscriptionAllowed(clinicName string, ctx context.Context) bool {
	for _, name := range s.allowedClinics {
		if name == clinicName {
			return true
		}
	}
	return false
}

func (s *ServiceshipService) Subscribe(accountID int64, accountType string, subscriberName string, email string, billingPlanID int, keepActive bool, platform string, paymentToken string, customerToken string, paymentMethodID int64, ctx context.Context) error {
	plan, err := dbutils.GetBillingPlanByID(billingPlanID, s.dbClient, ctx)
	if err != nil {
		return err
	}
	// create subscription
	sub, err := dbutils.CreateSubscription(accountID, accountType, subscriberName, email, billingPlanID, keepActive, s.dbClient, ctx)
	if err != nil {
		return err
	}
	var paymentID int64
	// make one time payment and create subscription if not keep active
	if !keepActive {
		success := false
		fee, _ := decimal.NewFromFloat(float64(plan.Fee)).Round(2).Float64()
		success, paymentID, err = external.GetChargingService().OneTimeCharge(accountType, accountID, float32(fee), sub.ID, platform, paymentToken, customerToken)
		if !success {
			dbutils.DeleteSubscription(sub.ID, s.dbClient, ctx)
			return fmt.Errorf("payment failed, clinic %d, platform: %s, payment token: %s, customer token: %s, err: %v",
				accountID, platform, paymentToken, customerToken, err)
		}
		if err != nil {
			common.Errorf("error occurred during subscription payment", err)
		}

		err = external.GetChargingService().CreateSubscription(accountType, accountID, float32(fee), sub.ID, getSubscriptionEndTime(plan, sub.StartTime), plan.BillingCycle, string(plan.Interval), paymentMethodID, 1)
	}

	confirmTemplate := &pb.SubscriptionConfirmationEmailTemplate{
		ProviderName:   subscriberName,
		AccountName:    subscriberName,
		InvoiceNumber:  strconv.FormatInt(paymentID, 10),
		PurchasedDate:  time.Now().Format("01/02/2006"),
		ActiveTillDate: sub.EndTime.Format("01/02/2006"),
		RenewalDate:    sub.EndTime.AddDate(0, 0, -7).Format("01/02/2006"),
		Charge:         "$" + strconv.Itoa(int(plan.Fee)),
		Subtotal:       "$" + strconv.Itoa(int(plan.Fee)),
		Payable:        strconv.Itoa(int(plan.Fee)),
		Payment:        "$" + strconv.Itoa(int(plan.Fee)),
	}

	switch plan.BillingCycle {
	case 1:
		confirmTemplate.IsMonthly = "true"
	case 6:
		confirmTemplate.IsHalfYear = "true"
	case 12:
		confirmTemplate.IsAnnual = "true"
	}

	err = publisher.GetPublisher().SendSubscriptionSuccessEmail(confirmTemplate, sub.Email, strconv.Itoa(sub.ID))
	if err != nil {
		common.Error(err)
		return nil
	}
	return err
}

func (s *ServiceshipService) GetAccountSubscriptions(accountID int64, accountType string, includeOutdated bool, ctx context.Context) ([]*ent.AccountSubscription, error) {
	// by clinic id, regardless of expiring or not
	return dbutils.GetAccountSubscriptions(accountID, accountType, includeOutdated, s.dbClient, ctx)
}

func (s *ServiceshipService) GetChargingSubscription(subscriptionID int, ctx context.Context) (*external.SubscriptionInfo, error) {
	accountID, accountType := s.getSubscriptionAccountInfo(subscriptionID, ctx)
	return external.GetChargingService().GetSubscription(subscriptionID, accountType, accountID)
}

func (s *ServiceshipService) GetServiceships(ctx context.Context) ([]*ent.Serviceship, error) {
	return dbutils.GetServiceships(s.dbClient, ctx)
}

func (s *ServiceshipService) GetServiceshipsByType(svcType string, ctx context.Context) ([]*ent.Serviceship, error) {
	return dbutils.GetServiceshipsByType(svcType, s.dbClient, ctx)
}

func (s *ServiceshipService) AddBillingPlan(svcshipID int, fee float32, billingCycle int32, interval string, ctx context.Context) error {
	return dbutils.AddBillingPlan(svcshipID, fee, billingCycle, interval, s.dbClient, ctx)
}

func (s *ServiceshipService) CreateBillingPlanSet(svcshipID int, fee []float32, billingCycle []int32, interval []string, effective int64, ctx context.Context) error {
	effectiveTime := time.Unix(effective, 0)
	subs, err := dbutils.MaybeMatchNewestBillingPlans(billingCycle, fee, interval, effectiveTime, svcshipID, s.dbClient, ctx)
	if err != nil {
		return err
	}
	for _, sub := range subs {
		plan, err := dbutils.GetBillingPlanBySubscriptionID(sub.ID, s.dbClient, ctx)
		if err != nil {
			common.Error(err)
			continue
		}
		err = external.GetChargingService().UpdateSubscriptionFee(sub.ID, string(sub.AccountType), sub.AccountID, plan.Fee)
		if err != nil {
			common.Error(err)
		}
	}
	return dbutils.CreateBillingPlanSet(svcshipID, fee, billingCycle, interval, effectiveTime, s.dbClient, ctx)
}

func (s *ServiceshipService) GetLatestBillingPlanSet(svcshipID int, ctx context.Context) ([]*ent.ServiceshipBillingPlan, error) {
	return dbutils.GetLatestEffectiveBillingPlanSet(svcshipID, s.dbClient, ctx)
}

func (s *ServiceshipService) GetBillingPlanByID(billingPlanID int, ctx context.Context) (*ent.ServiceshipBillingPlan, error) {
	return dbutils.GetBillingPlanByID(billingPlanID, s.dbClient, ctx)
}

func (s *ServiceshipService) GetBillingPlanBySubscriptionID(subscriptionID int, ctx context.Context) (*ent.ServiceshipBillingPlan, error) {
	return dbutils.GetBillingPlanBySubscriptionID(subscriptionID, s.dbClient, ctx)
}

func (s *ServiceshipService) GetPaymentMethods(accountID int64, accountType string, ctx context.Context) ([]*external.PaymentMethod, error) {
	return external.GetChargingService().GetPaymentMethods(accountType, accountID)
}

func (s *ServiceshipService) CreatePaymentMethod(accountID int64, accountType string, platform string, paymentToken string, customerToken string, cardType string, expiryDate string, lastFour string, ctx context.Context) (int64, error) {
	method := &external.PaymentMethod{
		AccountID:     accountID,
		AccountType:   accountType,
		Type:          "subscription",
		TokenPlatform: platform,
		PaymentToken:  paymentToken,
		CardType:      cardType,
		ExpiryDate:    expiryDate,
		LastFour:      lastFour,
		CustomerToken: customerToken,
		Subscription:  true,
	}
	return external.GetChargingService().CreatePaymentMethod(method)
}

func (s *ServiceshipService) UpdateSubscriptionPaymentMethod(subscriptionID int, paymentID int64, ctx context.Context) error {
	accountID, accountType := s.getSubscriptionAccountInfo(subscriptionID, ctx)
	return external.GetChargingService().UpdatePaymentMethod(subscriptionID, paymentID, accountType, accountID)
}

func (s *ServiceshipService) DeletePaymentMethod(paymentMethodID int64, accountID int64, accountType string, ctx context.Context) error {
	// return error if an active subscription is currently using the subscription
	subs, err := dbutils.GetAccountSubscriptions(accountID, accountType, false, s.dbClient, ctx)
	if err != nil {
		return err
	}

	for _, sub := range subs {
		s, err := external.GetChargingService().GetSubscription(sub.ID, accountType, accountID)
		if err != nil {
			common.Error(err)
		}

		if s.PaymentMethod.ID == paymentMethodID {
			return fmt.Errorf("payment method already in use, subscription id: %d", sub.ID)
		}
	}

	return external.GetChargingService().DeletePaymentMethod(paymentMethodID, accountType, accountID)
}

func (s *ServiceshipService) PauseAutoRenew(subscriptionID int, ctx context.Context) error {
	// set auto-renew to false
	//err := dbutils.SetAutoRenew(subscriptionID, false, s.dbClient, ctx)
	//if err != nil {
	//	return err
	//}
	// update subscription status
	accountID, accountType := s.getSubscriptionAccountInfo(subscriptionID, ctx)
	return external.GetChargingService().PauseSubscription(subscriptionID, accountType, accountID)
}

func (s *ServiceshipService) ResumeAutoRenew(subscriptionID int, platform string, paymentToken string, customerToken string, paymentMethodID int64, ctx context.Context) error {
	// if end date is within a week, will make a one time charge to extend the subscription first
	sub, err := dbutils.GetSubscriptionByID(subscriptionID, s.dbClient, ctx)
	if err != nil {
		return err
	}
	plan, err := sub.Edges.ServiceshipBillingPlanOrErr()
	if err != nil {
		return err
	}
	now := time.Now()
	if sub.EndTime.Before(now) {
		return fmt.Errorf("subscription expired at %s, create new subscription instead", sub.EndTime.String())
	}
	accountID, accountType := s.getSubscriptionAccountInfo(subscriptionID, ctx)
	if sub.EndTime.After(now.AddDate(0, 0, -7)) {
		fee, _ := decimal.NewFromFloat(float64(plan.Fee)).Round(2).Float64()
		success, _, err := external.GetChargingService().OneTimeCharge(accountType, accountID, float32(fee), sub.ID, platform, paymentToken, customerToken)
		if err != nil {
			common.Error(err)
		}
		if !success {
			return fmt.Errorf("payment failed for resuming subscription")
		}
		err = dbutils.AddSubscriptionEndTime(plan.BillingCycle, sub.ID, s.dbClient, ctx)
		if err != nil {
			return err
		}
		return external.GetChargingService().ResumeSubscription(subscriptionID, accountType, accountID, getSubscriptionEndTime(plan, sub.EndTime), paymentMethodID)
	} else {
		return external.GetChargingService().ResumeSubscription(subscriptionID, accountType, accountID, sub.EndTime.AddDate(0, 0, -7), paymentMethodID)
	}
}

// UpdateSubscriptionBillingPlan will take effect in the next billing cycle
func (s *ServiceshipService) UpdateSubscriptionBillingPlan(subscriptionID int, billingPlanID int, ctx context.Context) error {
	sub, err := dbutils.GetSubscriptionByID(subscriptionID, s.dbClient, ctx)
	if err != nil {
		return err
	}
	oldPlan, err := sub.Edges.ServiceshipBillingPlanOrErr()
	if err != nil {
		return err
	}
	plan, err := dbutils.GetBillingPlanByID(billingPlanID, s.dbClient, ctx)
	if err != nil {
		return err
	}
	if sub.EndTime.Before(time.Now()) {
		return errors.New("subscription already expired at " + sub.EndTime.String())
	}
	err = dbutils.UpdateSubscriptionBillingPlan(subscriptionID, billingPlanID, s.dbClient, ctx)
	if err != nil {
		return err
	}
	fee, _ := decimal.NewFromFloat(float64(plan.Fee)).Round(2).Float64()
	err = external.GetChargingService().UpdateSubscriptionBillingCycle(subscriptionID, string(sub.AccountType), sub.AccountID, plan.BillingCycle, float32(fee))
	if err != nil {
		err1 := dbutils.UpdateSubscriptionBillingPlan(subscriptionID, oldPlan.ID, s.dbClient, ctx)
		if err1 == nil {
			return err
		}
		common.Errorf("error resetting billing plan", err1)
		return err
	}
	return nil
}

func (s *ServiceshipService) UpdateSubscriptionEmail(email string, subscriptionID int, ctx context.Context) error {
	return dbutils.UpdateSubscriptionEmail(subscriptionID, email, s.dbClient, ctx)
}

func (s *ServiceshipService) UpdateAllSubscriptionEmail(email string, accountID int64, accountType string, ctx context.Context) error {
	return dbutils.UpdateAllSubscriptionsEmail(accountID, accountType, email, s.dbClient, ctx)
}

func (s *ServiceshipService) watchSubscriptions(ctx context.Context) {
	defer s.recoverSubscriptionWatcher(ctx)
	//sentry.CaptureMessage("Hello from subscription_watcher#" + strconv.Itoa(WatcherRecoverCount))
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1)
	// routine will start at 00:00 every day
	duration := time.Until(start)
	for {
		<-time.Tick(duration)
		//sentry.CaptureMessage("subscription watcher start scanning for outdated subscriptions at " + time.Now().String())
		common.Infof("subscription watcher start scanning for outdated subscriptions at " + time.Now().String())
		subs, err := dbutils.GetSubscriptionsByEndTime(start, s.dbClient, ctx)
		if err != nil {
			common.Error(err)
			sentry.CaptureMessage("error getting expiring subscriptions: " + err.Error())
		} else {
			for _, sub := range subs {
				template := &pb.SubscriptionCancellationEmailTemplate{
					ProviderName: sub.SubscriberName,
					CancelDate:   sub.EndTime.Format("01/02/2006"),
				}
				err = publisher.GetPublisher().SendSubscriptionCancellationEmail(template, sub.Email, strconv.Itoa(sub.ID))
				if err != nil {
					common.Error(err)
					sentry.CaptureMessage("error sending cancellation email: " + err.Error())
				}
			}
		}
		now = time.Now()
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1)
		// routine will start at 00:00 every day
		duration = time.Until(start)
	}
}

func (s *ServiceshipService) recoverSubscriptionWatcher(ctx context.Context) {
	if r := recover(); r != nil {
		common.Error(fmt.Errorf("panic in watcher %v\n", r))
		sentry.CaptureMessage("watcher panicked at " + time.Now().String())
		WatcherRecoverCount += 1
		if WatcherRecoverCount == MaxRecoverCount {
			WatcherRecoverCount = 0
			sentry.CaptureMessage("watcher recovers too many times")
			common.Fatal(fmt.Errorf("too many recovers"))
		} else {
			go s.watchSubscriptions(ctx)
		}
	}
}

func (s *ServiceshipService) GetSubscriptionTransactionInfos(subscriptionID int, ctx context.Context) (*external.Transactions, error) {
	id, t := s.getSubscriptionAccountInfo(subscriptionID, ctx)
	return external.GetChargingService().GetSubscriptionTransactionInfos(subscriptionID, t, id)
}

func (s *ServiceshipService) getSubscriptionAccountInfo(subscriptionID int, ctx context.Context) (int64, string) {
	sub, err := dbutils.GetSubscriptionByID(subscriptionID, s.dbClient, ctx)
	if err != nil {
		return 0, ""
	}
	return sub.AccountID, string(sub.AccountType)
}

func getSubscriptionEndTime(plan *ent.ServiceshipBillingPlan, start time.Time) time.Time {
	switch plan.Interval {
	case serviceshipbillingplan.IntervalMonthly:
		return start.AddDate(0, int(plan.BillingCycle), -7)
	case serviceshipbillingplan.IntervalDaily:
		return start.AddDate(0, 0, int(plan.BillingCycle)-7)
	}
	return time.Time{}
}
