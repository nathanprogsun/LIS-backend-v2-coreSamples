package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/user"
	pb "coresamples/proto"
	"fmt"
)

func UpsertUserCDC(data *pb.UserCDCUpdate_Data, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	if data.UserId == 0 {
		return fmt.Errorf("cdc contact message doesn't have user id %v", data)
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}

	isTwoFactorAuthenticationEnabled := data.IsTwoFactorAuthenticationEnabled == 1
	importedUserWithSaltPassword := data.ImportedUserWithSaltPassword == 1
	isActive := data.IsActive == 1

	err = client.User.Create().SetID(int(data.UserId)).
		SetUserName(data.Username).
		SetEmailUserID(data.EmailUserId).
		SetPassword(data.Password).
		SetTwoFactorAuthenticationSecret(data.TwoFactorAuthenticationSecret).
		SetIsTwoFactorAuthenticationEnabled(isTwoFactorAuthenticationEnabled).
		SetUserGroup(data.UserGroup).
		SetImportedUserWithSaltPassword(importedUserWithSaltPassword).
		SetIsActive(isActive).OnConflict().UpdateNewValues().Exec(ctx)
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

func DeleteUserCDC(data *pb.UserCDCUpdate_Data, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	if data.UserId == 0 {
		return fmt.Errorf("cdc contact message doesn't have contact id %v", data)
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}

	_, err = client.User.Delete().Where(user.ID(int(data.UserId))).Exec(ctx)
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
