package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/clinic"
	"coresamples/ent/customer"
	pb "coresamples/proto"
	"fmt"
	"time"
)

// AddSettingToClinic ensures that the clinic-setting relationship exists within a transaction.
func AddSettingToClinic(settingId int, clinicId int, client *ent.Client, ctx context.Context) error {
	// ✅ Start a transaction
	tx, err := client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// ✅ Step 1: Temporarily disable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to disable foreign key checks: %w", err))
	}

	// ✅ Step 2: Add new setting using Ent ORM
	err = tx.Clinic.UpdateOneID(clinicId).
		AddClinicSettingIDs(settingId).
		Exec(ctx)
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to insert clinic setting (clinic_id: %d, setting_id: %d): %w", clinicId, settingId, err))
	}

	// ✅ Step 3: Re-enable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=1")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to re-enable foreign key checks: %w", err))
	}

	// ✅ Commit the transaction
	return tx.Commit()
}

// RemoveSettingFromClinic safely removes a setting-clinic relationship within a transaction.
func RemoveSettingFromClinic(settingId int, clinicId int, client *ent.Client, ctx context.Context) error {
	// ✅ Start a transaction
	tx, err := client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// ✅ Step 1: Temporarily disable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to disable foreign key checks: %w", err))
	}

	// ✅ Step 2: Remove setting using Ent ORM
	err = tx.Clinic.UpdateOneID(clinicId).
		RemoveClinicSettingIDs(settingId).
		Exec(ctx)
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to remove clinic setting (clinic_id: %d, setting_id: %d): %w", clinicId, settingId, err))
	}

	// ✅ Step 3: Re-enable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=1")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to re-enable foreign key checks: %w", err))
	}

	// ✅ Commit the transaction
	return tx.Commit()
}

// UpdateSettingToClinic updates a clinic’s setting while handling foreign key constraints within a transaction.
func UpdateSettingToClinic(settingId int, oldSettingId int, clinicId int, client *ent.Client, ctx context.Context) error {
	// ✅ Start a transaction
	tx, err := client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// ✅ Step 1: Temporarily disable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to disable foreign key checks: %w", err))
	}

	// ✅ Step 2: Perform update using Ent ORM API
	err = tx.Clinic.UpdateOneID(clinicId).
		RemoveClinicSettingIDs(oldSettingId).
		AddClinicSettingIDs(settingId).
		Exec(ctx)
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to update clinic setting (clinic_id: %d, old_setting_id: %d, new_setting_id: %d): %w", clinicId, oldSettingId, settingId, err))
	}

	// ✅ Step 3: Re-enable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=1")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to re-enable foreign key checks: %w", err))
	}

	// ✅ Commit the transaction
	return tx.Commit()
}

func UpdateClinicToCustomerCDC(data *pb.ClinicToCustomerCDCUpdate_Data, old *pb.ClinicToCustomerCDCUpdate_Data, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}

	// old.B, aka customer id is not zero, meaning clinic with data.A as ID has changed its customer
	err = tx.Clinic.Update().Where(clinic.IDEQ(int(data.A))).RemoveCustomerIDs(int(old.B)).Exec(ctx)
	if err != nil {
		return Rollback(tx, err)
	}

	err = tx.Clinic.Update().Where(clinic.IDEQ(int(data.A))).AddCustomerIDs(int(data.B)).Exec(ctx)
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

func UpdateCustomerToClinicCDC(data *pb.ClinicToCustomerCDCUpdate_Data, old *pb.ClinicToCustomerCDCUpdate_Data, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}

	err = tx.Customer.Update().Where(customer.IDEQ(int(data.B))).RemoveClinicIDs(int(old.A)).Exec(ctx)
	if err != nil {
		return Rollback(tx, err)
	}

	err = tx.Customer.Update().Where(customer.IDEQ(int(data.B))).AddClinicIDs(int(data.A)).Exec(ctx)
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

func CreateClinicToCustomerCDC(data *pb.ClinicToCustomerCDCUpdate_Data, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}

	err = tx.Clinic.Update().Where(clinic.IDEQ(int(data.A))).AddCustomerIDs(int(data.B)).Exec(ctx)
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

