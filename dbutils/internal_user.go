package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/internaluser"
	pb "coresamples/proto"
	"fmt"
)

func GetSalesByEmail(email string, client *ent.Client, ctx context.Context) (*ent.InternalUser, error) {
	return client.InternalUser.Query().Where(
		internaluser.InternalUserEmailEqualFold(email),
	).First(ctx)
}

func UpsertInternalUserCDC(data *pb.InternalUserCDCUpdate_Data, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	if data.InternalUserId == 0 {
		return fmt.Errorf("cdc customer message doesn't have customer id %v", data)
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}

	internalUserIsFullTime := data.InternalUserIsFullTime == 1
	isActive := data.IsActive == 1

	err = tx.InternalUser.Create().SetID(int(data.InternalUserId)).
		SetInternalUserRole(data.InternalUserRole).
		SetInternalUserName(data.InternalUserName).
		SetInternalUserFirstname(data.InternalUserFirstname).
		SetInternalUserLastname(data.InternalUserLastname).
		SetInternalUserMiddleName(data.InternalUserMiddlename).
		SetInternalUserIsFullTime(internalUserIsFullTime).
		SetInternalUserEmail(data.InternalUserEmail).
		SetInternalUserPhone(data.InternalUserPhone).
		SetIsActive(isActive).
		SetUserID(int(data.UserId)).
		SetInternalUserType(data.InternalUserType).
		SetInternalUserRoleID(int(data.InternalUserRoleId)).
		OnConflict().UpdateNewValues().
		Exec(ctx)
	if err != nil {
		return Rollback(tx, err)
	}

	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=1;")
	if err != nil {
		return Rollback(tx, err)
	}
	err = tx.Commit()
	if err != nil {
		return Rollback(tx, err)
	}
	return nil
}

func DeleteInternalUserCDC(data *pb.InternalUserCDCUpdate_Data, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	if data.InternalUserId == 0 {
		return fmt.Errorf("cdc internal_user message doesn't have internal user id %v", data)
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}

	_, err = client.InternalUser.Delete().
		Where(internaluser.ID(int(data.InternalUserId))).
		Exec(ctx)
	if err != nil {
		return Rollback(tx, err)
	}

	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=1;")
	if err != nil {
		return Rollback(tx, err)
	}
	err = tx.Commit()
	if err != nil {
		return Rollback(tx, err)
	}
	return nil
}
