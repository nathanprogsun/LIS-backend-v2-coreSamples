package service

import (
	"context"
	"coresamples/common"
	"coresamples/dbutils"
	"coresamples/ent"
	"coresamples/external"
	"coresamples/model"
	pb "coresamples/proto"
	"coresamples/tasks"
	"coresamples/util"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type ISampleService interface {
	ListSampleTests(sampleId int, ctx context.Context) ([]int32, error)
	ListSamples(sampleIds []int, ctx context.Context) ([]*ent.Sample, error)
	GetSampleTests(sampleIds []int32, ctx context.Context) ([]*pb.SampleTest, error)
	GetSampleTubes(sampleId int32, ctx context.Context) ([]*ent.Tube, error) // tube with tube types
	GetSampleTubesCount(sampleId int32, ctx context.Context) (map[string]int32, error)
	GetdailyCollectionSamples(startTime, endTime string, ctx context.Context) ([]*pb.SampleCollection, error)
	GetdailyCheckNonReceivedSamples(startTime, endTime string, ctx context.Context) ([]*pb.Sample_NonReceived, error)
	GetMultiSampleTubesCount(sampleIds []int32, ctx context.Context) ([]map[string]int32, error)
	ListSamplesByAccessionIDs(accessionIds []string, ctx context.Context) ([]*ent.Sample, error)
	GetSampleReceiveRecords(sampleId int, ctx context.Context) ([]*ent.TubeReceive, error)
	GetSampleReceiveRecordsBatch(sampleIds []int, ctx context.Context) (map[int][]*ent.TubeReceive, error)
	ModifySampleReceiveRecord(record *ent.TubeReceive, ctx context.Context) (*ent.TubeReceive, error)
	GetSampleTestsInstrument(sampleIds []int32, ctx context.Context) (map[int32][]string, error)
	GetSampleTypeViaTubeType(tubeType string, ctx context.Context) ([]*ent.TubeType, error)
	GetTubeTypeViaSampleTypeCode(sampleTypeCode string, ctx context.Context) (*ent.SampleType, error)
	GetTubeTypeViaSampleTypeEnum(sampleTypeEnum string, ctx context.Context) (*ent.SampleType, error)
	GetSampleMinimumInfoByAccessionIds(accessionIds []string, ctx context.Context) ([]*ent.Sample, error)
	UpdateSampleRequiredTubeCount(sampleId int32, tubeType string, requiredCnt int32, requiredBy string, ctx context.Context) error
	ReceiveSampleTubes(sampleId int32, tubeDetails []*model.SampleTubeDetails, receivedBy string, receivedTime time.Time, isRedraw bool, ctx context.Context) error
	GenerateSampleID(*pb.EmptyRequest, context.Context) (*pb.GenerateSampleIdResponse, error)
	GenerateBarcodeForSampleID(*pb.GenerateBarcodeForSampleIdRequest, context.Context) (*pb.GenerateBarcodeForSampleIdResponse, error)
}

type SampleTests struct {
	tubeIds []int32
	//details: test_instrument, test_sample_type, test_assay_name, test_group_name, test_duplicate_assay_name, test_turnaround_days
	testsWithDetails []*ent.Test
	sample           *ent.Sample
}

type SampleService struct {
	Service
	asynqClient tasks.AsynqClient
}

func NewSampleService(dbClient *ent.Client, redisClient *common.RedisClient, asynqClient tasks.AsynqClient) ISampleService {
	s := &SampleService{
		Service:     InitService(dbClient, redisClient),
		asynqClient: asynqClient,
	}
	return s
}

func (s *SampleService) ListSampleTests(sampleId int, ctx context.Context) ([]int32, error) {
	//sub_orders have been deprecated
	var ids []int32
	tests, err := dbutils.GetTestsBySampleId(sampleId, s.dbClient, ctx)
	if err != nil {
		return nil, err
	}
	for _, test := range tests {
		ids = append(ids, int32(test.ID))
	}
	return ids, nil
}

func (s *SampleService) ListSamples(sampleIds []int, ctx context.Context) ([]*ent.Sample, error) {
	return dbutils.GetSamplesByIds(sampleIds, s.dbClient, ctx)
}

//func (s *SampleService) CreateSample(input *ent.Sample) (*ent.Sample, error) {
//	//time.Parse(time.RFC3339)
//	return dbutils.CreateSample(input, s.dbClient, ctx)
//}

func (s *SampleService) GetSampleTests(sampleIds []int32, ctx context.Context) ([]*pb.SampleTest, error) {
	var ans []*pb.SampleTest
	for _, sampleId := range sampleIds {
		//TODO: cache this in redis
		// Construct the Redis key
		redisKey := dbutils.KeyGetSampleTests(int(sampleId))
		// Attempt to fetch from Redis
		redisResult, err := s.redisClient.Get(ctx, redisKey).Result()
		if err == nil {
			// Redis hit: Unmarshal and append the result
			var cachedSampleTest pb.SampleTest
			if unmarshalErr := json.Unmarshal([]byte(redisResult), &cachedSampleTest); unmarshalErr == nil {
				ans = append(ans, &cachedSampleTest)
				continue
			}
		}
		var sampleTest *pb.SampleTest

		var testAns []*pb.TestS
		var tubeIds []*pb.TubeID
		details := []string{
			"test_sample_type",
			"test_instrument",
			"test_assay_name",
			"test_group_name",
			"test_duplicate_asssay_name",
			"test_turnaround_days",
		}
		sample, err := dbutils.GetSampleWithDetailsBySampleId(int(sampleId), details, s.dbClient, ctx)
		if err != nil {
			//TODO: 153?
			return nil, err
		}
		//TODO: force fetch order from 153? if it's empty?
		order, err := sample.Edges.OrderOrErr()
		if err != nil {
			return nil, err
		}
		tubes, err := sample.Edges.TubesOrErr()
		if err != nil {
			return nil, err
		}
		tests, err := order.Edges.TestsOrErr()
		if err != nil {
			return nil, err
		}
		for _, test := range tests {
			var testInstrument string
			var testSampleType string
			var testAssayName string
			var testGroupName string
			var testDuplicateAssayName string
			var testTurnaroundDays string
			testDetails, err := test.Edges.TestDetailsOrErr()
			if err != nil {
				common.Error(err)
				continue
			}
			// fill out the values of interests
			for _, detail := range testDetails {
				detailName := detail.TestDetailName
				detailVal := detail.TestDetailsValue
				switch detailName {
				case "test_instrument":
					testInstrument = detailVal
				case "test_sample_type":
					testSampleType = detailVal
				case "test_assay_name":
					testAssayName = detailVal
				case "test_group_name":
					testGroupName = detailVal
				case "test_duplicate_asssay_name":
					testDuplicateAssayName = detailVal
				case "test_turnaround_days":
					testTurnaroundDays = detailVal
				}
			}
			sampletest := &pb.TestS{
				TestId:                 int32(test.ID),
				TestNames:              test.TestName,
				TestCodes:              test.TestCode,
				TestInstrument:         testInstrument,
				TestType:               testSampleType,
				TestAssayName:          testAssayName,
				TestGroupName:          testGroupName,
				TestDuplicateAssayName: testDuplicateAssayName,
				TestTurnaroundDays:     testTurnaroundDays,
			}
			testAns = append(testAns, sampletest)
		}
		for _, tube := range tubes {
			tubeIds = append(tubeIds, &pb.TubeID{TubeId: tube.TubeID})
		}

		sampleTest = &pb.SampleTest{
			Tests:                testAns,
			TubeIds:              tubeIds,
			SampleCollectionTime: sample.SampleCollectionTime.String(),
			PatientId:            int32(sample.PatientID),
			AccessionId:          sample.AccessionID,
			SampleId:             strconv.Itoa(int(sampleId)),
		}

		serialized, marshalErr := json.Marshal(sampleTest)
		if marshalErr == nil {
			s.redisClient.SetEX(ctx, redisKey, serialized, time.Hour)
		}
		ans = append(ans, sampleTest)
	}
	return ans, nil
}

// GetDailyCollectionSamples handles the gRPC request
func (s *SampleService) GetdailyCollectionSamples(startTime, endTime string, ctx context.Context) ([]*pb.SampleCollection, error) {
	start_Time, err := time.Parse("2006-01-02 15:04:05", startTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	end_Time, err := time.Parse("2006-01-02 15:04:05", endTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	samples, err := dbutils.GetdailyCollectionSamples(start_Time, end_Time, s.dbClient, ctx)
	if err != nil {
		return nil, err
	}

	var sampleCollections []*pb.SampleCollection

	for idx, sp := range samples {
		sc := &pb.SampleCollection{
			SampleId: strconv.Itoa(sp.ID),
		}

		err = util.Swap(samples[idx].Edges.Patient, &sc.Patient)

		if err != nil {
			sentry.CaptureMessage(err.Error())
			return nil, err
		}

		err = util.Swap(samples[idx].Edges.Patient.Edges.PatientContacts, &sc.Patient.PatientContact)

		if err != nil {
			sentry.CaptureMessage(err.Error())
			return nil, err
		}

		sampleCollections = append(sampleCollections, sc)
	}

	return sampleCollections, nil
}

func (s *SampleService) GetdailyCheckNonReceivedSamples(startTime, endTime string, ctx context.Context) ([]*pb.Sample_NonReceived, error) {
	// Parse input time strings
	start_Time, err := time.Parse("2006-01-02 15:04:05", startTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	end_Time, err := time.Parse("2006-01-02 15:04:05", endTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	// Get samples from database
	samples, err := dbutils.GetdailyCheckNonReceivedSamples(start_Time, end_Time, s.dbClient, ctx)
	if err != nil {
		return nil, err
	}

	var sampleNonRecCollections []*pb.Sample_NonReceived

	// Process each sample
	for idx, sp := range samples {
		sampleNonRec := &pb.Sample_NonReceived{
			SampleId:    strconv.Itoa(sp.ID),
			AccessionId: sp.AccessionID,
		}

		err = util.Swap(samples[idx].Edges.Patient, &sampleNonRec.Patient)

		if err != nil {
			sentry.CaptureMessage(err.Error())
			return nil, err
		}

		err = util.Swap(samples[idx].Edges.Patient.Edges.PatientContacts, &sampleNonRec.Patient.PatientContact)

		if err != nil {
			sentry.CaptureMessage(err.Error())
			return nil, err
		}

		sampleNonRecCollections = append(sampleNonRecCollections, sampleNonRec)
	}

	return sampleNonRecCollections, nil
}

func (s *SampleService) GenerateSampleID(_ *pb.EmptyRequest, ctx context.Context) (*pb.GenerateSampleIdResponse, error) {
	// Create a new row in sample_id_generate table to get an auto-incremented ID
	sampleIDGen, err := dbutils.GenerateSampleID(s.dbClient, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate sample ID: %w", err)
	}

	return &pb.GenerateSampleIdResponse{
		SampleId: int32(sampleIDGen.ID),
	}, nil
}

func (s *SampleService) generateAndSaveBarcode(datePrefix string, sampleId int, ctx context.Context) (string, error) {
	if s.dbClient == nil {
		return "", fmt.Errorf("client is nil")
	}

	// Check if the sample ID exists
	exists, err := dbutils.CheckSampleIDExists(sampleId, s.dbClient, ctx)

	if err != nil {
		return "", fmt.Errorf("failed to check sample ID: %w", err)
	}
	if !exists {
		return "", fmt.Errorf("invalid sample ID: %d", sampleId)
	}

	startOfDay := datePrefix + "6001"
	endOfDay := datePrefix + "9999"

	lastBarcode, err := dbutils.GetLastBarcodeInRange(startOfDay, endOfDay, s.dbClient, ctx)

	suffix := 6001
	if err == nil && lastBarcode != nil {
		numericSuffix, err := strconv.Atoi(lastBarcode.Barcode[len(datePrefix):])
		if err != nil {
			return "", fmt.Errorf("failed to parse last barcode suffix: %w", err)
		}
		suffix = numericSuffix + 1
	}

	maxRetries := 100 // Prevent infinite loops
	retryCount := 0

	for retryCount < maxRetries {
		var barcodeSuffix string
		if suffix <= 9999 {
			barcodeSuffix = fmt.Sprintf("%04d", suffix)
		} else {
			alphaIndex := (suffix - 10000) / 999
			numberPart := ((suffix - 10000) % 999) + 1
			if (suffix - 9999) > 26*999 {
				return "", fmt.Errorf("maximum barcodes for the day reached")
			}
			barcodeSuffix = fmt.Sprintf("%c%03d", 'A'+alphaIndex, numberPart)
		}

		barcode := datePrefix + barcodeSuffix

		// Try to update the record with the new barcode
		err := dbutils.UpdateSampleBarcode(sampleId, barcode, s.dbClient, ctx)

		if err == nil {
			return barcode, nil
		}

		if ent.IsConstraintError(err) {
			suffix++
			retryCount++
			continue
		}
		return "", err
	}

	return "", fmt.Errorf("maximum barcodes for the day reached")
}

func (s *SampleService) GenerateBarcodeForSampleID(req *pb.GenerateBarcodeForSampleIdRequest, ctx context.Context) (*pb.GenerateBarcodeForSampleIdResponse, error) {
	// Check if barcode already exists
	existingBarcode, err := dbutils.GetBarcodeForSampleID(int(req.SampleId), s.dbClient, ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("failed to check existing barcode: %w", err)
	}
	if existingBarcode != "" {
		return &pb.GenerateBarcodeForSampleIdResponse{
			Barcode: existingBarcode,
		}, nil
	}

	// Get current time in PT timezone
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return nil, fmt.Errorf("failed to load timezone: %w", err)
	}
	now := time.Now().In(loc)

	// Format date components
	datePrefix := fmt.Sprintf("%02d%02d%02d", now.Year()%100, now.Month(), now.Day())

	// Generate and save new barcode
	barcode, err := s.generateAndSaveBarcode(datePrefix, int(req.SampleId), ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate barcode: %w", err)
	}

	return &pb.GenerateBarcodeForSampleIdResponse{
		Barcode: barcode,
	}, nil
}

func (s *SampleService) GetSampleTubes(sampleId int32, ctx context.Context) ([]*ent.Tube, error) {
	sample, err := dbutils.GetSampleWithTubes(int(sampleId), s.dbClient, ctx)
	if err != nil {
		return nil, err
	}
	return sample.Edges.TubesOrErr()
}

func (s *SampleService) GetSampleTubesCount(sampleId int32, ctx context.Context) (map[string]int32, error) {
	requiredTubes, err := dbutils.GetSampleRequiredTubes(int(sampleId), s.dbClient, ctx)
	if err != nil {
		return nil, err
	}
	if len(requiredTubes) == 0 {
		resp, err := external.GetOrderService().GetRequiredTubesBySampleId(int(sampleId))
		if err != nil {
			return nil, err
		}
		for tube, cnt := range resp.NumberOfTubes {
			err = s.UpdateSampleRequiredTubeCount(sampleId, tube, cnt, "order_service", ctx)
			if err != nil {
				return nil, err
			}
		}
		return resp.NumberOfTubes, nil
	}
	numberOfTubes := make(map[string]int32)
	for _, tube := range requiredTubes {
		numberOfTubes[tube.TubeType] = int32(tube.RequiredCount)
	}
	return numberOfTubes, nil
}

func (s *SampleService) GetMultiSampleTubesCount(sampleIds []int32, ctx context.Context) ([]map[string]int32, error) {
	// add redis for this
	var res []map[string]int32
	for _, sampleId := range sampleIds {
		redisKey := dbutils.KeySampleTubeCounts(int(sampleId))

		redisResult, err := s.redisClient.Get(ctx, redisKey).Result()

		if err == nil {
			// Redis hit: Parse the cached result
			var cachedData map[string]int32
			if unmarshalErr := json.Unmarshal([]byte(redisResult), &cachedData); unmarshalErr == nil {
				res = append(res, cachedData)
				continue
			}
		}
		noOfTubes, err := s.GetSampleTubesCount(sampleId, ctx)
		if err != nil {
			return nil, err
		}
		// Serialize and cache the data in Redis
		serialized, marshalErr := json.Marshal(noOfTubes)
		if marshalErr == nil {
			s.redisClient.SetEX(ctx, redisKey, serialized, time.Minute*10)
		}
		res = append(res, noOfTubes)
	}
	return res, nil
}

func (s *SampleService) ListSamplesByAccessionIDs(accessionIds []string, ctx context.Context) ([]*ent.Sample, error) {
	// what is accession id?
	return dbutils.GetSamplesWithAccessionIds(accessionIds, s.dbClient, ctx)
}

func (s *SampleService) GetSampleReceiveRecords(sampleId int, ctx context.Context) ([]*ent.TubeReceive, error) {
	records, err := dbutils.GetSampleReceiveRecordBySampleId(sampleId, s.dbClient, ctx)
	if ent.IsNotFound(err) {
		common.Infof("Get empty sample received record with sample id %d", sampleId)
		return records, nil
	}
	return records, err
}

func (s *SampleService) GetSampleReceiveRecordsBatch(sampleIds []int, ctx context.Context) (map[int][]*ent.TubeReceive, error) {
	ans := map[int][]*ent.TubeReceive{}
	records, err := dbutils.GetSampleReceiveRecordBySampleIds(sampleIds, s.dbClient, ctx)
	if records == nil || ent.IsNotFound(err) {
		common.Infof("Get empty sample received record with sample id %v", sampleIds)
	}
	for _, record := range records {
		if _, ok := ans[record.SampleID]; !ok {
			ans[record.SampleID] = []*ent.TubeReceive{}
		}
		ans[record.SampleID] = append(ans[record.SampleID], record)
	}
	return ans, nil
}

func (s *SampleService) ModifySampleReceiveRecord(record *ent.TubeReceive, ctx context.Context) (*ent.TubeReceive, error) {
	//TODO: audit log?
	if record.ReceivedCount != 0 {
		err := dbutils.UpdateTubeReceiveRecord(record, s.dbClient, ctx)
		if err != nil {
			return nil, err
		}
		return dbutils.GetReceiveRecordById(record.ID, s.dbClient, ctx)
	} else {
		ret, err := dbutils.DeleteReceiveRecordById(record.ID, s.dbClient, ctx)
		if err != nil {
			return nil, err
		}
		err = dbutils.UpdateResendStatusBySampleIdAndTubeType(record.SampleID, record.TubeType, false, s.dbClient, ctx)
		//delete lab send order from redis cache
		s.redisClient.Del(ctx, dbutils.KeyLabOrderSendRecord(record.SampleID, record.TubeType))
		return ret, nil
	}
}

func (s *SampleService) GetSampleTestsInstrument(sampleIds []int32, ctx context.Context) (map[int32][]string, error) {
	var ret map[int32][]string
	for _, sampleId := range sampleIds {
		var instruments []string
		sample, err := dbutils.GetSampleWithDetailsBySampleId(int(sampleId), []string{"test_instrument"}, s.dbClient, ctx)
		if err != nil {
			//TODO: create this from old db? prisma4
			return nil, err
		}
		order, err := sample.Edges.OrderOrErr()
		if err != nil {
			return nil, err
		}
		tests, err := order.Edges.TestsOrErr()
		if err != nil {
			return nil, err
		}
		for _, test := range tests {
			testDetails, err := test.Edges.TestDetailsOrErr()
			if err != nil {
				common.Error(err)
				continue
			}
			for _, detail := range testDetails {
				if detail.TestDetailName == "test_instrument" {
					instruments = append(instruments, detail.TestDetailsValue)
				}
			}
		}
		ret[sampleId] = instruments
	}
	return ret, nil
}

func (s *SampleService) GetSampleTypeViaTubeType(tubeType string, ctx context.Context) ([]*ent.TubeType, error) {
	return dbutils.GetTubeTypeInfoByTubeTypeEnum(tubeType, s.dbClient, ctx)
}

func (s *SampleService) GetTubeTypeViaSampleTypeCode(sampleTypeCode string, ctx context.Context) (*ent.SampleType, error) {
	//TODO: add redis here
	// Construct the Redis key
	redisKey := dbutils.KeyTubeTypeViaSampleTypeCode(sampleTypeCode)
	// Attempt to fetch from Redis
	redisResult, err := s.redisClient.Get(ctx, redisKey).Result()
	if err == nil {
		var sampleType ent.SampleType
		if unmarshalErr := json.Unmarshal([]byte(redisResult), &sampleType); unmarshalErr == nil {
			return &sampleType, nil
		}
	}
	// Cache miss: Fetch from the database
	sampleType, err := dbutils.GetTubeTypeInfoBySampleTypeCode(sampleTypeCode, s.dbClient, ctx)
	if err != nil {
		return nil, err
	}
	serialized, marshalErr := json.Marshal(sampleType)
	if marshalErr == nil {
		s.redisClient.SetEX(ctx, redisKey, serialized, time.Minute*10)
	}
	return sampleType, err
}

func (s *SampleService) GetTubeTypeViaSampleTypeEnum(sampleTypeEnum string, ctx context.Context) (*ent.SampleType, error) {
	return dbutils.GetTubeTypeBySampleTypeEnum(sampleTypeEnum, s.dbClient, ctx)
}

func (s *SampleService) GetSampleMinimumInfoByAccessionIds(accessionIds []string, ctx context.Context) ([]*ent.Sample, error) {
	return dbutils.GetMiniSampleByAccessionIds(accessionIds, s.dbClient, ctx)
}

func (s *SampleService) UpdateSampleRequiredTubeCount(sampleId int32, tubeType string, requiredCnt int32, requiredBy string, ctx context.Context) error {
	tubeReq, err := dbutils.FindTubeRequirement(int(sampleId), tubeType, s.dbClient, ctx)
	if err == nil && tubeReq != nil {
		tubeReq, err = tubeReq.Update().
			SetRequiredCount(int(requiredCnt)).
			SetRequiredBy(requiredBy).
			SetModifiedBy(requiredBy).Save(ctx)
		if err != nil {
			return err
		}
	} else {
		err = dbutils.CreateTubeRequirement(int(sampleId), tubeType, requiredCnt, requiredBy, s.dbClient, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SampleService) ReceiveSampleTubes(sampleId int32, tubeDetails []*model.SampleTubeDetails, receivedBy string, receivedTime time.Time, isRedraw bool, ctx context.Context) error {
	trackingId := uuid.NewString()
	tx, err := s.dbClient.Tx(ctx)
	if err != nil {
		return err
	}
	for _, detail := range tubeDetails {
		if detail.ReceiveCount < 0 {
			return fmt.Errorf("received tube count of negative number")
		}
		_, err = dbutils.CreateTubeReceive(sampleId, detail.TubeType, detail.CollectionTime, detail.ReceiveCount, receivedBy, receivedTime, isRedraw, tx.Client(), ctx)
		if err != nil {
			return dbutils.Rollback(tx, err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return dbutils.Rollback(tx, err)
	}
	flagOrderTask, err := tasks.NewFlagOrderOnReceivingTask(&tasks.SampleTubeReceiveTask{SampleId: sampleId})
	taskInfo, err := s.asynqClient.Enqueue(flagOrderTask,
		asynq.Queue(common.Env.AsynqQueueName),
		asynq.MaxRetry(10))

	if err != nil {
		common.Error(err)
		sentry.CaptureException(err)
		return err
	}
	common.LogTaskInfo(taskInfo)

	sendOrderTask, err := tasks.NewSendOrderOnReceivingTask(&tasks.SampleTubeReceiveTask{
		SampleId:     sampleId,
		TubeDetails:  tubeDetails,
		ReceivedTime: receivedTime,
		IsRedraw:     isRedraw,
	})
	taskInfo, err = s.asynqClient.Enqueue(sendOrderTask,
		asynq.Queue(common.Env.AsynqQueueName),
		asynq.MaxRetry(100))

	if err != nil {
		common.Error(err)
		sentry.CaptureException(err)
		return err
	}
	common.LogTaskInfo(taskInfo)

	//send sample receive general event
	var details []*pb.TubeDetail
	for _, detail := range tubeDetails {
		details = append(details, &pb.TubeDetail{
			TubeType:       detail.TubeType,
			CollectionTime: detail.CollectionTime.Format(time.RFC3339),
			ReceiveCount:   detail.ReceiveCount,
		})
	}
	event := &pb.GeneralEvent{
		SampleId: sampleId,
		AddonColumn: &pb.EventAddonColumn{
			TubeDetails:  details,
			ReceivedBy:   receivedBy,
			ReceivedTime: receivedTime.Format(time.RFC3339),
			IsRedraw:     isRedraw,
		},
		EventId:       trackingId,
		EventProvider: "lis-shipping",
		EventName:     "receive_sample_tubes",
	}
	task, err := tasks.NewSendSampleReceiveGeneralEventTask(event)
	if err != nil {
		return err
	}
	taskInfo, err = s.asynqClient.Enqueue(task,
		asynq.Queue(common.Env.AsynqQueueName),
		asynq.MaxRetry(10))

	if err != nil {
		common.Error(err)
		sentry.CaptureException(err)
		return err
	}
	common.LogTaskInfo(taskInfo)
	return nil
}
