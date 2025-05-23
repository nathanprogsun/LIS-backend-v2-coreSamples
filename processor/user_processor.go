package processor

import (
	"context"
	"coresamples/common"
	"coresamples/dbutils"
	"coresamples/ent"
	"coresamples/tasks"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"strconv"
)

type UserProcessor struct {
	Processor
	//rs            *redsync.Redsync
}

func NewUserProcessor(dbClient *ent.Client, redisClient *common.RedisClient) *UserProcessor {
	return &UserProcessor{
		Processor: InitProcessor(dbClient, redisClient, context.Background()),
	}
}

func (p *UserProcessor) HandleHubspotEvent(ctx context.Context, t *asynq.Task) error {
	task := &tasks.HubspotEventTask{}
	if err := json.Unmarshal(t.Payload(), task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	// if common.Env.DryRun {
	// 	return nil
	// }
	schema := task.Event.Schema
	if schema == nil {
		common.Debugf("received nil schema %v", task)
		return fmt.Errorf("received nil schema %v", asynq.SkipRetry)
	}
	customerId, err := strconv.Atoi(schema.ProviderId)
	if err != nil {
		return fmt.Errorf("unable to parse customer id %v %w", err, asynq.SkipRetry)
	}
	// We'll support assigning cam to customer directly in the future
	//camName := schema.CamName
	salesEmail := schema.OwnerEmail
	customer, err := dbutils.GetCustomerByCustomerID(customerId, p.dbClient, p.ctx)
	if err != nil {
		return fmt.Errorf("unable to find customer %v %w", err, asynq.SkipRetry)
	}
	sales, err := dbutils.GetSalesByEmail(salesEmail, p.dbClient, p.ctx)
	if err != nil {
		return fmt.Errorf("unable to find sales %v %w", err, asynq.SkipRetry)
	}
	return customer.Update().SetSales(sales).Exec(ctx)
}
