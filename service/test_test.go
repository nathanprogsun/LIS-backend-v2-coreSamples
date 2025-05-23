package service

import (
	"context"
	"coresamples/common"
	"coresamples/dbutils"
	"coresamples/ent/enttest"
	pb "coresamples/proto"
	"coresamples/publisher"
	"coresamples/util"
	"fmt"
	"strconv"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stvp/tempredis"
)

func setupTestTest(t *testing.T) (*TestService, *tempredis.Server, context.Context) {
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

	s := &TestService{
		Service: InitService(
			dbClient,
			common.NewRedisClient(redisClient, redisClient)),
	}
	// This will start a mock publisher in the memory
	publisher.InitMockPublisher()

	return s, server, ctx
}

func cleanUpTestTest(svc *TestService, s *tempredis.Server) {
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

func TestGetTest(t *testing.T) {
	svc, redisServer, ctx := setupTestTest(t)
	defer cleanUpTestTest(svc, redisServer)

	// Create test data
	test1, err := svc.dbClient.Test.Create().
		SetTestName("Test1").
		SetTestCode("Test1 testcode").
		SetDisplayName("Test1 testDisplayName").
		SetTestDescription("Test1 testDescription").
		SetAssayName("Test1 assayName").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create test1: %v", err)
	}
	_, err = svc.dbClient.Test.Create().
		SetTestName("Test2").
		SetTestCode("Test2 testcode").
		SetDisplayName("Test2 testDisplayName").
		SetTestDescription("Test2 testDescription").
		SetAssayName("Test2 assayName").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create test2: %v", err)
	}

	// 1️⃣ Case: Get all tests (no testIds), simulate Redis miss
	allTests, err := svc.GetTest(nil, ctx)
	if err != nil {
		t.Fatalf("unexpected error in GetTest(nil): %v", err)
	}
	if len(allTests) < 2 {
		t.Errorf("expected at least 2 tests, got %d", len(allTests))
	}

	// 2️⃣ Case: Get specific testId, simulate Redis miss again
	svc.redisClient.Del(ctx, dbutils.KeyGetTestByTestId(test1.ID)) // force cache miss

	testsById, err := svc.GetTest([]int{test1.ID}, ctx)
	if err != nil {
		t.Fatalf("unexpected error for specific testId: %v", err)
	}
	if len(testsById) != 1 || testsById[0].ID != test1.ID {
		t.Errorf("expected 1 test with ID %d, got %v", test1.ID, testsById)
	}

	// 3️⃣ Case: Get from Redis cache (warm hit)
	cached, err := svc.GetTest([]int{test1.ID}, ctx)
	if err != nil || len(cached) != 1 || cached[0].ID != test1.ID {
		t.Errorf("expected cache hit for ID %d, got error %v and result %+v", test1.ID, err, cached)
	}
}