func DeleteClinicToCustomerCDC(data *pb.ClinicToCustomerCDCUpdate_Data, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}

	err = tx.Clinic.Update().Where(clinic.IDEQ(int(data.A))).RemoveCustomerIDs(int(data.B)).Exec(ctx)
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

func UpsertClinic(ctx context.Context, data *pb.ClinicCDCUpdate_Data, client *ent.Client) error {
	if data == nil || client == nil {
		return fmt.Errorf("data:%v or ent.Client:%v is nil", data, client)
	}

	signTimeStr, err := time.Parse("2006-01-02 15:04:05", data.ClinicSignupTime)
	if err != nil {
		return fmt.Errorf("parse time:%v failed, err:%s", data.ClinicSignupTime, err.Error())
	}

	updatedTimeStr, err := time.Parse("2006-01-02 15:04:05", data.ClinicUpdatedTime)
	if err != nil {
		return fmt.Errorf("parse time:%v failed, err:%s", data.ClinicSignupTime, err.Error())
	}
	// ✅ Start a transaction
	tx, err := client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// ✅ Step 1: Temporarily disable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to disable foreign key checks: %w", err))
	}

	isActive := data.IsActive == 1

	// ✅ Step 2: Perform update using Ent ORM API
	err = tx.Clinic.Create().SetID(int(data.ClinicId)).
		SetClinicName(data.ClinicName).
		SetUserID(int(data.UserId)).
		SetIsActive(isActive).
		SetClinicAccountID(int(data.ClinicAccountId)).
		SetClinicNameOldSystem(data.ClinicNameOldSystem).
		SetClinicSignupTime(signTimeStr).SetClinicUpdatedTime(updatedTimeStr).
		OnConflict().UpdateNewValues().
		Exec(ctx)
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to update clinic (clinic_id: %d, ): %w", data.ClinicId, err))
	}

	// ✅ Step 3: Re-enable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=1")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to re-enable foreign key checks: %w", err))
	}

	// ✅ Commit the transaction
	return tx.Commit()
}

func DeleteClinic(ctx context.Context, clinicID int32, client *ent.Client) error {
	if client == nil {
		return fmt.Errorf("ent.Client is nil")
	}
	// ✅ Start a transaction
	tx, err := client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// ✅ Step 1: Temporarily disable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to disable foreign key checks: %w", err))
	}

	// ✅ Step 2: Perform update using Ent ORM API
	err = tx.Clinic.DeleteOneID(int(clinicID)).Exec(ctx)
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to delete clinic (clinic_id: %d): %w", clinicID, err))
	}

	// ✅ Step 3: Re-enable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=1")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to re-enable foreign key checks: %w", err))
	}

	// ✅ Commit the transaction
	return tx.Commit()
}

func CreateClinicToPatientCDC(data *pb.ClinicToPatientCDCUpdate_Data, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}

	err = tx.Clinic.Update().Where(clinic.IDEQ(int(data.A))).AddClinicPatientIDs(int(data.B)).Exec(ctx)
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

func UpdateClinicToPatientCDC(data *pb.ClinicToPatientCDCUpdate_Data, old *pb.ClinicToPatientCDCUpdate_Data, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}

	// Remove the old patient from the clinic
	err = tx.Clinic.Update().Where(clinic.IDEQ(int(old.A))).RemoveClinicPatientIDs(int(old.B)).Exec(ctx)
	if err != nil {
		return Rollback(tx, err)
	}

	// Add the new patient to the clinic
	err = tx.Clinic.Update().Where(clinic.IDEQ(int(data.A))).AddClinicPatientIDs(int(data.B)).Exec(ctx)
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

func DeleteClinicToPatientCDC(data *pb.ClinicToPatientCDCUpdate_Data, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}

	err = tx.Clinic.Update().Where(clinic.IDEQ(int(data.A))).RemoveClinicPatientIDs(int(data.B)).Exec(ctx)
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
