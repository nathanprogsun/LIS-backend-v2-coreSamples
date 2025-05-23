package processor

import (
	"context"
	"coresamples/common"
	"coresamples/ent"
)

type Processor struct {
	dbClient    *ent.Client
	redisClient *common.RedisClient
	ctx         context.Context
}

func InitProcessor(dbClient *ent.Client, redisClient *common.RedisClient, ctx context.Context) Processor {
	return Processor{
		dbClient:    dbClient,
		redisClient: redisClient,
		ctx:         ctx,
	}
}
