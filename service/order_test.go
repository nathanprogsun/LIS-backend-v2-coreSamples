package service

import (
	"context"
	"coresamples/common"
	"coresamples/ent/enttest"
	"coresamples/publisher"
	"github.com/go-redis/redis/v8"
	"github.com/stvp/tempredis"
	"testing"
)

func setupOrderTest(t *testing.T) (*OrderService, *tempredis.Server) {
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

	redisClient := redis.NewClient(&redis.Options{
		Network: "unix",
		Addr:    server.Socket(),
	})

	s := &OrderService{
		Service: InitService(
			dbClient,
			common.NewRedisClient(redisClient, redisClient)),
	}

	// This will start a mock publisher in the memory
	publisher.InitMockPublisher()
	return s, server
}

func cleanupOrderTest(svc *OrderService, server *tempredis.Server) {
	var err error
	if err = server.Kill(); err != nil {
		common.Error(err)
	}
	if svc.dbClient != nil {
		if err = svc.dbClient.Close(); err != nil {
			common.Error(err)
		}
	}
	publisher.GetPublisher().GetWriter().Close()
}

func TestGetOrder(t *testing.T) {
	svc, redisServer := setupOrderTest(t)
	defer cleanupOrderTest(svc, redisServer)

	ctx := context.Background()
	order, err := svc.dbClient.OrderInfo.Create().
		SetOrderConfirmationNumber("test_order").Save(ctx)
	flag, err := svc.dbClient.OrderFlag.Create().
		AddFlaggedOrders(order).
		SetOrderFlagName("test_order_flag").Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	o, err := svc.GetOrder(order.ID, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if o.ID != order.ID {
		t.Fatal("unmatched order id")
	}
	if len(o.Edges.OrderFlags) != 1 {
		t.Fatal("order should have one flag")
	} else {
		if o.Edges.OrderFlags[0].ID != flag.ID {
			t.Fatal("unmatched order flag id")
		}
	}
}

func TestCancelOrder(t *testing.T) {

}

func TestRestoreCanceledOrder(t *testing.T) {

}

func TestAddOrderFlag(t *testing.T) {

}

// important!
func TestFlagOrder(t *testing.T) {

}

func TestFlagOrdersWithSampleId(t *testing.T) {

}

func TestUnflagOrder(t *testing.T) {

}

func TestGetOrderStatusForDisplay(t *testing.T) {

}

func TestRerunSampleTests(t *testing.T) {

}

func TestRestoreOrderStatus(t *testing.T) {

}

func TestTriggerOrderTransmissionOnSampleReceiving(t *testing.T) {

}

func TestDispatchRemoveSampleOrder(t *testing.T) {

}
