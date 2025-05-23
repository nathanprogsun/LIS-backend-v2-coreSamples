package dbutils

import (
	"context"
	"coresamples/common"
	"coresamples/ent"
	"coresamples/ent/address"
	"coresamples/ent/customeraddressonclinics"
	"coresamples/model"
	pb "coresamples/proto"
	"fmt"
)

func UpsertAddress(data *pb.AddressCDCUpdate_AddressData, client *ent.Client, ctx context.Context) error {
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

	addressConfirmed := data.AddressConfirmed == 1
	isPrimaryAddress := data.IsPrimaryAddress == 1
	applyToAllGroupMember := data.ApplyToAllGroupMember == 1
	isGroupAddress := data.IsGroupAddress == 1
	useAsDefaultCreateAddress := data.UseAsDefaultCreateAddress == 1
	useGroupAddress := data.UseGroupAddress == 1

	// ✅ Step 3: Perform UPSERT using Ent ORM
	err = tx.Address.Create().
		SetID(int(data.AddressId)).
		SetAddressType(data.AddressType).
		SetStreetAddress(data.StreetAddress).
		SetAptPo(data.AptPo).
		SetCity(data.City).
		SetCountry(data.Country).
		SetAddressConfirmed(addressConfirmed).
		SetIsPrimaryAddress(isPrimaryAddress).
		SetCustomerID(int(data.CustomerId)).
		SetPatientID(int(data.PatientId)).
		SetClinicID(int(data.ClinicId)).
		SetInternalUserID(int(data.InternalUserId)).
		SetAddressLevel(int(data.AddressLevel)).
		SetAddressLevelName(data.AddressLevelName).
		SetGroupAddressID(int(data.GroupAddressId)).
		SetApplyToAllGroupMember(applyToAllGroupMember).
		SetIsGroupAddress(isGroupAddress).
		SetUseAsDefaultCreateAddress(useAsDefaultCreateAddress).
		SetUseGroupAddress(useGroupAddress).
		OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to upsert address: %w", err))
	}

	// ✅ Step 4: Re-enable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=1")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to re-enable foreign key checks: %w", err))
	}

	// ✅ Step 5: Commit transaction
	return tx.Commit()
}

func DeleteAddress(addressId int, client *ent.Client, ctx context.Context) error {
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

	// ✅ Step 3: Delete the address
	_, err = tx.Address.Delete().
		Where(address.IDEQ(addressId)).
		Exec(ctx)
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to remove address (address_id: %d): %w", addressId, err))
	}

	// ✅ Step 4: Re-enable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=1")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to re-enable foreign key checks: %w", err))
	}

	// ✅ Step 5: Commit transaction
	return tx.Commit()
}

func FetchCustomerClinicAddresses(ctx context.Context, db *ent.Client, customerID, clinicID int) []*model.Address {
	var result []*model.Address

	// Query CustomerAddressOnClinics
	relations, err := db.CustomerAddressOnClinics.
		Query().
		Where(
			customeraddressonclinics.CustomerIDEQ(customerID),
			customeraddressonclinics.ClinicIDEQ(clinicID),
		).
		All(ctx)
	if err != nil {
		common.Error(err)
		return result
	}

	// 2. Get all address_id
	var addressIDs []int
	for _, rel := range relations {
		addressIDs = append(addressIDs, rel.AddressID)
	}

	if len(addressIDs) == 0 {
		return result
	}

	// 3. query address
	addressEntities, err := db.Address.
		Query().
		Where(address.IDIn(addressIDs...)).
		All(ctx)
	if err != nil {
		common.Error(err)
		return result
	}

	for _, addr := range addressEntities {
		result = append(result, &model.Address{
			AddressID:        int32(addr.ID),
			AddressType:      addr.AddressType,
			StreetAddress:    addr.StreetAddress,
			AptPO:            addr.AptPo,
			City:             addr.City,
			State:            addr.State,
			Zipcode:          addr.Zipcode,
			Country:          addr.Country,
			AddressConfirmed: addr.AddressConfirmed,
			IsPrimaryAddress: addr.IsPrimaryAddress,
		})
	}

	return result
}

func CreateAddress(ctx context.Context, db *ent.Client, address *model.Address) (addressID int32, err error) {
	createdAddress, err := db.Address.
		Create().
		SetAddressType(address.AddressType).
		SetStreetAddress(address.StreetAddress).
		SetAptPo(address.AptPO).
		SetCity(address.City).
		SetState(address.State).
		SetZipcode(address.Zipcode).
		SetCountry(address.Country).
		SetAddressConfirmed(address.AddressConfirmed).
		SetIsPrimaryAddress(address.IsPrimaryAddress).
		Save(ctx)
	if err != nil {
		return 0, err
	}
	return int32(createdAddress.ID), nil
}

func AddAddressToCustomerClinic(ctx context.Context, db *ent.Client, data *model.CustomerAddressOnClinicsCreation) error {
	err := db.CustomerAddressOnClinics.
		Create().
		SetAddressID(int(data.AddressID)).
		SetAddressType(data.AddressType).
		SetClinicID(int(data.ClinicID)).
		SetCustomerID(int(data.CustomerID)).
		Exec(ctx)
	return err
}
