package dbutils

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"

	"coresamples/ent"
	"coresamples/ent/labordersendhistory"
	"coresamples/ent/orderinfo"
	"coresamples/ent/sample"
	"coresamples/ent/sampleidgenerate"
	"coresamples/ent/testdetail"
	"coresamples/ent/tubereceive"
)

func GetdailyCollectionSamples(
	startTime time.Time,
	endTime time.Time,
	client *ent.Client,
	ctx context.Context,
) ([]*ent.Sample, error) {
	// Check for nil client first
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	// Query samples within the date range and preload relationships
	samples, err := client.Sample.Query().
		Where(
			sample.SampleReceivedTimeGTE(startTime),
			sample.SampleReceivedTimeLT(endTime),
		).
		Select(sample.FieldID).
		WithPatient(func(q *ent.PatientQuery) {
			q.WithPatientContacts()
		}).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to query samples: %w", err)
	}

	return samples, nil
}

// GetdailyCheckNonReceivedSamples retrieves non-received samples
func GetdailyCheckNonReceivedSamples(
	startTime time.Time,
	endTime time.Time,
	client *ent.Client,
	ctx context.Context,
) ([]*ent.Sample, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	samples, err := client.Sample.Query().
		Where(
			sample.SampleReceivedTimeIsNil(), // Filter samples where received time is NULL
			sample.HasOrderWith(
				orderinfo.OrderCreateTimeGTE(startTime),
				orderinfo.OrderCreateTimeLT(endTime),
			),
		).
		Select(sample.FieldID, sample.FieldAccessionID).
		WithPatient(func(q *ent.PatientQuery) {
			q.WithPatientContacts()
		}).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to query samples: %w", err)
	}

	return samples, nil
}

// GenerateSampleID generates a new sample ID
func GenerateSampleID(client *ent.Client, ctx context.Context) (*ent.SampleIDGenerate, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	return client.SampleIDGenerate.Create().Save(ctx)
}

// GetBarcodeForSampleID retrieves the barcode for a given sample ID
func GetBarcodeForSampleID(sampleId int, client *ent.Client, ctx context.Context) (string, error) {
	if client == nil {
		return "", fmt.Errorf("client is nil")
	}
	sampleIDGen, err := client.SampleIDGenerate.Query().
		Where(sampleidgenerate.IDEQ(sampleId)).
		First(ctx)
	if err != nil {
		return "", err
	}
	return sampleIDGen.Barcode, nil
}

