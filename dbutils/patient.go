package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/patient"
	pb "coresamples/proto"
	"fmt"
	"time"
)

func GetGuestPatient(accessionId string, birthday string,
	firstName string, lastName string, client *ent.Client, ctx context.Context) (*ent.Patient, error) {
	// Step 1: Query the Sample table to find the patient_id using the accessionId.
	sample, err := GetSampleByAccessionId(accessionId, client, ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("no sample found for the given accession ID")
		}
		return nil, fmt.Errorf("failed to query sample: %w", err)
	}

	// Step 2: Use the patient_id from the Sample table to find the patient.
	patientRecord, err := client.Patient.Query().Where(
		patient.IDEQ(sample.PatientID),
		patient.PatientFirstNameEqualFold(firstName),
		patient.PatientLastNameEQ(lastName),
		patient.Or(
			patient.PatientBirthdateEQ(birthday),
			patient.PatientBirthdateEQ(""),
			patient.PatientBirthdateIsNil(),
		),
	).Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("patient not found for the given info")
		}
		return nil, fmt.Errorf("failed to query patient: %w", err)
	}

	// Return the patient record.
	return patientRecord, nil
}

func UpsertPatient(ctx context.Context, data *pb.PatientCDCUpdate_Data, client *ent.Client) error {
	if data == nil || client == nil {
		return fmt.Errorf("data:%v or ent.Client:%v is nil", data, client)
	}

	createTimeStr, err := time.Parse("2006-01-02 15:04:05", data.PaientCreateTime)
	if err != nil {
		return fmt.Errorf("parse time:%v failed, err:%s", data.PaientCreateTime, err.Error())
	}

	serviceDateStr, err := time.Parse("2006-01-02 15:04:05", data.PatientServiceDate)
	if err != nil {
		return fmt.Errorf("parse time:%v failed, err:%s", data.PatientServiceDate, err.Error())
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
	patientFlagged := data.PatientFlagged == 1

	// ✅ Step 2: Perform update using Ent ORM API
	err = tx.Patient.Create().SetID(int(data.PatientId)).
		SetUserID(int(data.UserId)).
		SetPatientType(data.PatientType).
		SetOriginalPatientID(int(data.OriginalPatientId)).
		SetPatientGender(data.PatientGender).
		SetPatientFirstName(data.PatientFirstName).
		SetPatientLastName(data.PatientLastName).
		SetPatientMiddleName(data.PatientMiddleName).
		SetPatientMedicalRecordNumber(data.PatientMedicalRecordNumber).
		SetPatientLegalFirstname(data.PatientLegalFirstname).
		SetPatientLegalLastname(data.PatientLegalLastname).
		SetPatientHonorific(data.PatientHonorific).
		SetPatientSuffix(data.PatientSuffix).
		SetPatientMarital(data.PatientMarital).
		SetPatientEthnicity(data.PatientEthnicity).
		SetPatientBirthdate(data.PatientBirthdate).
		SetPatientSsn(data.PatientSsn).
		SetPatientHeight(data.PatientHeight).
		SetPatientWeight(data.PatientWeight).
		SetOfficeallyID(int(data.OfficeallyId)).
		SetPatientNyWaiveFormIssueStatus(data.PatientNyWaiveFormIssueStatus).
		SetPatientCreateTime(createTimeStr).
		SetCustomerID(int(data.CustomerId)).
		SetIsActive(isActive).
		SetPatientFlagged(patientFlagged).
		SetPatientServiceDate(serviceDateStr).
		SetPatientDescription(data.PatientDescription).
		SetPatientLanguage(data.PatientLanguage).
		OnConflict().UpdateNewValues().Exec(ctx)
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to update patient (patientId: %d): %w", data.PatientId, err))
	}

	// ✅ Step 3: Re-enable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=1")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to re-enable foreign key checks: %w", err))
	}

	// ✅ Commit the transaction
	return tx.Commit()
}

func DeletePatient(ctx context.Context, patientID int32, client *ent.Client) error {
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
	err = tx.Clinic.DeleteOneID(int(patientID)).Exec(ctx)
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to delete patient(patient_id: %d): %w", patientID, err))
	}

	// ✅ Step 3: Re-enable foreign key checks
	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=1")
	if err != nil {
		return Rollback(tx, fmt.Errorf("failed to re-enable foreign key checks: %w", err))
	}

	// ✅ Commit the transaction
	return tx.Commit()
}
