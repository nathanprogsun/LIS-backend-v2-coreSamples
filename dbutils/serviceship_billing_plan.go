package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/accountsubscription"
	"coresamples/ent/serviceship"
	"coresamples/ent/serviceshipbillingplan"
	"coresamples/util"
	"errors"
	"time"
)

func AddBillingPlan(serviceshipID int, fee float32, billingCycle int32, interval string, client *ent.Client, ctx context.Context) error {
	latestPlan, err := client.ServiceshipBillingPlan.Query().
		Where(serviceshipbillingplan.HasServiceshipWith(serviceship.IDEQ(serviceshipID))).
		Order(ent.Desc(serviceshipbillingplan.FieldEffectiveTime)).
		First(ctx)
	if err != nil {
		return err
	}
	effectiveTime := latestPlan.EffectiveTime
	cycleExist, err := client.ServiceshipBillingPlan.Query().
		Where(serviceshipbillingplan.And(
			serviceshipbillingplan.EffectiveTimeEQ(effectiveTime),
			serviceshipbillingplan.BillingCycleEQ(billingCycle),
			serviceshipbillingplan.HasServiceshipWith(serviceship.IDEQ(serviceshipID)),
		)).Exist(ctx)
	if err != nil {
		return nil
	}
	if cycleExist {
		return errors.New("billing cycle for the latest billing plan already exists")
	}
	_, err = client.ServiceshipBillingPlan.Create().
		SetServiceshipID(serviceshipID).
		SetFee(fee).
		SetBillingCycle(billingCycle).
		SetInterval(serviceshipbillingplan.Interval(interval)).
		SetEffectiveTime(effectiveTime).Save(ctx)
	return err
}

func CreateBillingPlanSet(serviceshipID int, fee []float32, billingCycle []int32, interval []string, effectiveTime time.Time, client *ent.Client, ctx context.Context) error {
	if len(fee) != len(billingCycle) {
		return errors.New("billing cycles should match monthly fee")
	}
	if !util.ElementsUniqueInt32(billingCycle) {
		return errors.New("billing cycles should be unique")
	}
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	for idx := range fee {
		_, err = tx.ServiceshipBillingPlan.Create().
			SetServiceshipID(serviceshipID).
			SetFee(fee[idx]).
			SetBillingCycle(billingCycle[idx]).
			SetInterval(serviceshipbillingplan.Interval(interval[idx])).
			SetEffectiveTime(effectiveTime).
			Save(ctx)
		if err != nil {
			return Rollback(tx, err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return Rollback(tx, err)
	}
	return nil
}

func GetBillingPlanByID(billingPlanID int, client *ent.Client, ctx context.Context) (*ent.ServiceshipBillingPlan, error) {
	return client.ServiceshipBillingPlan.Get(ctx, billingPlanID)
}

func GetBillingPlanBySubscriptionID(subscriptionID int, client *ent.Client, ctx context.Context) (*ent.ServiceshipBillingPlan, error) {
	return client.ServiceshipBillingPlan.Query().Where(
		serviceshipbillingplan.HasAccountSubscriptionWith(
			accountsubscription.ID(subscriptionID)),
	).WithServiceship().Only(ctx)
}

func GetLatestBillingPlanSet(serviceshipID int, client *ent.Client, ctx context.Context) ([]*ent.ServiceshipBillingPlan, error) {
	latestPlan, err := client.ServiceshipBillingPlan.Query().
		Where(serviceshipbillingplan.HasServiceshipWith(serviceship.IDEQ(serviceshipID))).
		Order(ent.Desc(serviceshipbillingplan.FieldEffectiveTime)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	effectiveTime := latestPlan.EffectiveTime
	return client.ServiceshipBillingPlan.Query().
		Where(serviceshipbillingplan.And(
			serviceshipbillingplan.HasServiceshipWith(serviceship.IDEQ(serviceshipID)),
			serviceshipbillingplan.EffectiveTimeEQ(effectiveTime),
		)).All(ctx)
}

func GetLatestEffectiveBillingPlanSet(serviceshipID int, client *ent.Client, ctx context.Context) ([]*ent.ServiceshipBillingPlan, error) {
	latestPlan, err := client.ServiceshipBillingPlan.Query().
		Where(serviceshipbillingplan.And(
			serviceshipbillingplan.HasServiceshipWith(serviceship.IDEQ(serviceshipID)),
			serviceshipbillingplan.EffectiveTimeLTE(time.Now()),
		)).
		Order(ent.Desc(serviceshipbillingplan.FieldEffectiveTime)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	effectiveTime := latestPlan.EffectiveTime
	return client.ServiceshipBillingPlan.Query().
		Where(serviceshipbillingplan.And(
			serviceshipbillingplan.HasServiceshipWith(serviceship.IDEQ(serviceshipID)),
			serviceshipbillingplan.EffectiveTimeEQ(effectiveTime),
		)).All(ctx)
}