func GetTestsBySampleId(sampleId int, client *ent.Client, ctx context.Context) ([]*ent.Test, error) {
	//TODO: redis?
	samp, err := client.Sample.Query().
		Where(sample.IDEQ(sampleId)).
		WithOrder(func(q *ent.OrderInfoQuery) {
			q.WithTests()
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	order, err := samp.Edges.OrderOrErr()
	if err != nil {
		return nil, err
	}
	return order.Edges.TestsOrErr()
}

func GetSampleWithDetailsBySampleId(sampleId int, details []string, client *ent.Client, ctx context.Context) (*ent.Sample, error) {
	//TODO: redis?
	return client.Sample.Query().
		Where(sample.IDEQ(sampleId)).
		WithOrder(func(q *ent.OrderInfoQuery) {
			q.WithTests(func(q *ent.TestQuery) {
				q.WithTestDetails(func(q *ent.TestDetailQuery) {
					q.Where(testdetail.TestDetailNameIn(details...))
				})
			})
		}).
		WithTubes().
		Only(ctx)
}

func GetSamplesByIds(sampleIds []int, client *ent.Client, ctx context.Context) ([]*ent.Sample, error) {
	return client.Sample.Query().Where(func(s *sql.Selector) {
		s.Where(sql.InInts(sample.FieldID, sampleIds...))
	}).WithOrder(func(q *ent.OrderInfoQuery) {
		q.WithOrderFlags()
	}).All(ctx)
}

func GetSampleById(sampleId int, client *ent.Client, ctx context.Context) (*ent.Sample, error) {
	return client.Sample.Query().Where(sample.IDEQ(sampleId)).First(ctx)
}

func CreateSample(input *ent.Sample, client *ent.Client, ctx context.Context) (*ent.Sample, error) {
	return client.Sample.Create().
		SetOrderID(input.OrderID).
		SetAccessionID(input.AccessionID).
		SetTubeCount(input.TubeCount).
		SetSampleDescription(input.SampleDescription).
		SetPatientID(input.PatientID).
		SetSampleCollectionTime(input.SampleCollectionTime).
		SetSampleReceivedTime(input.SampleReceivedTime).
		Save(ctx)
}

func GetSampleWithTubes(sampleId int, client *ent.Client, ctx context.Context) (*ent.Sample, error) {
	return client.Sample.Query().
		Where(sample.IDEQ(sampleId)).
		WithTubes(func(q *ent.TubeQuery) {
			q.WithTubeType()
		}).
		Only(ctx)
}

func GetSamplesWithAccessionIds(accessionIds []string, client *ent.Client, ctx context.Context) ([]*ent.Sample, error) {
	//TODO: add patient
	return client.Sample.
		Query().
		Where(sample.AccessionIDIn(accessionIds...)).
		WithOrder(func(q *ent.OrderInfoQuery) {
			q.WithOrderFlags()
		}).
		All(ctx)
}

func GetSampleReceiveRecordBySampleId(sampleId int, client *ent.Client, ctx context.Context) ([]*ent.TubeReceive, error) {
	return client.TubeReceive.
		Query().
		Where(tubereceive.SampleIDEQ(sampleId)).
		All(ctx)
}

func GetSampleReceiveRecordBySampleIds(sampleIds []int, client *ent.Client, ctx context.Context) ([]*ent.TubeReceive, error) {
	return client.TubeReceive.
		Query().
		Where(tubereceive.SampleIDIn(sampleIds...)).
		All(ctx)
}

func UpdateTubeReceiveRecord(record *ent.TubeReceive, client *ent.Client, ctx context.Context) error {
	update := client.TubeReceive.
		Update().
		Where(tubereceive.IDEQ(record.ID))
	if record.SampleID != 0 {
		update.SetSampleID(record.SampleID)
	}
	if record.TubeType != "" {
		update.SetTubeType(record.TubeType)
	}
	if !record.CollectionTime.IsZero() {
		update.SetCollectionTime(record.CollectionTime)
	}
	if record.ReceivedCount != 0 {
		update.SetReceivedCount(record.ReceivedCount)
	}
	if record.ReceivedBy != "" {
		update.SetReceivedBy(record.ReceivedBy)
	}
	if record.ModifiedBy != "" {
		update.SetModifiedBy(record.ModifiedBy)
	}
	if !record.ReceivedTime.IsZero() {
		update.SetReceivedTime(record.ReceivedTime)
	}
	_, err := update.Save(ctx)
	return err
}

func GetReceiveRecordById(id int, client *ent.Client, ctx context.Context) (*ent.TubeReceive, error) {
	return client.TubeReceive.Get(ctx, id)
}

func DeleteReceiveRecordById(id int, client *ent.Client, ctx context.Context) (*ent.TubeReceive, error) {
	tx, err := client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	txClient := tx.Client()
	record, err := GetReceiveRecordById(id, txClient, ctx)
	if err != nil {
		return nil, Rollback(tx, err)
	}
	err = txClient.TubeReceive.DeleteOne(record).Exec(ctx)
	if err != nil {
		return nil, Rollback(tx, err)
	}
	return record, tx.Commit()
}

func UpdateResendStatusBySampleIdAndTubeType(sampleId int, tubeType string, resendBlocked bool, client *ent.Client, ctx context.Context) error {
	return client.LabOrderSendHistory.
		Update().Where(labordersendhistory.SampleIDEQ(sampleId),
		labordersendhistory.TubeTypeEQ(tubeType)).
		SetIsResendBlocked(resendBlocked).Exec(ctx)
}

func GetMiniSampleByAccessionIds(accessionIds []string, client *ent.Client, ctx context.Context) ([]*ent.Sample, error) {
	return client.Sample.Query().
		Where(sql.FieldIn(sample.FieldAccessionID, accessionIds...)).
		Select(sample.FieldID, sample.FieldCustomerID, sample.FieldPatientID, sample.FieldAccessionID).
		All(ctx)
}

func GetSampleByAccessionId(accessionId string, client *ent.Client, ctx context.Context) (*ent.Sample, error) {
	return client.Sample.Query().Where(sample.AccessionIDEQ(accessionId)).First(ctx)
}

func UpdateSampleBarcode(sampleId int, barcode string, client *ent.Client, ctx context.Context) error {

	return client.SampleIDGenerate.Update().
		Where(sampleidgenerate.IDEQ(sampleId)).
		SetBarcode(barcode).
		Exec(ctx)
}

// CheckSampleIDExists verifies if a sample with the given ID exists in the database
func CheckSampleIDExists(sampleId int, client *ent.Client, ctx context.Context) (bool, error) {
	exists, err := client.SampleIDGenerate.Query().
		Where(sampleidgenerate.ID(sampleId)).
		Exist(ctx)

	if err != nil {
		return false, fmt.Errorf("failed to check sample existence: %w", err)
	}

	return exists, nil
}

func GetLastBarcodeInRange(startOfDay, endOfDay string, client *ent.Client, ctx context.Context) (*ent.SampleIDGenerate, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	return client.SampleIDGenerate.Query().
		Where(
			sampleidgenerate.And(
				sampleidgenerate.BarcodeGTE(startOfDay),
				sampleidgenerate.BarcodeLTE(endOfDay),
			),
		).
		Order(ent.Desc(sampleidgenerate.FieldBarcode)).
		First(ctx)
}
