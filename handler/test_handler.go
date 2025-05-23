package handler

import (
	"context"
	"coresamples/common"
	pb "coresamples/proto"
	"coresamples/service"
	"coresamples/util"
	"strconv"

	"github.com/getsentry/sentry-go"
)

type TestHandler struct {
	TestService service.ITestService
}

func (th *TestHandler) GetTest(ctx context.Context, request *pb.GetTestRequest, response *pb.GetTestResponse) error {
	var ids []int
	for _, id := range request.TestIds {
		testId, err := strconv.Atoi(id)
		if err != nil {
			continue
		}
		ids = append(ids, testId)
	}
	tests, err := th.TestService.GetTest(ids, ctx)
	if err != nil {
		sentry.CaptureMessage(err.Error())
		return err
	}
	err = util.Swap(tests, &response.Test)

	if err != nil {
		sentry.CaptureMessage(err.Error())
		return err
	}

	for idx := range tests {
		err = util.Swap(tests[idx].Edges.TestDetails, &response.Test[idx].TestDetails)
		if err != nil {
			sentry.CaptureMessage(err.Error())
			return err
		}
	}
	return nil
}

func (th *TestHandler) GetTestField(ctx context.Context, request *pb.GetTestFieldRequest, response *pb.GetTestResponse) error {
	var ids []int
	for _, id := range request.TestIds {
		testId, err := strconv.Atoi(id)
		if err != nil {
			continue
		}
		ids = append(ids, testId)
	}
	tests, err := th.TestService.GetTestField(ids, request.TestDetailNames, ctx)
	if err != nil {
		sentry.CaptureMessage(err.Error())
		return err
	}
	err = util.Swap(tests, &response.Test)

	if err != nil {
		sentry.CaptureMessage(err.Error())
		return err
	}

	for idx := range tests {
		err = util.Swap(tests[idx].Edges.TestDetails, &response.Test[idx].TestDetails)
		if err != nil {
			sentry.CaptureMessage(err.Error())
			return err
		}
	}
	return nil
}

func (th *TestHandler) CreateTest(ctx context.Context, request *pb.CreateTestRequest, response *pb.CreateTestResponse) error {
	test, err := th.TestService.CreateTest(request, ctx)
	if err != nil {
		common.Error(err)
		return err
	}
	response.CreatedAt = test.CreatedTime.Format("2006-01-02 15:04:05")
	response.UpdatedAt = test.UpdatedTime.Format("2006-01-02 15:04:05")

	err = util.Swap(test, response)

	if err != nil {
		sentry.CaptureMessage(err.Error())
		return err
	}
	return nil
}

func (th *TestHandler) GetTestIDsFromTestCodes(ctx context.Context, request *pb.GetTestIDsFromTestCodesRequest, response *pb.GetTestIDsFromTestCodesResponse) error {
	res, err := th.TestService.GetTestIDsFromTestCodes(request.TestCodes, ctx)
	if err != nil {
		common.Error(err)
		return err
	}
	for code, ids := range res {
		if len(ids) == 0 {
			continue
		}
		response.Response = append(response.Response,
			&pb.TestCodetoTestIDsList{
				TestCode: code,
				TestIds:  ids,
			})
	}
	return nil
}

func (th *TestHandler) GetDuplicateAssayGroupTest(ctx context.Context, request *pb.GetDuplicateAssayGroupTestRequest, response *pb.GetDuplicateAssayGroupTestResponse) error {
	testId, err := strconv.Atoi(request.TestId)
	if err != nil {
		common.Error(err)
		return err
	}
	ids, err := th.TestService.GetDuplicateAssayGroupTest(testId, ctx)
	if err != nil {
		common.Error(err)
		return err
	}
	response.DuplicateTests = ids
	return nil
}

func (th *TestHandler) GetTestTubeTypes(ctx context.Context, request *pb.GetTestTubeTypesRequest, response *pb.GetTestTubeTypesResponse) error {
	var ids []int
	for _, id := range request.TestIds {
		ids = append(ids, int(id))
	}

	TestTubeInfos, err := th.TestService.GetTestTubeTypes(ids, ctx)
	if err != nil {
		common.Error(err)
		return err
	}
	response.TestTubeInfos = TestTubeInfos
	return nil
}
