package service

import (
	"context"
	"coresamples/common"
	"coresamples/ent/enttest"
	pb "coresamples/proto"
	"coresamples/publisher"
	"coresamples/tasks"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stvp/tempredis"
)

type mockAsynqClient struct {
	mock.Mock
	tasks.AsynqClient
}

func setupSampleTest(t *testing.T) (*SampleService, *tempredis.Server, context.Context) {
	dataSource := "file:ent?mode=memory&_fk=1"
	dbClient := enttest.Open(t, "sqlite3", dataSource)
	ctx := context.Background()
	err := dbClient.Schema.Create(ctx)
	if err != nil {
		common.Fatalf("failed opening connection to MySQL", err)
	}
	server, err := tempredis.Start(tempredis.Config{
		"port": "0",
	})
	if err != nil {
		common.Fatalf("Failed to start tempredis: %v", err)
	}

	common.InitZapLogger("debug")

	redisClient := redis.NewClient(&redis.Options{
		Network: "unix",
		Addr:    server.Socket(),
	})

	mockAsynq := &mockAsynqClient{}

	s := &SampleService{
		Service: InitService(
			dbClient,
			common.NewRedisClient(redisClient, redisClient)),
		asynqClient: mockAsynq,
	}
	// This will start a mock publisher in the memory
	publisher.InitMockPublisher()

	return s, server, ctx
}

func cleanUpSampleTest(svc *SampleService, s *tempredis.Server) {
	var err error
	if err = s.Kill(); err != nil {
		common.Error(err)
	}
	if svc.dbClient != nil {
		if err = svc.dbClient.Close(); err != nil {
			common.Error(err)
		}
	}

	publisher.GetPublisher().GetWriter().Close()
}

func TestGetTubeTypeViaSampleTypeEnum(t *testing.T) {
	svc, redisServer, ctx := setupSampleTest(t)
	defer cleanUpSampleTest(svc, redisServer)

	sampleTypeEnum := "testSampleTypeEnum"
	sampleTypeName := "testSampleTypeName"
	sampleTypeCode := "testSampleTypeCode"
	sampleTypeDescription := "testSampleTypeDescription"
	primarySampleTypeGroup := "testPrimarySampleTypeGroup"

	// Seed the database
	sampleType1, err := svc.dbClient.SampleType.Create().
		SetSampleTypeName(sampleTypeName).
		SetSampleTypeEnum(sampleTypeEnum).
		SetSampleTypeCode(sampleTypeCode).
		SetSampleTypeDescription(sampleTypeDescription).
		SetPrimarySampleTypeGroup(primarySampleTypeGroup).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to seed sample type: %v", err)
	}

	_, err = svc.dbClient.SampleType.Create().
		SetSampleTypeName("testSampleTypeName2").
		SetSampleTypeEnum(sampleTypeEnum).
		SetSampleTypeCode(sampleTypeCode).
		SetSampleTypeDescription(sampleTypeDescription).
		SetPrimarySampleTypeGroup(primarySampleTypeGroup).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to seed sample type: %v", err)
	}

	tubeTypes := []struct {
		TubeName       string
		TubeTypeEnum   string
		TubeTypeSymbol string
	}{
		{"testTubeA", "testTubeEnumA", "testTubeSymbolA"},
		{"testTubeB", "testTubeEnumB", "testTubeSymbolB"},
	}

	for _, tube := range tubeTypes {
		_, err := svc.dbClient.TubeType.Create().
			SetTubeName(tube.TubeName).
			SetTubeTypeEnum(tube.TubeTypeEnum).
			SetTubeTypeSymbol(tube.TubeTypeSymbol).
			AddSampleTypeIDs(sampleType1.ID).
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to seed tube type: %v", err)
		}
	}

	// Run the service method
	actualSampleType, err := svc.GetTubeTypeViaSampleTypeEnum(sampleTypeEnum, ctx)
	if err != nil {
		t.Fatalf("GetTubeTypeViaSampleTypeEnum failed: %v", err)
	}

	// Verify main fields
	if actualSampleType.SampleTypeEnum != sampleTypeEnum ||
		actualSampleType.SampleTypeName != sampleTypeName ||
		actualSampleType.SampleTypeCode != sampleTypeCode ||
		actualSampleType.SampleTypeDescription != sampleTypeDescription ||
		actualSampleType.PrimarySampleTypeGroup != primarySampleTypeGroup {
		t.Fatal("sample type fields do not match expected values")
	}

	// Verify TubeTypes
	if len(actualSampleType.Edges.TubeTypes) != len(tubeTypes) {
		t.Fatalf("expected %d TubeTypes, got: %d", len(tubeTypes), len(actualSampleType.Edges.TubeTypes))
	}

	for i, actualTube := range actualSampleType.Edges.TubeTypes {
		expectedTube := tubeTypes[i]
		if actualTube.TubeName != expectedTube.TubeName ||
			actualTube.TubeTypeEnum != expectedTube.TubeTypeEnum ||
			actualTube.TubeTypeSymbol != expectedTube.TubeTypeSymbol {
			t.Fatalf("TubeType mismatch: expected %+v, got %+v", expectedTube, actualTube)
		}
	}
}

