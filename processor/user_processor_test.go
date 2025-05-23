package processor

import (
	"context"
	"coresamples/common"
	"coresamples/ent/enttest"
	pb "coresamples/proto"
	"coresamples/publisher"
	"coresamples/tasks"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stvp/tempredis"
	"strconv"
	"testing"
)

func setupUserProcessorTest(t *testing.T) (*UserProcessor, *tempredis.Server, tasks.AsynqClient) {
	dataSource := "file:ent?mode=memory&_fk=1"
	dbClient := enttest.Open(t, "sqlite3", dataSource)
	err := dbClient.Schema.Create(context.Background())
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
	asynqClient := tasks.NewMockAsynqClient()

	// This will start a mock publisher in the memory
	publisher.InitMockPublisher()

	client := common.NewRedisClient(redisClient, redisClient)

	processor := &UserProcessor{
		Processor: InitProcessor(dbClient,
			client,
			context.Background()),
	}

	return processor, server, asynqClient
}

func cleanUpUserProcessorTest(p *UserProcessor, s *tempredis.Server, asynqClient tasks.AsynqClient) {
	var err error
	if err = s.Kill(); err != nil {
		common.Error(err)
	}
	if p.dbClient != nil {
		if err = p.dbClient.Close(); err != nil {
			common.Error(err)
		}
	}

	publisher.GetPublisher().GetWriter().Close()
	asynqClient.Close()
}

func TestHandleHubspotEvent(t *testing.T) {
	processor, server, asynqClient := setupUserProcessorTest(t)
	defer cleanUpUserProcessorTest(processor, server, asynqClient)

	sales1 := processor.dbClient.InternalUser.
		Create().
		SetInternalUserEmail("sales1@gmail.com").
		SetInternalUserRole("sales").
		SaveX(processor.ctx)
	sales2 := processor.dbClient.InternalUser.
		Create().
		SetInternalUserEmail("sales2@gmail.com").
		SetInternalUserRole("sales").
		SaveX(processor.ctx)
	customer := processor.dbClient.Customer.Create().SetSales(sales1).SaveX(processor.ctx)

	assert.Equal(t, customer.SalesID, sales1.ID)

	task, _ := tasks.NewHubspotEventTask(&tasks.HubspotEventTask{
		Event: &pb.HubspotEvent{
			Schema: &pb.HubspotEvent_Schema{
				ProviderId: strconv.Itoa(customer.ID),
				OwnerEmail: sales2.InternalUserEmail,
			},
		},
	})

	assert.NoError(t, processor.HandleHubspotEvent(processor.ctx, task))
	customer = customer.Update().SaveX(processor.ctx)
	assert.Equal(t, customer.SalesID, sales2.ID)
}