func TestGetTestField(t *testing.T) {
	svc, redisServer, ctx := setupTestTest(t)
	defer cleanUpTestTest(svc, redisServer)

	// Seed Test1
	test1, err := svc.dbClient.Test.Create().
		SetTestName("Test1").
		SetTestCode("Test1 testcode").
		SetDisplayName("Test1 testDisplayName").
		SetTestDescription("Test1 testDescription").
		SetAssayName("Test1 assayName").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create test1: %v", err)
	}

	detailNames := []string{"test_detail1", "test_detail2", "test_detail3"}
	for _, name := range detailNames {
		_, err := svc.dbClient.TestDetail.Create().
			SetTestDetailName(name).
			SetTestID(test1.ID).
			SetTestDetailsValue(fmt.Sprintf("value_for_%s", name)).
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to create test detail for %s: %v", name, err)
		}
	}

	// Seed Test2
	test2, err := svc.dbClient.Test.Create().
		SetTestName("Test2").
		SetTestCode("Test2 testcode").
		SetDisplayName("Test2 testDisplayName").
		SetTestDescription("Test2 testDescription").
		SetAssayName("Test2 assayName").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create test2: %v", err)
	}

	_, err = svc.dbClient.TestDetail.Create().
		SetTestDetailName("test_detail1").
		SetTestDetailsValue("Test2 test_detail1").
		SetTestID(test2.ID).
		Save(ctx)

	if err != nil {
		t.Fatalf("failed to create test2: %v", err)
	}

	// Only request specific detail names
	requestedDetails := []string{"test_detail1", "test_detail2"}

	// Run the test
	results, err := svc.GetTestField([]int{test1.ID, test2.ID}, requestedDetails, ctx)
	if err != nil {
		t.Fatalf("unexpected error from GetTestField: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// Validate each test
	for _, test := range results {
		if len(test.Edges.TestDetails) == 0 {
			t.Errorf("test %d returned no TestDetails", test.ID)
		}
		for _, td := range test.Edges.TestDetails {
			found := false
			for _, expected := range requestedDetails {
				if td.TestDetailName == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("unexpected TestDetail name: %s in test %d", td.TestDetailName, test.ID)
			}
		}
	}
}

func TestCreateTest(t *testing.T) {
	svc, redisServer, ctx := setupTestTest(t)
	defer cleanUpTestTest(svc, redisServer)

	req := &pb.CreateTestRequest{
		IsActive:        true,
		TestName:        "Test ABC",
		TestCode:        "ABC123",
		DisplayName:     "ABC Display",
		TestDescription: "This is a test description",
		AssayName:       "ABC Assay",
	}

	created, err := svc.CreateTest(req, ctx)
	if err != nil {
		t.Fatalf("CreateTest failed: %v", err)
	}

	// Check values
	if created.TestName != req.TestName {
		t.Errorf("expected TestName %q, got %q", req.TestName, created.TestName)
	}
	if created.TestCode != req.TestCode {
		t.Errorf("expected TestCode %q, got %q", req.TestCode, created.TestCode)
	}
	if created.DisplayName != req.DisplayName {
		t.Errorf("expected DisplayName %q, got %v", req.DisplayName, created.DisplayName)
	}
	if created.TestDescription == nil || *created.TestDescription != req.TestDescription {
		t.Errorf("expected TestDescription %q, got %v", req.TestDescription, created.TestDescription)
	}
	if created.AssayName == nil || *created.AssayName != req.AssayName {
		t.Errorf("expected AssayName %q, got %v", req.AssayName, created.AssayName)
	}
	if !created.IsActive {
		t.Errorf("expected IsActive true, got false")
	}
}

func TestGetTestIDsFromTestCodes(t *testing.T) {
	svc, redisServer, ctx := setupTestTest(t)
	defer cleanUpTestTest(svc, redisServer)

	// Seed test1 with full fields
	test1, err := svc.dbClient.Test.Create().
		SetTestName("Test1").
		SetTestCode("CODE123").
		SetDisplayName("Display1").
		SetTestDescription("Description1").
		SetAssayName("Assay1").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create test1: %v", err)
	}

	// Seed test2 with full fields
	test2, err := svc.dbClient.Test.Create().
		SetTestName("Test2").
		SetTestCode("CODE123").
		SetDisplayName("Display2").
		SetTestDescription("Description2").
		SetAssayName("Assay2").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create test2: %v", err)
	}

	// Seed test3 with different code
	test3, err := svc.dbClient.Test.Create().
		SetTestName("Test3").
		SetTestCode("CODE999").
		SetDisplayName("Display3").
		SetTestDescription("Description3").
		SetAssayName("Assay3").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create test3: %v", err)
	}

	// Run the method
	result, err := svc.GetTestIDsFromTestCodes([]string{"CODE123", "CODE999"}, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Validate CODE123 → [test1.ID, test2.ID]
	expected123 := []int32{int32(test1.ID), int32(test2.ID)}
	got123 := result["CODE123"]
	if !util.SliceEqual(expected123, got123) {
		t.Errorf("CODE123 mismatch: expected %v, got %v", expected123, got123)
	}

	// Validate CODE999 → [test3.ID]
	expected999 := []int32{int32(test3.ID)}
	got999 := result["CODE999"]
	if !util.SliceEqual(expected999, got999) {
		t.Errorf("CODE999 mismatch: expected %v, got %v", expected999, got999)
	}
}