func TestGetTubeTypeViaSampleTypeCode(t *testing.T) {
	svc, redisServer, ctx := setupSampleTest(t)
	defer cleanUpSampleTest(svc, redisServer)

	sampleTypeEnum := "testSampleTypeEnum"
	sampleTypeName := "testSampleTypeName"
	sampleTypeCode := "testSampleTypeCode"
	sampleTypeDescription := "testSampleTypeDescription"
	primarySampleTypeGroup := "testPrimarySampleTypeGroup"

	// Seed the database
	sampleType, err := svc.dbClient.SampleType.Create().
		SetSampleTypeName(sampleTypeName).
		SetSampleTypeEnum(sampleTypeEnum).
		SetSampleTypeCode(sampleTypeCode).
		SetSampleTypeDescription(sampleTypeDescription).
		SetPrimarySampleTypeGroup(primarySampleTypeGroup).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to seed sample type: %v", err)
	}

	tubeTypes := []struct {
		TubeName       string
		TubeTypeEnum   string
		TubeTypeSymbol string
	}{
		{"testTubeC", "testTubeEnumC", "testTubeSymbolC"},
		{"testTubeD", "testTubeEnumD", "testTubeSymbolD"},
		{"testTubeE", "testTubeEnumE", "testTubeSymbolE"},
	}

	for _, tube := range tubeTypes {
		_, err := svc.dbClient.TubeType.Create().
			SetTubeName(tube.TubeName).
			SetTubeTypeEnum(tube.TubeTypeEnum).
			SetTubeTypeSymbol(tube.TubeTypeSymbol).
			AddSampleTypeIDs(sampleType.ID).
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to seed tube type: %v", err)
		}
	}

	// Run the service method
	actualSampleType, err := svc.GetTubeTypeViaSampleTypeCode(sampleTypeCode, ctx)

	// Verify main fields
	if actualSampleType.SampleTypeEnum != sampleTypeEnum ||
		actualSampleType.SampleTypeName != sampleTypeName ||
		actualSampleType.SampleTypeCode != sampleTypeCode ||
		actualSampleType.SampleTypeDescription != sampleTypeDescription ||
		actualSampleType.PrimarySampleTypeGroup != primarySampleTypeGroup {
		t.Fatal("sample type fields do not match expected values")
	}

	// Verify TubeTypes
	if len(actualSampleType.Edges.TubeTypes) != len(tubeTypes) {
		t.Fatalf("expected %d TubeTypes, got: %d", len(tubeTypes), len(actualSampleType.Edges.TubeTypes))
	}

	for i, actualTube := range actualSampleType.Edges.TubeTypes {
		expectedTube := tubeTypes[i]
		if actualTube.TubeName != expectedTube.TubeName ||
			actualTube.TubeTypeEnum != expectedTube.TubeTypeEnum ||
			actualTube.TubeTypeSymbol != expectedTube.TubeTypeSymbol {
			t.Fatalf("TubeType mismatch: expected %+v, got %+v", expectedTube, actualTube)
		}
	}
}

