package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/accountsubscription"
	"coresamples/ent/serviceship"
	"coresamples/ent/serviceshipbillingplan"
	"errors"
	"strconv"
	"time"
)

// CreateSubscription Create a new subscription for user
// The key idea is to ensure at most only one membership per subscription under one account, whether expired or not
// And there cannot be two non-single memberships under one account
// It returns error if the account wants to subscribe to another non-single unexpired membership,
// or the account has an unexpired subscription to the required membership
// If the account has an expired service membership, it updates the start and end date
func CreateSubscription(
	accountID int64,
	accountType string,
	subscriberName string,
	email string,
	membershipBillingPlanId int,
	keepActive bool,
	client *ent.Client,
	ctx context.Context) (*ent.AccountSubscription, error) {
	txClient, err := client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	// get billing plan and corresponding membership
	billingPlan, err := txClient.ServiceshipBillingPlan.Query().
		Where(serviceshipbillingplan.IDEQ(membershipBillingPlanId)).WithServiceship().Only(ctx)
	if err != nil {
		return nil, Rollback(txClient, err)
	}
	svcship, err := billingPlan.Edges.ServiceshipOrErr()
	if err != nil {
		return nil, Rollback(txClient, err)
	}
	svcshipID := svcship.ID

	// Check if the account has/had a subscription to this membership
	subscription, err := txClient.AccountSubscription.Query().Where(
		accountsubscription.And(
			accountsubscription.AccountID(accountID),
			accountsubscription.AccountTypeEQ(accountsubscription.AccountType(accountType)),
			accountsubscription.HasServiceshipWith(serviceship.IDEQ(svcshipID)),
		)).First(ctx)
	now := time.Now()
	var sub *ent.AccountSubscription
	if subscription != nil {
		if subscription.EndTime.After(time.Now()) || subscription.EndTime.IsZero() {
			return nil, Rollback(txClient, errors.New("account "+strconv.Itoa(int(accountID))+" already has the membership subscribed"))
		}
		// has an expired subscription
		builder := subscription.Update().
			SetSubscriberName(subscriberName).
			SetEmail(email).
			SetStartTime(now).
			SetServiceshipID(svcshipID).
			SetServiceshipBillingPlan(billingPlan)
		if !keepActive {
			builder.SetEndTime(now.AddDate(0, int(billingPlan.BillingCycle), 0))
		}
		sub, err = builder.Save(ctx)
		if err != nil {
			return nil, Rollback(txClient, err)
		}
	} else {
		builder := txClient.AccountSubscription.Create().
			SetAccountID(accountID).
			SetAccountType(accountsubscription.AccountType(accountType)).
			SetSubscriberName(subscriberName).
			SetEmail(email).
			SetStartTime(now).
			SetServiceshipID(svcshipID).
			SetServiceshipBillingPlan(billingPlan)
		if !keepActive {
			builder.SetEndTime(now.AddDate(0, int(billingPlan.BillingCycle), 0))
		}
		sub, err = builder.Save(ctx)
		if err != nil {
			return nil, Rollback(txClient, err)
		}
	}
	err = txClient.Commit()
	if err != nil {
		return nil, Rollback(txClient, err)
	}
	return sub, nil
}

// CancelSubscription this is soft delete, it merely sets the end time to now and turn off auto-renewal
// It does nothing if a valid subscription does not exist
func CancelSubscription(accountID int64, accountType string, svcshipID int, client *ent.Client, ctx context.Context) error {
	return client.AccountSubscription.Update().Where(
		accountsubscription.And(
			accountsubscription.AccountIDEQ(accountID),
			accountsubscription.AccountTypeEQ(accountsubscription.AccountType(accountType)),
			accountsubscription.EndTimeGT(time.Now()),
			accountsubscription.HasServiceshipWith(serviceship.IDEQ(svcshipID)),
		)).
		SetEndTime(time.Now()).
		Exec(ctx)
}

func DeleteSubscription(subscriptionID int, client *ent.Client, ctx context.Context) error {
	_, err := client.AccountSubscription.Delete().Where(accountsubscription.ID(subscriptionID)).Exec(ctx)
	return err
}

// UpdateSubscriptionBillingPlan does nothing if the subscription does not exist or has expired
func UpdateSubscriptionBillingPlan(subscriptionID int, billingPlanID int, client *ent.Client, ctx context.Context) error {
	return client.AccountSubscription.Update().Where(
		accountsubscription.IDEQ(subscriptionID),
		accountsubscription.Or(
			accountsubscription.EndTimeGT(time.Now()),
			accountsubscription.EndTimeEQ(time.Time{}),
		),
	).SetServiceshipBillingPlanID(billingPlanID).Exec(ctx)
}

func UpdateSubscriptionEmail(subscriptionID int, email string, client *ent.Client, ctx context.Context) error {
	return client.AccountSubscription.Update().Where(
		accountsubscription.ID(subscriptionID),
	).SetEmail(email).Exec(ctx)
}

func UpdateAllSubscriptionsEmail(accountID int64, accountType string, email string, client *ent.Client, ctx context.Context) error {
	return client.AccountSubscription.Update().Where(
		accountsubscription.AccountIDEQ(accountID),
		accountsubscription.AccountTypeEQ(accountsubscription.AccountType(accountType)),
		accountsubscription.Or(
			accountsubscription.EndTimeGT(time.Now()),
			accountsubscription.EndTimeEQ(time.Time{}),
		),
	).SetEmail(email).Exec(ctx)
}