func TestGetDuplicateAssayGroupTest(t *testing.T) {
	svc, redisServer, ctx := setupTestTest(t)
	defer cleanUpTestTest(svc, redisServer)

	const detailName = "test_duplicate_asssay_name"

	// ---- CASE 1: Test has valid duplicate_assay reference ----

	// Create reference test (this ID will be used as detail value)
	refTest, err := svc.dbClient.Test.Create().
		SetTestName("RefTest").
		SetTestCode("REFCODE").
		SetDisplayName("Ref Display").
		SetTestDescription("Ref Desc").
		SetAssayName("Ref Assay").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create refTest: %v", err)
	}

	// Create test that points to refTest.ID via detail
	targetTest, err := svc.dbClient.Test.Create().
		SetTestName("Target").
		SetTestCode("TARGET").
		SetDisplayName("Target Display").
		SetTestDescription("Target Desc").
		SetAssayName("Target Assay").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create targetTest: %v", err)
	}

	_, err = svc.dbClient.TestDetail.Create().
		SetTestID(targetTest.ID).
		SetTestDetailName(detailName).
		SetTestDetailsValue(strconv.Itoa(refTest.ID)).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create detail for targetTest: %v", err)
	}

	// Create another test with same detail value
	peerTest, err := svc.dbClient.Test.Create().
		SetTestName("Peer").
		SetTestCode("PEERCODE").
		SetDisplayName("Peer Display").
		SetTestDescription("Peer Desc").
		SetAssayName("Peer Assay").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create peerTest: %v", err)
	}

	_, err = svc.dbClient.TestDetail.Create().
		SetTestID(peerTest.ID).
		SetTestDetailName(detailName).
		SetTestDetailsValue(strconv.Itoa(refTest.ID)).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create detail for peerTest: %v", err)
	}

	// Run test
	res, err := svc.GetDuplicateAssayGroupTest(targetTest.ID, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []int32{int32(targetTest.ID), int32(peerTest.ID)}
	if !util.SliceEqual(res, expected) && !util.SliceEqual(res, []int32{expected[1], expected[0]}) {
		t.Errorf("unexpected result, got %v, expected %v", res, expected)
	}

	// ---- CASE 2: No TestDetail present ----
	noDetailTest, err := svc.dbClient.Test.Create().
		SetTestName("NoDetail").
		SetTestCode("ND").
		SetDisplayName("NoDisplay").
		SetTestDescription("NoDesc").
		SetAssayName("NoAssay").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create noDetailTest: %v", err)
	}

	res, err = svc.GetDuplicateAssayGroupTest(noDetailTest.ID, ctx)
	if err != nil {
		t.Fatalf("unexpected error for no detail test: %v", err)
	}
	if res != nil {
		t.Errorf("expected nil result for no detail test, got %v", res)
	}

	// ---- CASE 3: Detail exists but value is "0" ----
	zeroValueTest, err := svc.dbClient.Test.Create().
		SetTestName("Zero").
		SetTestCode("ZERO").
		SetDisplayName("ZeroDisplay").
		SetTestDescription("ZeroDesc").
		SetAssayName("ZeroAssay").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create zeroValueTest: %v", err)
	}

	_, _ = svc.dbClient.TestDetail.Create().
		SetTestID(zeroValueTest.ID).
		SetTestDetailName(detailName).
		SetTestDetailsValue("0").
		Save(ctx)

	res, err = svc.GetDuplicateAssayGroupTest(zeroValueTest.ID, ctx)
	if err != nil {
		t.Fatalf("unexpected error for zero value test: %v", err)
	}
	expectedSingle := []int32{int32(zeroValueTest.ID)}
	if !util.SliceEqual(res, expectedSingle) {
		t.Errorf("expected only self ID, got %v", res)
	}
}

func TestGetTestTubeTypes(t *testing.T) {
	svc, redisServer, ctx := setupTestTest(t)
	defer cleanUpTestTest(svc, redisServer)

	// Create SampleType
	sampleType, err := svc.dbClient.SampleType.Create().
		SetSampleTypeName("Serum").
		SetSampleTypeEnum("SERUM_ENUM").
		SetSampleTypeCode("SERUM123").
		SetSampleTypeDescription("Serum type").
		SetPrimarySampleTypeGroup("GroupA").
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create SampleType: %v", err)
	}

	// Define tube types
	tubeTypes := []struct {
		TubeName       string
		TubeTypeEnum   string
		TubeTypeSymbol string
	}{
		{"Tube A", "TUBE_A", "TA"},
		{"Tube B", "TUBE_B", "TB"},
	}

	// Create TubeTypes using the format you provided
	for _, tube := range tubeTypes {
		_, err := svc.dbClient.TubeType.Create().
			SetTubeName(tube.TubeName).
			SetTubeTypeEnum(tube.TubeTypeEnum).
			SetTubeTypeSymbol(tube.TubeTypeSymbol).
			AddSampleTypeIDs(sampleType.ID).
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to create TubeType: %v", err)
		}
	}

	// Create Test
	test, err := svc.dbClient.Test.Create().
		SetTestName("Test1").
		SetTestCode("T1").
		SetDisplayName("T1 Disp").
		SetTestDescription("desc").
		SetAssayName("assay").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create Test: %v", err)
	}

	// Link TestDetail with sampleTypeCode
	_, err = svc.dbClient.TestDetail.Create().
		SetTestID(test.ID).
		SetTestDetailName("sample_type_code").
		SetTestDetailsValue("SERUM123").
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create TestDetail: %v", err)
	}

	// Run the function
	res, err := svc.GetTestTubeTypes([]int{test.ID}, ctx)
	if err != nil {
		t.Fatalf("GetTestTubeTypes error: %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if res[0].TestId != int32(test.ID) {
		t.Errorf("unexpected test ID: got %d", res[0].TestId)
	}
	if len(res[0].SampleTypes) != 1 {
		t.Errorf("expected 1 sample type, got %d", len(res[0].SampleTypes))
	}
	sample := res[0].SampleTypes[0]
	if sample.SampleType != "SERUM_ENUM" {
		t.Errorf("unexpected SampleType: got %s", sample.SampleType)
	}
	if !util.Contains(sample.TubeType, "TUBE_A") || !util.Contains(sample.TubeType, "TUBE_B") {
		t.Errorf("missing expected tube types: got %v", sample.TubeType)
	}
}