func TestGetSampleTypeViaTubeType(t *testing.T) {
	svc, redisServer, ctx := setupSampleTest(t)
	defer cleanUpSampleTest(svc, redisServer)

	sampleTypeEnum := "testSampleTypeEnum"
	sampleTypeName := "testSampleTypeName"
	sampleTypeCode := "testSampleTypeCode"
	sampleTypeDescription := "testSampleTypeDescription"
	primarySampleTypeGroup := "testPrimarySampleTypeGroup"

	// Seed the database
	sampleType, err := svc.dbClient.SampleType.Create().
		SetSampleTypeName(sampleTypeName).
		SetSampleTypeEnum(sampleTypeEnum).
		SetSampleTypeCode(sampleTypeCode).
		SetSampleTypeDescription(sampleTypeDescription).
		SetPrimarySampleTypeGroup(primarySampleTypeGroup).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to seed sample type: %v", err)
	}

	// Seed tube types
	tubeTypes := []struct {
		TubeName       string
		TubeTypeEnum   string
		TubeTypeSymbol string
	}{
		{"testTubeAA", "testTubeEnumA", "testTubeSymbolAA"},
		{"testTubeAB", "testTubeEnumA", "testTubeSymbolAB"},
		{"testTubeAC", "testTubeEnumA", "testTubeSymbolAC"},
	}

	expected := make(map[string]struct {
		TubeTypeEnum   string
		TubeTypeSymbol string
	})
	for _, tube := range tubeTypes {
		expected[tube.TubeName] = struct {
			TubeTypeEnum   string
			TubeTypeSymbol string
		}{tube.TubeTypeEnum, tube.TubeTypeSymbol}

		if _, err := svc.dbClient.TubeType.Create().
			SetTubeName(tube.TubeName).
			SetTubeTypeEnum(tube.TubeTypeEnum).
			SetTubeTypeSymbol(tube.TubeTypeSymbol).
			AddSampleTypeIDs(sampleType.ID).
			Save(ctx); err != nil {
			t.Fatalf("failed to seed tube type: %v", err)
		}
	}

	// Run the service method
	actualTubeTypes, err := svc.GetSampleTypeViaTubeType("testTubeEnumA", ctx)
	if err != nil {
		t.Fatalf("GetSampleTypeViaTubeType failed: %v", err)
	}

	// Verify the results
	if len(actualTubeTypes) != len(tubeTypes) {
		t.Fatalf("expected %d TubeTypes, got %d", len(tubeTypes), len(actualTubeTypes))
	}

	for _, actual := range actualTubeTypes {
		if expectedTube, ok := expected[actual.TubeName]; !ok {
			t.Fatalf("unexpected TubeName Got: %s", actual.TubeName)
		} else {
			if actual.TubeTypeEnum != expectedTube.TubeTypeEnum {
				t.Fatalf("for TubeName %s, expected TubeTypeEnum: %s, got: %s", actual.TubeName, expectedTube.TubeTypeEnum, actual.TubeTypeEnum)
			}
			if actual.TubeTypeSymbol != expectedTube.TubeTypeSymbol {
				t.Fatalf("for TubeName %s, expected TubeTypeSymbol: %s, got: %s", actual.TubeName, expectedTube.TubeTypeSymbol, actual.TubeTypeSymbol)
			}
		}

		// Validate linked SampleType
		actualSample := actual.Edges.SampleTypes[0]
		if actualSample == nil {
			t.Fatalf("TubeName %s: missing SampleType", actual.TubeName)
			continue
		}
		if actualSample.SampleTypeName != sampleTypeName ||
			actualSample.SampleTypeEnum != sampleTypeEnum ||
			actualSample.SampleTypeCode != sampleTypeCode ||
			actualSample.SampleTypeDescription != sampleTypeDescription ||
			actualSample.PrimarySampleTypeGroup != primarySampleTypeGroup {
			t.Fatalf("TubeName %s: mismatched SampleType fields", actual.TubeName)
		}
	}
}

func TestGenerateSampleID(t *testing.T) {
	svc, redisServer, ctx := setupSampleTest(t)
	defer cleanUpSampleTest(svc, redisServer)

	t.Run("Successful generation", func(t *testing.T) {
		resp, err := svc.GenerateSampleID(&pb.EmptyRequest{}, ctx)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Greater(t, resp.SampleId, int32(0))
	})
}

