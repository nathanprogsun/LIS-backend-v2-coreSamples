package dbutils

import (
	"context"
	"coresamples/common"
	"coresamples/ent"
	"coresamples/ent/contact"
	"coresamples/ent/customercontactonclinics"
	"coresamples/model"
	pb "coresamples/proto"
	"fmt"
)

func UpsertContactCDC(data *pb.ContactCDCUpdate_ContactData, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	if data.ContactId == 0 {
		return fmt.Errorf("cdc contact message doesn't have contact id %v", data)
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}
	isPrimaryContact := data.IsPrimaryContact == 1
	is2FaContact := data.Is_2FaContact == 1
	applyToAllGroupMember := data.ApplyToAllGroupMember == 1
	isGroupContact := data.IsGroupContact == 1
	useAsDefaultCreateContact := data.UseAsDefaultCreateContact == 1
	useGroupContact := data.UseGroupContact == 1

	err = tx.Contact.Create().
		SetID(int(data.ContactId)).
		SetContactDescription(data.ContactDescription).
		SetContactDetails(data.ContactDetails).
		SetContactType(data.ContactType).
		SetIsPrimaryContact(isPrimaryContact).
		SetIs2faContact(is2FaContact).
		SetCustomerID(int(data.CustomerId)).
		SetPatientID(int(data.PatientId)).
		SetClinicID(int(data.ClinicId)).
		SetInternalUserID(int(data.InternalUserId)).
		SetUserID(int(data.UserId)).
		SetContactLevel(int(data.ContactLevel)).
		SetContactLevelName(data.ContactLevelName).
		SetGroupContactID(int(data.GroupContactId)).
		SetApplyToAllGroupMember(applyToAllGroupMember).
		SetIsGroupContact(isGroupContact).
		SetUseAsDefaultCreateContact(useAsDefaultCreateContact).
		SetUseGroupContact(useGroupContact).
		OnConflict().UpdateNewValues().Exec(ctx)
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

func DeleteContactCDC(data *pb.ContactCDCUpdate_ContactData, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	if data.ContactId == 0 {
		return fmt.Errorf("cdc contact message doesn't have contact id %v", data)
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}

	_, err = client.Contact.Delete().Where(contact.ID(int(data.ContactId))).Exec(ctx)
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

func FetchCustomerClinicContacts(ctx context.Context, db *ent.Client, customerID, clinicID int) []*model.Contact {
	var result []*model.Contact

	relations, err := db.CustomerContactOnClinics.
		Query().
		Where(
			customercontactonclinics.CustomerIDEQ(customerID),
			customercontactonclinics.ClinicIDEQ(clinicID),
		).
		All(ctx)
	if err != nil {
		common.Error(err)
		return result
	}

	var contactIDs []int
	for _, rel := range relations {
		contactIDs = append(contactIDs, rel.ContactID)
	}

	if len(contactIDs) == 0 {
		return result
	}

	contactEntities, err := db.Contact.
		Query().
		Where(contact.IDIn(contactIDs...)).
		All(ctx)
	if err != nil {
		common.Error(err)
		return result
	}

	for _, ct := range contactEntities {
		result = append(result, &model.Contact{
			ContactID:          int32(ct.ID),
			ContactDescription: ct.ContactDescription,
			ContactDetails:     ct.ContactDetails,
			ContactType:        ct.ContactType,
			IsPrimaryContact:   ct.IsPrimaryContact,
		})
	}

	return result
}

func CreateContact(ctx context.Context, db *ent.Client, contact *model.Contact) (contactID int32, err error) {
	createdContact, err := db.Contact.
		Create().
		SetContactDescription(contact.ContactDescription).
		SetContactDetails(contact.ContactDetails).
		SetContactType(contact.ContactType).
		SetIsPrimaryContact(contact.IsPrimaryContact).
		Save(ctx)
	if err != nil {
		return 0, err
	}
	return int32(createdContact.ID), nil
}

func AddContactToCustomerClinic(ctx context.Context, db *ent.Client, data *model.CustomerContactOnClinicsCreation) error {
	err := db.CustomerContactOnClinics.
		Create().
		SetContactID(int(data.ContactID)).
		SetContactType(data.ContactType).
		SetClinicID(int(data.ClinicID)).
		SetCustomerID(int(data.CustomerID)).
		Exec(ctx)
	return err
}
