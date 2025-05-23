package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/pendingordercredits"
)

func CreateOrderCredits(orderID int64, creditID int64, clinicID int64, client *ent.Client, ctx context.Context) error {
	_, err := client.PendingOrderCredits.Create().SetOrderID(orderID).SetCreditID(creditID).SetClinicID(clinicID).Save(ctx)
	return err
}

func GetCreditByOrder(orderID int64, client *ent.Client, ctx context.Context) (*ent.PendingOrderCredits, error) {
	return client.PendingOrderCredits.Query().Where(
		pendingordercredits.OrderID(orderID),
	).Only(ctx)
}

func DeleteCreditsByID(id int, client *ent.Client, ctx context.Context) error {
	_, err := client.PendingOrderCredits.Delete().Where(
		pendingordercredits.ID(id),
	).Exec(ctx)
	return err
}