func GetSubscriptionByID(subscriptionID int, client *ent.Client, ctx context.Context) (*ent.AccountSubscription, error) {
	return client.AccountSubscription.Query().Where(accountsubscription.ID(subscriptionID)).WithServiceshipBillingPlan().WithServiceship().Only(ctx)
}

// GetSubscriptionsByEndTime find the subscriptions that ended within one day before the timestamp
// i.e. timestamp - 1 <= endtime < timestamp day
func GetSubscriptionsByEndTime(timestamp time.Time, client *ent.Client, ctx context.Context) ([]*ent.AccountSubscription, error) {
	return client.AccountSubscription.Query().Where(
		accountsubscription.And(
			accountsubscription.EndTimeGTE(timestamp.AddDate(0, 0, -1)),
			accountsubscription.EndTimeLT(timestamp)),
	).All(ctx)
}

func AddSubscriptionEndTime(months int32, subscriptionID int, client *ent.Client, ctx context.Context) error {
	txClient, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	subscription, err := txClient.AccountSubscription.Query().Where(
		accountsubscription.ID(subscriptionID),
	).Only(ctx)
	if err != nil {
		return Rollback(txClient, err)
	}

	timestamp := subscription.EndTime.AddDate(0, int(months), 0)
	err = subscription.Update().SetEndTime(timestamp).Exec(ctx)
	if err != nil {
		return Rollback(txClient, err)
	}

	err = txClient.Commit()
	if err != nil {
		return Rollback(txClient, err)
	}
	return nil
}

func MaybeMatchNewestBillingPlans(billingCycles []int32, fee []float32, intervals []string, effectiveTime time.Time, svcshipID int, client *ent.Client, ctx context.Context) ([]*ent.AccountSubscription, error) {
	if len(billingCycles) != len(fee) {
		return nil, errors.New("unmatched billing cycles and monthly fee")
	}
	tx, err := client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var res []*ent.AccountSubscription
	for idx := range billingCycles {
		currPlan, err := tx.ServiceshipBillingPlan.Query().Where(
			serviceshipbillingplan.And(
				serviceshipbillingplan.BillingCycleEQ(billingCycles[idx]),
				serviceshipbillingplan.FeeEQ(fee[idx]),
				serviceshipbillingplan.IntervalEQ(serviceshipbillingplan.Interval(intervals[idx])),
				serviceshipbillingplan.EffectiveTimeEQ(effectiveTime),
				serviceshipbillingplan.HasServiceshipWith(serviceship.ID(svcshipID)),
			),
		).Only(ctx)
		if err != nil {
			return nil, Rollback(tx, err)
		}
		subs, err := tx.AccountSubscription.Query().Where(
			accountsubscription.And(
				accountsubscription.EndTimeLTE(effectiveTime),
				accountsubscription.EndTimeNEQ(time.Time{}),
				accountsubscription.HasServiceshipWith(serviceship.ID(svcshipID)),
				accountsubscription.HasServiceshipBillingPlanWith(
					serviceshipbillingplan.And(
						serviceshipbillingplan.BillingCycleEQ(billingCycles[idx]),
						serviceshipbillingplan.IntervalEQ(serviceshipbillingplan.Interval(intervals[idx])),
						serviceshipbillingplan.FeeGT(fee[idx]),
					)),
			),
		).All(ctx)
		if err != nil {
			return nil, Rollback(tx, err)
		}
		if subs == nil || len(subs) == 0 {
			continue
		}
		res = append(res, subs...)
		for _, sub := range subs {
			err := tx.AccountSubscription.Update().
				Where(accountsubscription.IDEQ(sub.ID)).
				SetServiceshipBillingPlan(currPlan).
				Exec(ctx)
			if err != nil {
				return nil, Rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, Rollback(tx, err)
	}
	return res, nil
}

func GetAccountSubscriptions(accountID int64, accountType string, includeOutdated bool, client *ent.Client, ctx context.Context) ([]*ent.AccountSubscription, error) {
	if includeOutdated {
		return client.AccountSubscription.Query().
			Where(
				accountsubscription.AccountIDEQ(accountID),
				accountsubscription.AccountTypeEQ(accountsubscription.AccountType(accountType)),
			).WithServiceship().WithServiceshipBillingPlan().All(ctx)
	}
	return client.AccountSubscription.Query().
		Where(accountsubscription.And(
			accountsubscription.AccountIDEQ(accountID),
			accountsubscription.AccountTypeEQ(accountsubscription.AccountType(accountType)),
			accountsubscription.Or(
				accountsubscription.EndTimeGT(time.Now()),
				accountsubscription.EndTimeEQ(time.Time{}),
			),
		)).WithServiceship().WithServiceshipBillingPlan().All(ctx)
}

// GetAccountServiceshipSubscriptionsByType returns valid non-single service membership subscription of an account
func GetAccountServiceshipSubscriptionsByType(accountID int64, accountType string, serviceType string, client *ent.Client, ctx context.Context) ([]*ent.AccountSubscription, error) {
	return client.AccountSubscription.Query().
		Where(
			accountsubscription.And(
				accountsubscription.HasServiceshipWith(serviceship.TypeEQ(serviceship.Type(serviceType))),
				accountsubscription.AccountIDEQ(accountID),
				accountsubscription.AccountTypeEQ(accountsubscription.AccountType(accountType)),
				accountsubscription.EndTimeGT(time.Now()),
			)).WithServiceship().WithServiceshipBillingPlan().All(ctx)
}
