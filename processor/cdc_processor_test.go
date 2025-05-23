package processor

import (
	"context"
	"coresamples/common"
	"coresamples/ent"
	"coresamples/ent/customersettingonclinics"
	"coresamples/ent/enttest"
	pb "coresamples/proto"
	"coresamples/publisher"
	"coresamples/tasks"
	"encoding/json"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stvp/tempredis"
)

func setupCDCProcessorTest(t *testing.T) (*CDCProcessor, *tempredis.Server, tasks.AsynqClient) {
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

	processor := &CDCProcessor{
		Processor: InitProcessor(dbClient,
			client,
			context.Background()),
		rs: nil,
	}

	return processor, server, asynqClient
}

func cleanUpCDCProcessorTest(p *CDCProcessor, s *tempredis.Server, asynqClient tasks.AsynqClient) {
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

func TestHandleCustomerSettingOnClinicsUpdates(t *testing.T) {
	processor, server, asynqClient := setupCDCProcessorTest(t)
	defer cleanUpCDCProcessorTest(processor, server, asynqClient)

	testKafkaMessage := []byte(`{
		"database": "lis_core_v7",
		"table": "customersettingonclinics",
		"type": "update",
		"ts": 1739321116,
		"xid": 13925989445,
		"commit": true,
		"data": {
			"customer_id" : 1009,
			"clinic_id": 71557,
			"setting_id": 2,
			"setting_name": "test_setting1"
		},
		"old": {
			"setting_id": 1
		}
	}`)

	event := &pb.CustomerSettingOnClinicsCDCUpdate{}
	err := json.Unmarshal(testKafkaMessage, event)
	if err != nil {
		t.Fatalf("failed to unmarshal message: %v", err)
	}

	ctx := context.Background()
	_, err = processor.dbClient.ExecContext(ctx, "PRAGMA foreign_keys=OFF")
	if err != nil {
		t.Fatal(err)
	}

	err = processor.dbClient.CustomerSettingOnClinics.Create().
		SetCustomerID(1009).
		SetClinicID(71557).
		SetSettingID(1).SetSettingName("test_setting1").
		Exec(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = processor.dbClient.CustomerSettingOnClinics.Create().
		SetCustomerID(1009).
		SetClinicID(71557).
		SetSettingID(3).SetSettingName("test_setting3").
		Exec(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = processor.dbClient.CustomerSettingOnClinics.Create().
		SetCustomerID(int(event.Data.CustomerId)).
		SetClinicID(int(event.Data.ClinicId)).
		SetSettingID(int(event.Data.SettingId)).
		SetSettingName(event.Data.SettingName).
		OnConflictColumns("customer_id", "clinic_id", "setting_name").
		Update(func(u *ent.CustomerSettingOnClinicsUpsert) {
			u.SetSettingID(int(event.Data.SettingId)) // ✅ Only update setting_name if conflict occurs
		}).
		Exec(ctx)

	if err != nil {
		t.Fatal(err)
	}

	updatedSetting, err := processor.dbClient.CustomerSettingOnClinics.Query().
		Where(
			customersettingonclinics.CustomerIDEQ(int(event.Data.CustomerId)),
			customersettingonclinics.ClinicIDEQ(int(event.Data.ClinicId)),
			customersettingonclinics.SettingNameEQ(event.Data.SettingName),
		).First(ctx)
	if err != nil {
		t.Fatalf("Failed to query updated customer setting: %v", err)
	}

	if updatedSetting.SettingID != 2 {
		t.Fatalf("Expected SettingID to be updated to 2, got %d", updatedSetting.SettingID)
	}

	unmodifiedSetting, err := processor.dbClient.CustomerSettingOnClinics.Query().
		Where(
			customersettingonclinics.CustomerIDEQ(int(event.Data.CustomerId)),
			customersettingonclinics.ClinicIDEQ(int(event.Data.ClinicId)),
			customersettingonclinics.SettingNameEQ("test_setting3"),
		).First(ctx)
	if err != nil {
		t.Fatalf("Failed to query the testing setting for customer clinic: %v", err)
	}

	if unmodifiedSetting.SettingID != 3 {
		t.Fatalf("Expected SettingID to be unmodified as 3, got %d", updatedSetting.SettingID)
	}

	_, err = processor.dbClient.ExecContext(ctx, "PRAGMA foreign_keys=ON")
	if err != nil {
		t.Fatal(err)
	}
}
