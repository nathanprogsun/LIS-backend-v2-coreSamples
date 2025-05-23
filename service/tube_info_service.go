package service

import (
	"context"
	"coresamples/common"
	"coresamples/dbutils"
	"coresamples/ent"
)

type TubeInfoService struct {
	Service
}

func (s *TubeInfoService) Init(dbClient *ent.Client, redisClient *common.RedisClient) {
	s.Service = InitService(dbClient, redisClient)
}

func (s *TubeInfoService) GetTestsByBloodType(bloodType bool, ctx context.Context) ([]int32, error) {
	tests, err := dbutils.GetTestsByBloodType(bloodType, s.dbClient, ctx)
	if err != nil {
		return []int32{}, nil
	}
	var ids []int32
	for _, test := range tests {
		ids = append(ids, int32(test.ID))
	}
	return ids, nil
}

func (s *TubeInfoService) GetTube(tubeID string, ctx context.Context) (*ent.Tube, error) {
	return dbutils.GetTubeByTubeID(tubeID, s.dbClient, ctx)
}
