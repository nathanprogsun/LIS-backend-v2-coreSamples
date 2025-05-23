package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/customersettingonclinics"
	"coresamples/ent/setting"
	pb "coresamples/proto"
	"fmt"
	"time"
)

func UpsertCustomerSettingOnClinics(customerID int, clinicID int, settingID int, settingName string, client *ent.Client, ctx context.Context) error {
	// ✅ Step 1: Start a transaction
	tx, err := client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// ✅ Step 2: Temporarily disable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to disable foreign key checks: %w", err))
	}

	// ✅ Step 3: Perform UPSERT using Ent ORM
	err = tx.CustomerSettingOnClinics.Create().
		SetCustomerID(customerID).
		SetClinicID(clinicID).
		SetSettingID(settingID).
		SetSettingName(settingName).
		OnConflictColumns(customersettingonclinics.FieldCustomerID, customersettingonclinics.FieldClinicID, customersettingonclinics.FieldSettingName).
		Update(func(u *ent.CustomerSettingOnClinicsUpsert) {
			u.SetSettingID(settingID) // ✅ Only update setting_name if conflict occurs
		}).
		Exec(ctx)
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to insert customer setting: %w", err))
	}

	// ✅ Step 4: Re-enable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=1")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to re-enable foreign key checks: %w", err))
	}

	// ✅ Step 5: Commit transaction
	return tx.Commit()
}

func RemoveCustomerSettingOnClinics(customerID int, clinicID int, settingID int, client *ent.Client, ctx context.Context) error {
	// ✅ Step 1: Start a transaction
	tx, err := client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// ✅ Step 2: Disable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to disable foreign key checks: %w", err))
	}

	// ✅ Step 3: Delete the customer-setting relationship
	_, err = tx.CustomerSettingOnClinics.Delete().
		Where(
			customersettingonclinics.CustomerIDEQ(customerID),
			customersettingonclinics.ClinicIDEQ(clinicID),
			customersettingonclinics.SettingIDEQ(settingID),
		).
		Exec(ctx)
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to remove customer setting: %w", err))
	}

	// ✅ Step 4: Re-enable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=1")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to re-enable foreign key checks: %w", err))
	}

	// ✅ Step 5: Commit transaction
	return tx.Commit()
}

func UpsertSetting(data *pb.SettingCDCUpdate_SettingData, client *ent.Client, ctx context.Context) error {
	// ✅ Step 1: Start a transaction
	tx, err := client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// ✅ Step 2: Ensure valid Setting ID
	if data.SettingId == 0 {
		return fmt.Errorf("invalid SettingID in CDC data: %v", data)
	}

	// ✅ Step 3: Disable foreign key checks temporarily
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to disable foreign key checks: %w", err))
	}

	isActive := data.IsActive == 1
	applyToAllGroupMember := data.ApplyToAllGroupMember == 1
	isOfficial := data.IsOfficial == 1
	useGroupSetting := data.UseGroupSetting == 1

	// ✅ Step 4: Perform UPSERT using Ent ORM
	err = tx.Setting.Create().
		SetID(int(data.SettingId)).
		SetSettingName(data.SettingName).
		SetSettingGroup(data.SettingGroup).
		SetSettingDescription(data.SettingDescription).
		SetSettingValue(data.SettingValue).
		SetSettingType(data.SettingType).
		SetSettingValueUpdatedTime(time.Now()). // Auto-update timestamp
		SetIsActive(isActive).
		SetApplyToAllGroupMember(applyToAllGroupMember).
		SetIsOfficial(isOfficial).
		SetSettingLevel(int(data.SettingLevel)).
		SetSettingLevelName(data.SettingLevelName).
		SetUseGroupSetting(useGroupSetting).
		OnConflict().UpdateNewValues(). // ✅ If conflict, update with new values
		Exec(ctx)

	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to upsert setting: %w", err))
	}

	// ✅ Step 5: Re-enable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=1;")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to re-enable foreign key checks: %w", err))
	}

	// ✅ Step 6: Commit transaction
	return tx.Commit()
}

func DeleteSetting(settingID int, client *ent.Client, ctx context.Context) error {
	// ✅ Step 1: Start a transaction
	tx, err := client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// ✅ Step 2: Temporarily disable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to disable foreign key checks: %w", err))
	}

	// ✅ Step 3: Delete the setting
	_, err = tx.Setting.Delete().
		Where(setting.IDEQ(settingID)).
		Exec(ctx)
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to remove setting (setting_id: %d): %w", settingID, err))
	}

	// ✅ Step 4: Re-enable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=1")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to re-enable foreign key checks: %w", err))
	}

	// ✅ Step 5: Commit transaction
	return tx.Commit()
}
