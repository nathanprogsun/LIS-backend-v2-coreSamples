package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/serviceship"
	"coresamples/ent/serviceshipbillingplan"
	"coresamples/util"
	"errors"
	"time"
)

func CreateServiceship(
	tag string,
	serviceType string,
	fee []float32,
	billingCycles []int32,
	billingIntervals []string,
	client *ent.Client,
	ctx context.Context) error {
	if !util.ElementsUniqueInt32(billingCycles) {
		return errors.New("billing cycles should be unique")
	}
	if len(fee) != len(billingCycles) || len(fee) != len(billingIntervals) {
		return errors.New("monthly fee should match billing cycles")
	}
	txClient, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	svcship, err := txClient.Serviceship.Create().
		SetTag(tag).SetType(serviceship.Type(serviceType)).Save(ctx)
	if err != nil {
		return Rollback(txClient, err)
	}

	now := time.Now()

	for idx := range fee {
		_, err = txClient.ServiceshipBillingPlan.Create().
			SetServiceship(svcship).
			SetFee(fee[idx]).
			SetBillingCycle(billingCycles[idx]).
			SetInterval(serviceshipbillingplan.Interval(billingIntervals[idx])).
			SetEffectiveTime(now).Save(ctx)
		if err != nil {
			return Rollback(txClient, err)
		}
	}

	err = txClient.Commit()
	if err != nil {
		return Rollback(txClient, err)
	}
	return nil
}

func GetServiceships(client *ent.Client, ctx context.Context) ([]*ent.Serviceship, error) {
	return client.Serviceship.Query().All(ctx)
}

func GetServiceshipByID(id int, client *ent.Client, ctx context.Context) (*ent.Serviceship, error) {
	return client.Serviceship.Query().Where(serviceship.IDEQ(id)).Only(ctx)
}

func GetServiceshipByTag(tag string, svcType string, client *ent.Client, ctx context.Context) (*ent.Serviceship, error) {
	return client.Serviceship.Query().Where(
		serviceship.TagEQ(tag),
		serviceship.TypeEQ(serviceship.Type(svcType)),
	).Only(ctx)
}

func GetServiceshipsByType(svcType string, client *ent.Client, ctx context.Context) ([]*ent.Serviceship, error) {
	return client.Serviceship.Query().Where(
		serviceship.TypeEQ(serviceship.Type(svcType)),
	).All(ctx)
}
