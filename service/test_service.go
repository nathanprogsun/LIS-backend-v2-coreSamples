package service

import (
	"context"
	"coresamples/common"
	"coresamples/dbutils"
	"coresamples/ent"
	pb "coresamples/proto"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const duplicateAssayTestDetailName = "test_duplicate_asssay_name"

type ITestService interface {
	GetTest(testIds []int, ctx context.Context) ([]*ent.Test, error)
	GetTestField(testIds []int, testDetailNames []string, ctx context.Context) ([]*ent.Test, error)
	CreateTest(test *pb.CreateTestRequest, ctx context.Context) (*ent.Test, error)
	GetTestIDsFromTestCodes(testCodes []string, ctx context.Context) (map[string][]int32, error)
	GetDuplicateAssayGroupTest(testId int, ctx context.Context) ([]int32, error)
	GetTestTubeTypes(testIds []int, ctx context.Context) ([]*pb.GetTestTubeTypesResponse_TestTubeInfo, error)
}

type TestService struct {
	Service
}

func NewTestService(dbClient *ent.Client, redisClient *common.RedisClient) ITestService {
	s := &TestService{
		Service: InitService(dbClient, redisClient),
	}
	return s
}

func (s *TestService) GetTest(testIds []int, ctx context.Context) ([]*ent.Test, error) {
	trackingID := uuid.New().String()
	if len(testIds) == 0 {
		redisKey := dbutils.KeyGetTestAll()
		redisResult, err := s.redisClient.Get(ctx, redisKey).Result()
		if err == nil && redisResult != "" {
			var cachedTests []*ent.Test
			if err := json.Unmarshal([]byte(redisResult), &cachedTests); err == nil {
				return cachedTests, nil
			} else {
				common.ErrorLogger("[tracking:%s] failed to unmarshal get_test_all: %v", trackingID, err)
			}
		}

		// Cache miss or error: query DB
		allTests, err := dbutils.GetAllTests(s.dbClient, ctx)
		if err != nil {
			common.ErrorLogger("[tracking:%s] failed to fetch all tests from DB: %v", trackingID, err)
			return nil, err
		}

		// Store in Redis
		if jsonBytes, err := json.Marshal(allTests); err == nil {
			_ = s.redisClient.SetEX(ctx, redisKey, jsonBytes, 24*time.Hour).Err()
		}

		return allTests, nil
	}

	// Case 2: Specific testIds provided
	var resultArray []*ent.Test
	for _, id := range testIds {
		redisKey := dbutils.KeyGetTestByTestId(id)
		redisResult, err := s.redisClient.Get(ctx, redisKey).Result()

		if err == nil && redisResult != "" {
			var cachedTest ent.Test
			if err := json.Unmarshal([]byte(redisResult), &cachedTest); err == nil {
				resultArray = append(resultArray, &cachedTest)
				continue
			} else {
				common.ErrorLogger("[tracking:%s] failed to unmarshal test_id:%d: %v", trackingID, id, err)
			}
		}

		// Cache miss or unmarshal fail: fetch from DB
		testEnt, err := dbutils.GetTestByTestId(id, s.dbClient, ctx)
		if err != nil {
			common.ErrorLogger("[tracking:%s] failed to fetch test_id:%d from DB: %v", trackingID, id, err)
			continue
		}

		resultArray = append(resultArray, testEnt)

		if jsonBytes, err := json.Marshal(testEnt); err == nil {
			_ = s.redisClient.SetEX(ctx, redisKey, jsonBytes, 24*time.Hour).Err()
		}
	}

	return resultArray, nil
}

func (s *TestService) GetTestField(testIds []int, testDetailNames []string, ctx context.Context) ([]*ent.Test, error) {
	return dbutils.GetTestsWithFields(testIds, testDetailNames, s.dbClient, ctx)
}

func (s *TestService) CreateTest(test *pb.CreateTestRequest, ctx context.Context) (*ent.Test, error) {
	return dbutils.CreateTest(test, s.dbClient, ctx)
}

func (s *TestService) GetTestIDsFromTestCodes(testCodes []string, ctx context.Context) (map[string][]int32, error) {
	ret := make(map[string][]int32)
	for _, code := range testCodes {
		res, err := dbutils.GetTestIdsByTestCode(code, s.dbClient, ctx)
		fmt.Print(res)
		if err != nil {
			common.Error(err)
			continue
		}
		var ids []int32
		for _, id := range res {
			ids = append(ids, int32(id))
		}
		ret[code] = ids
	}
	return ret, nil
}

func (s *TestService) GetDuplicateAssayGroupTest(testId int, ctx context.Context) ([]int32, error) {
	test, err := dbutils.GetTestByTestIdWithDetailName(testId, duplicateAssayTestDetailName, s.dbClient, ctx)
	if err != nil {
		return nil, err
	}

	if len(test.Edges.TestDetails) == 0 {
		return nil, nil
	}
	dupAssayTestId := test.Edges.TestDetails[0].TestDetailsValue
	if dupAssayTestId != "0" {
		dupAssayTestDetails, err := dbutils.GetTestDetailsByDetailValueAndDetailName(dupAssayTestId, duplicateAssayTestDetailName, s.dbClient, ctx)
		if err != nil {
			return nil, err
		}
		ids := make(map[int]bool)
		ids[testId] = true
		for _, detail := range dupAssayTestDetails {
			ids[detail.TestID] = true
		}
		var ret []int32
		for id := range ids {
			ret = append(ret, int32(id))
		}
		return ret, nil
	} else {
		return []int32{int32(testId)}, nil
	}
}

func (s *TestService) GetTestTubeTypes(testIds []int, ctx context.Context) ([]*pb.GetTestTubeTypesResponse_TestTubeInfo, error) {
	testsInfo, err := dbutils.GetTestsByTestIds(testIds, s.dbClient, ctx)
	if err != nil {
		return nil, err
	}
	var result []*pb.GetTestTubeTypesResponse_TestTubeInfo
	for _, test := range testsInfo {
		var sampleTypes []*pb.GetTestTubeTypesResponse_TestTubeInfo_SampleType
		for _, detail := range test.Edges.TestDetails {
			tubeTypeInfo, err := dbutils.GetTubeTypeInfoBySampleTypeCode(detail.TestDetailsValue, s.dbClient, ctx)
			if err != nil {
				return nil, err
			}
			if tubeTypeInfo != nil {
				var tubeTypeEnums []string
				for _, tube := range tubeTypeInfo.Edges.TubeTypes {
					tubeTypeEnums = append(tubeTypeEnums, tube.TubeTypeEnum)
				}

				sampleTypes = append(sampleTypes, &pb.GetTestTubeTypesResponse_TestTubeInfo_SampleType{
					SampleType: tubeTypeInfo.SampleTypeEnum,
					TubeType:   tubeTypeEnums,
				})
			}
		}

		result = append(result, &pb.GetTestTubeTypesResponse_TestTubeInfo{
			TestId:      int32(test.ID),
			SampleTypes: sampleTypes,
		})

	}
	return result, nil
}
