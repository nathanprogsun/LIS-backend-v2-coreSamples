package service

import (
	"context"
	"coresamples/common"
	"coresamples/ent"
)

type ITubeService interface {
	GetRequiredTubeVolume(testIds []int32, ctx context.Context) (response *InternalRequiredTubeVolumeResponse, err error)
	GetTestsByBloodType(bloodType bool, ctx context.Context) ([]int32, error)
	GetTube(tubeID string, ctx context.Context) (*ent.Tube, error)
}

type TubeService struct {
	*GetRequiredTubeVolumeService
	*TubeInfoService
}

func newTubeService(dbClient *ent.Client, redisClient *common.RedisClient) ITubeService {
	requiredTubeVolumeService := &GetRequiredTubeVolumeService{}
	testInfoService := &TubeInfoService{}
	requiredTubeVolumeService.Init(dbClient, redisClient)
	testInfoService.Init(dbClient, redisClient)
	return &TubeService{
		GetRequiredTubeVolumeService: requiredTubeVolumeService,
		TubeInfoService:              testInfoService,
	}
}