func TestGenerateBarcodeForSampleID(t *testing.T) {
	svc, redisServer, ctx := setupSampleTest(t)
	defer cleanUpSampleTest(svc, redisServer)

	t.Run("Valid sample ID", func(t *testing.T) {
		// Create a sample ID first
		sampleID, err := svc.GenerateSampleID(&pb.EmptyRequest{}, ctx)
		assert.NoError(t, err)
		assert.NotNil(t, sampleID)

		// Generate barcode for the sample ID
		req := &pb.GenerateBarcodeForSampleIdRequest{
			SampleId: sampleID.GetSampleId(),
		}
		resp, err := svc.GenerateBarcodeForSampleID(req, ctx)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.GetBarcode())
	})

	t.Run("Invalid sample ID", func(t *testing.T) {
		req := &pb.GenerateBarcodeForSampleIdRequest{
			SampleId: int32(-1),
		}
		resp, err := svc.GenerateBarcodeForSampleID(req, ctx)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid sample ID")
	})
}

func TestGetdailyCollectionSamples(t *testing.T) {
	svc, redisServer, ctx := setupSampleTest(t)
	defer cleanUpSampleTest(svc, redisServer)

	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	// Create patient with edges
	patient := svc.dbClient.Patient.Create().
		SaveX(ctx)

	// Create contact with edges
	svc.dbClient.Contact.Create().
		SetPatient(patient).
		SetContactType("phone").
		SetContactDetails("1234567890").
		ExecX(ctx)

	// Create sample with edges
	svc.dbClient.Sample.Create().
		SetPatient(patient).
		SetID(1).
		SetAccessionID("TEST001").
		SetDelayedHours(0).
		SetSampleReceivedTime(testTime).
		ExecX(ctx)

	tests := []struct {
		name      string
		startTime string
		endTime   string
		wantEmpty bool
	}{
		{
			name:      "Valid time range with sample",
			startTime: "2024-01-01 00:00:00",
			endTime:   "2024-01-01 23:59:59",
			wantEmpty: false,
		},
		{
			name:      "Valid time range without sample",
			startTime: "2023-01-01 00:00:00",
			endTime:   "2023-01-01 23:59:59",
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collections, _ := svc.GetdailyCollectionSamples(tt.startTime, tt.endTime, ctx)

			if tt.wantEmpty {
				assert.Empty(t, collections)
			} else {
				if assert.NotEmpty(t, collections, "Collections should not be empty for time range %s to %s", tt.startTime, tt.endTime) {
					assert.Equal(t, "1", collections[0].GetSampleId())
					assert.NotNil(t, collections[0].GetPatient())
					assert.NotEmpty(t, collections[0].GetPatient().GetPatientContact())
				}
			}
		})
	}
}

func TestGetdailyCheckNonReceivedSamples(t *testing.T) {
	svc, redisServer, ctx := setupSampleTest(t)
	defer cleanUpSampleTest(svc, redisServer)

	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	// Create patient with edges
	patient := svc.dbClient.Patient.Create().
		SaveX(ctx)

	// Create contact with edges
	svc.dbClient.Contact.Create().
		SetPatient(patient).
		SetContactType("phone").
		SetContactDetails("1234567890").
		ExecX(ctx)

	// Create order first with fields
	order := svc.dbClient.OrderInfo.Create().
		SetOrderCreateTime(testTime).
		SetOrderConfirmationNumber("TEST-ORDER-001").
		SaveX(ctx)

	// Create sample with all required edges but no received time
	svc.dbClient.Sample.Create().
		SetPatient(patient).
		SetID(1).
		SetAccessionID("TEST001").
		SetDelayedHours(0).
		SetOrder(order).
		ExecX(ctx)

	tests := []struct {
		name      string
		startTime string
		endTime   string
		wantEmpty bool
	}{
		{
			name:      "Valid time range with non-received sample",
			startTime: "2024-01-01 00:00:00",
			endTime:   "2024-01-01 23:59:59",
			wantEmpty: false,
		},
		{
			name:      "Valid time range without sample",
			startTime: "2023-01-01 00:00:00",
			endTime:   "2023-01-01 23:59:59",
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collections, _ := svc.GetdailyCheckNonReceivedSamples(tt.startTime, tt.endTime, ctx)

			if tt.wantEmpty {
				assert.Empty(t, collections)
			} else {
				if assert.NotEmpty(t, collections, "Collections should not be empty for time range %s to %s", tt.startTime, tt.endTime) {
					assert.Equal(t, "1", collections[0].GetSampleId())
					assert.Equal(t, "TEST001", collections[0].GetAccessionId())
					assert.NotNil(t, collections[0].GetPatient())
					assert.NotEmpty(t, collections[0].GetPatient().GetPatientContact())
				}
			}
		})
	}
}
