package dbutils

import (
	"context"
	"coresamples/common"
	"coresamples/ent"
	"coresamples/ent/customer"
	"coresamples/ent/patient"
	"coresamples/model"
	pb "coresamples/proto"
	"coresamples/util"
	"encoding/json"
	"fmt"
	"time"
)

func GetCustomerByCustomerID(customerId int, client *ent.Client, ctx context.Context) (*ent.Customer, error) {
	return client.Customer.Get(ctx, customerId)
}

func UpsertCustomerCDC(data *pb.CustomerCDCUpdate_Data, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	reqSubmitTime, err := time.Parse("2006-01-02 15:04:05", data.CustomerRequestSubmitTime)
	if err != nil {
		return err
	}
	signUpTime, err := time.Parse("2006-01-02 15:04:05", data.CustomerSignupTime)
	if err != nil {
		return err
	}
	fillTime, err := time.Parse("2006-01-02 15:04:05", data.OnboardingQuestionnaireFilledOn)
	if err != nil {
		return err
	}
	if data.CustomerId == 0 {
		return fmt.Errorf("cdc customer message doesn't have customer id %v", data)
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}

	isActive := data.IsActive == 1
	orderPlacementAllowed := data.OrderPlacementAllowed == 1
	betaProgramEnabled := data.BetaProgramEnabled == 1

	err = tx.Customer.Create().SetID(int(data.CustomerId)).
		SetUserID(int(data.UserId)).
		SetCustomerType(data.CustomerType).
		SetCustomerFirstName(data.CustomerFirstName).
		SetCustomerLastName(data.CustomerLastName).
		SetCustomerMiddleName(data.CustomerMiddleName).
		SetCustomerTypeID(data.CustomerTypeId).
		SetCustomerSuffix(data.CustomerSuffix).
		SetCustomerSamplesReceived(data.CustomerSamplesReceived).
		SetCustomerRequestSubmitTime(reqSubmitTime).
		SetCustomerSignupTime(signUpTime).
		SetIsActive(isActive).
		SetSalesID(int(data.SalesId)).
		SetCustomerNpiNumber(data.CustomerNpiNumber).
		SetReferralSource(data.ReferralSource).
		SetOrderPlacementAllowed(orderPlacementAllowed).
		SetBetaProgramEnabled(betaProgramEnabled).
		SetOnboardingQuestionnaireFilledOn(fillTime).
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

func DeleteCustomerCDC(data *pb.CustomerCDCUpdate_Data, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	if data.CustomerId == 0 {
		return fmt.Errorf("cdc customer message doesn't have customer id %v", data)
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}

	_, err = client.Customer.Delete().Where(customer.ID(int(data.CustomerId))).Exec(ctx)
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

func UpdatePatientToCustomerCDC(data *pb.CustomerToPatientCDCUpdate_Data, old *pb.CustomerToPatientCDCUpdate_Data, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}

	//old.A, aka customer id is not zero, meaning patient's customer has changed
	err = tx.Patient.Update().
		Where(patient.IDEQ(int(data.B))).RemovePatientCustomerIDs(int(old.A)).Exec(ctx)
	if err != nil {
		return Rollback(tx, err)
	}

	err = tx.Patient.Update().
		Where(patient.IDEQ(int(data.B))).AddPatientCustomerIDs(int(data.A)).Exec(ctx)
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

func UpdateCustomerToPatientCDC(data *pb.CustomerToPatientCDCUpdate_Data, old *pb.CustomerToPatientCDCUpdate_Data, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}
	//old.B, aka patient id is not zero, meaning customer's patient has changed
	err = tx.Customer.Update().
		Where(customer.IDEQ(int(data.A))).RemovePatientIDs(int(old.B)).Exec(ctx)
	if err != nil {
		return Rollback(tx, err)
	}

	err = tx.Customer.Update().
		Where(customer.IDEQ(int(data.A))).AddPatientIDs(int(data.B)).Exec(ctx)
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

func CreateCustomerToPatientCDC(data *pb.CustomerToPatientCDCUpdate_Data, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}

	err = tx.Customer.Update().
		Where(customer.IDEQ(int(data.A))).
		AddPatientIDs(int(data.B)).Exec(ctx)
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

func DeleteCustomerToPatientCDC(data *pb.CustomerToPatientCDCUpdate_Data, client *ent.Client, ctx context.Context) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		return Rollback(tx, err)
	}

	err = tx.Customer.Update().
		Where(customer.IDEQ(int(data.A))).RemovePatientIDs(int(data.B)).Exec(ctx)
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

func FetchModelFullCustomers(ctx context.Context, db *ent.Client, customerID *int, offset, limit int) ([]*model.FullCustomer, error) {
	var result []*model.FullCustomer

	query := db.Customer.Query().WithClinics()

	if customerID != nil {
		query = query.Where(customer.IDEQ(*customerID))
	} else {
		query = query.Offset(offset).Limit(limit)
	}

	customers, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	for _, cust := range customers {
		var clinics []*model.CustomerClinic

		for _, clinic := range cust.Edges.Clinics {
			addresses := FetchCustomerClinicAddresses(ctx, db, cust.ID, clinic.ID)
			contacts := FetchCustomerClinicContacts(ctx, db, cust.ID, clinic.ID)

			clinics = append(clinics, &model.CustomerClinic{
				ClinicID:          int32(clinic.ID),
				ClinicName:        clinic.ClinicName,
				UserID:            int32(clinic.UserID),
				IsActive:          clinic.IsActive,
				ClinicAccountID:   int32(clinic.ClinicAccountID),
				CustomerAddresses: addresses,
				CustomerContacts:  contacts,
			})
		}

		result = append(result, &model.FullCustomer{
			CustomerID:                int32(cust.ID),
			UserID:                    int32(cust.UserID),
			CustomerFirstName:         cust.CustomerFirstName,
			CustomerLastName:          cust.CustomerLastName,
			CustomerMiddleName:        cust.CustomerMiddleName,
			CustomerTypeID:            cust.CustomerTypeID,
			CustomerSuffix:            cust.CustomerSuffix,
			CustomerSamplesReceived:   cust.CustomerSamplesReceived,
			CustomerRequestSubmitTime: cust.CustomerRequestSubmitTime.Format(time.RFC3339),
			IsActive:                  cust.IsActive,
			Clinics:                   clinics,
			CustomerNPINumber:         cust.CustomerNpiNumber,
			SalesID:                   int32(cust.SalesID),
			CustomerSignupTime:        cust.CustomerSignupTime.Format(time.RFC3339),
		})
	}

	return result, nil
}

var FetchAndCacheCustomerClinicData = func(ctx context.Context, db *ent.Client, redisClient *common.RedisClient) ([]*model.CustomerClinicData, error) {
	customers, err := db.Customer.
		Query().
		WithClinics().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}

	var data []*model.CustomerClinicData

	for _, cust := range customers {
		fullName := util.AssembleFullName(cust.CustomerFirstName, "", cust.CustomerLastName)
		for _, clinic := range cust.Edges.Clinics {
			data = append(data, &model.CustomerClinicData{
				CustomerId:   int32(cust.ID),
				CustomerName: fullName,
				ClinicName:   clinic.ClinicName,
			})
		}
	}
	// Marshal data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %w", err)
	}

	if err := redisClient.Set(ctx, "customer_clinic_data", jsonData, 24*time.Hour).Err(); err != nil {
		return nil, fmt.Errorf("redis SET error: %w", err)
	}

	return data, nil
}

func GetNewCustomerID(ctx context.Context, db *ent.Client) (customerID int, err error) {
	existingCust, err := db.Customer.
		Query().
		Where(
			customer.IDGT(38000),
			customer.IDLT(99999),
		).
		Order(ent.Desc(customer.FieldID)).
		First(ctx)
	newCustomerID := 38001
	if err == nil && existingCust != nil {
		newCustomerID = existingCust.ID + 1
	}

	return newCustomerID, nil
}
