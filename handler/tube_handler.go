package handler

import (
	"context"
	pb "coresamples/proto"
	"coresamples/service"
	"coresamples/util"
	"github.com/getsentry/sentry-go"
)

type TubeHandler struct {
	TubeService service.ITubeService
}

func (th *TubeHandler) GetRequiredTubeVolume(ctx context.Context, request *pb.RequiredTubeVolumeRequest, response *pb.RequiredTubeVolumeResponse) error {
	testIds := request.GetTestIds()
	resp, err := th.TubeService.GetRequiredTubeVolume(testIds, ctx)

	if err != nil {
		response.Message = err.Error()
		sentry.CaptureMessage(response.Message)
		return err
	}
	response.Message = MsgSuccess
	err = util.Swap(resp, response)

	if err != nil {
		response.Message = err.Error()
		sentry.CaptureMessage(response.Message)
		return err
	}
	return nil
}

func (th *TubeHandler) GetTestsByBloodType(ctx context.Context, request *pb.BloodType, response *pb.TestIDs) error {
	tests, err := th.TubeService.GetTestsByBloodType(request.Blood, ctx)
	if err != nil {
		sentry.CaptureMessage(err.Error())
		return err
	}
	response.TestIds = tests
	return nil
}

func (th *TubeHandler) GetTube(ctx context.Context, req *pb.TubeID, resp *pb.Tube) error {
	tube, err := th.TubeService.GetTube(req.TubeId, ctx)
	if err != nil {
		return err
	}
	err = util.Swap(tube, resp)

	if err != nil {
		return err
	}

	err = util.Swap(tube.Edges.TubeType, &resp.TubeTypes)
	if err != nil {
		return err
	}
	return nil
}

func (th *TubeHandler) GetTubeTests(ctx context.Context, req *pb.TubeID, resp *pb.TubeTests) error {
	return nil
}
