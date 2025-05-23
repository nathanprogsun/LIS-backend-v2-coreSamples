package subscriber

import (
	"context"
	"coresamples/common"
	"coresamples/config"
	"coresamples/dbutils"
	"coresamples/ent"
	"coresamples/external"
	pb "coresamples/proto"
	"coresamples/publisher"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

type MembershipEventHandler struct {
	dbClient *ent.Client
	ctx      context.Context
}

func (oh *MembershipEventHandler) HandleMembershipGeneralEvent(event *pb.GeneralEvent) {
	err := oh.handleEventByName(event)
	if err != nil {
		common.Error(err)
	}
	err = oh.handleOrderPlacementEvent(event)
	if err != nil {
		common.Error(err)
	}
}

func (oh *MembershipEventHandler) handleOrderPlacementEvent(event *pb.GeneralEvent) error {
	if event.GetEventName() != "order_status_updates" {
		return nil
	}
	switch event.GetAddonColumn().GetOrderStatus() {
	case "new_order":
		amount := event.GetAddonColumn().Amount
		clinicID, err := strconv.ParseInt(event.GetAddonColumn().ClinicId, 10, 64)
		if err != nil {
			common.Infof("getting new order event with no clinic id %v", event)
			return nil
		}
		customerID, err := strconv.ParseInt(event.GetAddonColumn().CustomerId, 10, 64)
		if err != nil {
			return err
		}
		patientID, err := strconv.ParseInt(event.GetAddonColumn().PatientId, 10, 64)
		if err != nil {
			return err
		}
		orderID := event.GetAddonColumn().OrderId
		subs, err := dbutils.GetAccountServiceshipSubscriptionsByType(clinicID, "clinic", "membership", oh.dbClient, oh.ctx)
		if err != nil || subs == nil {
			return nil
		}
		for _, sub := range subs {
			svc, err := sub.Edges.ServiceshipOrErr()
			if err != nil {
				common.Errorf("unable to find service", err)
				continue
			}
			mem, ok := config.GetServiceshipConfig().Membership[svc.Tag]
			if !ok {
				common.Error(fmt.Errorf("unable to find membership"))
				continue
			}
			credits, err := oh.getCredits(mem, amount, patientID, clinicID, event.EventTime)
			if err != nil {
				common.Error(err)
				continue
			}
			creditID, err := external.GetAccountingService().AddCredits(decimal.NewFromFloat(float64(credits)).String(), orderID, customerID)
			if err != nil {
				sentry.CaptureMessage(err.Error())
				common.Error(err)
				continue
			}
			common.Infof("event %s: add credits of %f for clinic %d, customer %d, patient %d", event.EventId, credits, clinicID, customerID, patientID)
			err = dbutils.CreateOrderCredits(orderID, creditID, clinicID, oh.dbClient, oh.ctx)
			if err != nil {
				common.Error(err)
			}
		}
	case "order_canceled":
		orderID := event.GetAddonColumn().OrderId
		clinicID, err := strconv.ParseInt(event.GetAddonColumn().ClinicId, 10, 64)
		if err != nil {
			return err
		}
		rec, err := dbutils.GetCreditByOrder(orderID, oh.dbClient, oh.ctx)
		if err != nil || rec == nil {
			return nil
		}
		err = dbutils.DeleteCreditsByID(rec.ID, oh.dbClient, oh.ctx)
		if err != nil {
			return err
		}
		common.Infof("event %s: void credits for clinic %d", event.EventId, clinicID)
		return external.GetAccountingService().VoidCredits(rec.CreditID, clinicID, "clinic")
	}
	return nil
}

func (oh *MembershipEventHandler) handleEventByName(event *pb.GeneralEvent) error {
	switch event.GetEventName() {
	case "subscription_payment_updates":
		if event.GetAddonColumn().GetChargeType() != "subscription" || event.GetAddonColumn().GetAccountType() != "clinic" {
			return nil
		}
		clinicID, err := strconv.ParseInt(event.GetAddonColumn().GetAccountId(), 10, 64)
		if err != nil {
			return err
		}
		subscriptionID, err := strconv.ParseInt(event.GetAddonColumn().GetChargeTypeId(), 10, 64)
		if err != nil {
			return err
		}
		sub, err := dbutils.GetSubscriptionByID(int(subscriptionID), oh.dbClient, oh.ctx)
		if err != nil {
			return err
		}
		plan, err := sub.Edges.ServiceshipBillingPlanOrErr()
		if err != nil {
			return err
		}
		status := event.GetAddonColumn().GetStatus()
		if status == "succeed" {
			// extend subscription
			err = dbutils.AddSubscriptionEndTime(plan.BillingCycle, int(subscriptionID), oh.dbClient, oh.ctx)
			if err != nil {
				return err
			}
			common.Infof("event %s: extend subscription %d by %d month for clinic %d", event.EventId, subscriptionID, plan.BillingCycle, clinicID)
		} else if status == "fail" {
			// pause subscription
			err = external.GetChargingService().PauseSubscription(int(subscriptionID), "clinic", clinicID)
			if err != nil {
				return err
			}
			// send notification
			template := &pb.PaymentUpdateEmailTemplate{
				ProviderName: sub.SubscriberName,
			}
			err = publisher.GetPublisher().SendPaymentUpdateEmail(template, publisher.EmailFrom, event.GetAddonColumn().GetChargeTypeId())
			if err != nil {
				common.Error(err)
			}
			return err
		}
	case "receive_sample_tubes":
		if event.GetEventProvider() != "lis-shipping" {
			return nil
		}
		orderID := event.GetAddonColumn().OrderId
		rec, err := dbutils.GetCreditByOrder(orderID, oh.dbClient, oh.ctx)
		if err != nil {
			return nil
		}
		err = dbutils.DeleteCreditsByID(rec.ID, oh.dbClient, oh.ctx)
		if err != nil {
			return err
		}
		err = external.GetAccountingService().ActivateCredits(rec.CreditID, rec.ClinicID, "clinic")
		if err != nil {
			return err
		}
		common.Infof("event %s: activate credits %d for clinic %d", event.EventId, rec.CreditID, rec.ClinicID)
		return nil
	}
	return nil
}

func (oh *MembershipEventHandler) getCredits(mem config.MembershipConfig, amount float32, patientID int64, clinicID int64, eventTime string) (float32, error) {
	creditRate := float32(mem.Bonus["credit_rate"].(float64))
	credits := amount * creditRate
	repeatCreditRate := float32(mem.Bonus["repeat_patient_bonus_rate"].(float64))
	if repeatCreditRate > 0 {
		orders, err := external.GetOrderService().GetLatestOrders(patientID)
		if err != nil {
			sentry.CaptureMessage(err.Error())
			return credits, nil
		}
		now, err := time.Parse(time.RFC3339, eventTime)
		if err != nil {
			return 0, err
		}
		for _, order := range orders {
			orderDate, err := time.Parse(time.RFC3339, order.OrderCreatedDate)
			if err != nil {
				common.Error(err)
				continue
			}
			// repeat patient
			if now.Before(orderDate.AddDate(0, 6, 0)) && now.After(orderDate) && order.ClinicID == clinicID {
				credits += amount * repeatCreditRate
				break
			}
		}
	}
	return credits, nil
}
