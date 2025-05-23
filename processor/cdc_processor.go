package processor

import (
	"context"
	"coresamples/common"
	"coresamples/dbutils"
	"coresamples/ent"
	"coresamples/tasks"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-redsync/redsync/v4"
	"github.com/hibiken/asynq"
)

type CDCProcessor struct {
	Processor
	rs *redsync.Redsync
}

func NewCDCProcessor(dbClient *ent.Client, redisClient *common.RedisClient, rs *redsync.Redsync) *CDCProcessor {
	return &CDCProcessor{
		Processor: InitProcessor(dbClient, redisClient, context.Background()),
		rs:        rs,
	}
}

func (h *CDCProcessor) HandleAddressUpdates(ctx context.Context, t *asynq.Task) error {
	task := &tasks.AddressCDCUpdateTask{}
	if err := json.Unmarshal(t.Payload(), task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	event := task.Event
	if event == nil {
		return fmt.Errorf("nil event in task :%w", asynq.SkipRetry)
	}

	if common.Env.DryRun {
		return nil
	}

	addressID := int(event.Data.AddressId)

	switch strings.ToLower(event.Type) {
	case "insert", "update":
		return dbutils.UpsertAddress(event.Data, h.dbClient, ctx)
	case "delete":
		return dbutils.DeleteAddress(addressID, h.dbClient, ctx)
	default:
		return fmt.Errorf("unknown cdc update type %s %w", event.Type, asynq.SkipRetry)
	}
}

func (h *CDCProcessor) HandleClinicUpdates(ctx context.Context, t *asynq.Task) error {
	if t == nil {
		return fmt.Errorf("t is nil")
	}
	task := &tasks.ClinicCDCUpdateTask{}
	if err := json.Unmarshal(t.Payload(), task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	event := task.Event
	if event == nil {
		return fmt.Errorf("nil event in task :%w", asynq.SkipRetry)
	}

	if common.Env.DryRun {
		return nil
	}

	switch strings.ToLower(event.Type) {
	case "insert", "update":
		return dbutils.UpsertClinic(ctx, event.Data, h.dbClient)
	case "delete":
		return dbutils.DeleteClinic(ctx, event.Data.ClinicId, h.dbClient)
	default:
		return fmt.Errorf("unknown cdc update type %s %w", event.Type, asynq.SkipRetry)
	}
}

func (h *CDCProcessor) HandleContactUpdates(ctx context.Context, t *asynq.Task) error {
	task := &tasks.ContactCDCUpdateTask{}
	if err := json.Unmarshal(t.Payload(), task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	event := task.Event
	if event == nil {
		return fmt.Errorf("nil event in task :%w", asynq.SkipRetry)
	}

	if common.Env.DryRun {
		return nil
	}

	if common.Env.DryRun {
		return nil
	}

	switch strings.ToLower(event.Type) {
	case "insert", "update":
		return dbutils.UpsertContactCDC(event.Data, h.dbClient, ctx)
	case "delete":
		return dbutils.DeleteContactCDC(event.Data, h.dbClient, ctx)
	default:
		return fmt.Errorf("unknown cdc update type %s %w", event.Type, asynq.SkipRetry)
	}
}

func (h *CDCProcessor) HandleCustomerUpdates(ctx context.Context, t *asynq.Task) error {
	task := &tasks.CustomerCDCUpdateTask{}
	if err := json.Unmarshal(t.Payload(), task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	event := task.Event
	if event == nil {
		return fmt.Errorf("nil event in task :%w", asynq.SkipRetry)
	}

	if common.Env.DryRun {
		return nil
	}

	switch strings.ToLower(event.Type) {
	case "insert", "update":
		return dbutils.UpsertCustomerCDC(event.Data, h.dbClient, ctx)
	case "delete":
		return dbutils.DeleteCustomerCDC(event.Data, h.dbClient, ctx)
	default:
		return fmt.Errorf("unknown cdc update type %s %w", event.Type, asynq.SkipRetry)
	}
}

func (h *CDCProcessor) HandleInternalUserUpdates(ctx context.Context, t *asynq.Task) error {
	task := &tasks.InternalUserCDCUpdateTask{}
	if err := json.Unmarshal(t.Payload(), task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	event := task.Event
	if event == nil {
		return fmt.Errorf("nil event in task :%w", asynq.SkipRetry)
	}

	// if common.Env.DryRun {
	// 	return nil
	// }

	switch strings.ToLower(event.Type) {
	case "insert", "update":
		return dbutils.UpsertInternalUserCDC(event.Data, h.dbClient, ctx)
	case "delete":
		return dbutils.DeleteInternalUserCDC(event.Data, h.dbClient, ctx)
	default:
		return fmt.Errorf("unknown cdc update type %s %w", event.Type, asynq.SkipRetry)
	}
}

func (h *CDCProcessor) HandlePatientUpdates(ctx context.Context, t *asynq.Task) error {
	task := &tasks.PatientCDCUpdateTask{}
	if err := json.Unmarshal(t.Payload(), task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	event := task.Event
	if event == nil {
		return fmt.Errorf("nil event in task :%w", asynq.SkipRetry)
	}

	if common.Env.DryRun {
		return nil
	}

	switch strings.ToLower(event.Type) {
	case "insert", "update":
		return dbutils.UpsertPatient(ctx, event.Data, h.dbClient)
	case "delete":
		return dbutils.DeletePatient(ctx, event.Data.PatientId, h.dbClient)
	default:
		return fmt.Errorf("unknown cdc update type %s %w", event.Type, asynq.SkipRetry)
	}
}

func (h *CDCProcessor) HandleSettingUpdates(ctx context.Context, t *asynq.Task) error {
	task := &tasks.SettingCDCUpdateTask{}
	if err := json.Unmarshal(t.Payload(), task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	event := task.Event
	if event == nil {
		return fmt.Errorf("event data is missing in task payload: %w", asynq.SkipRetry)
	}

	if common.Env.DryRun {
		return nil
	}

	settingID := int(event.Data.SettingId)

	switch strings.ToLower(event.Type) {
	case "insert", "update":
		return dbutils.UpsertSetting(event.Data, h.dbClient, ctx)

	case "delete":
		return dbutils.DeleteSetting(settingID, h.dbClient, ctx)

	default:
		return fmt.Errorf("unknown operation type: %s: %w", event.Type, asynq.SkipRetry)
	}
}

func (h *CDCProcessor) HandleUserUpdates(ctx context.Context, t *asynq.Task) error {
	task := &tasks.UserCDCUpdateTask{}
	if err := json.Unmarshal(t.Payload(), task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	event := task.Event
	if event == nil {
		return fmt.Errorf("nil event in task :%w", asynq.SkipRetry)
	}

	if common.Env.DryRun {
		return nil
	}

	switch strings.ToLower(event.Type) {
	case "insert", "update":
		return dbutils.UpsertUserCDC(event.Data, h.dbClient, ctx)
	case "delete":
		return dbutils.DeleteUserCDC(event.Data, h.dbClient, ctx)
	default:
		return fmt.Errorf("unknown cdc update type %s %w", event.Type, asynq.SkipRetry)
	}
}

func (h *CDCProcessor) HandleCustomerToPatientUpdates(ctx context.Context, t *asynq.Task) error {
	task := &tasks.CustomerToPatientCDCUpdateTask{}
	if err := json.Unmarshal(t.Payload(), task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	event := task.Event
	if event == nil {
		return fmt.Errorf("nil event in task :%w", asynq.SkipRetry)
	}

	// if common.Env.DryRun {
	// 	return nil
	// }

	switch strings.ToLower(event.Type) {
	case "insert":
		return dbutils.CreateCustomerToPatientCDC(event.Data, h.dbClient, ctx)
	case "delete":
		return dbutils.DeleteCustomerToPatientCDC(event.Data, h.dbClient, ctx)
	case "update":
		if event.Old.B != 0 {
			return dbutils.UpdateCustomerToPatientCDC(event.Data, event.Old, h.dbClient, ctx)
		} else if event.Old.A != 0 {
			return dbutils.UpdatePatientToCustomerCDC(event.Data, event.Old, h.dbClient, ctx)
		} else {
			return fmt.Errorf("invalid cdc update data %s %w", event.Type, asynq.SkipRetry)
		}
	default:
		return fmt.Errorf("unknown cdc update type %s %w", event.Type, asynq.SkipRetry)
	}
}

func (h *CDCProcessor) HandleCustomerSettingOnClinicsUpdates(ctx context.Context, t *asynq.Task) error {
	// Parse task payload
	task := &tasks.CustomerSettingOnClinicsCDCUpdateTask{}
	if err := json.Unmarshal(t.Payload(), task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// Validate event data
	event := task.Event
	if event == nil {
		return fmt.Errorf("event data is missing in task payload: %w", asynq.SkipRetry)
	}

	if common.Env.DryRun {
		return nil
	}

	customerID := int(event.Data.CustomerId)
	clinicID := int(event.Data.ClinicId)
	settingID := int(event.Data.SettingId)
	settingName := event.Data.SettingName

	switch strings.ToLower(event.Type) {
	case "insert", "update":
		// Upsert customer setting on clinic (insert or update)
		return dbutils.UpsertCustomerSettingOnClinics(customerID, clinicID, settingID, settingName, h.dbClient, ctx)

	case "delete":
		// Remove customer setting on clinic
		return dbutils.RemoveCustomerSettingOnClinics(customerID, clinicID, settingID, h.dbClient, ctx)

	default:
		return fmt.Errorf("unknown operation type: %s: %w", event.Type, asynq.SkipRetry)
	}
}

func (h *CDCProcessor) HandleClinicToCustomerUpdates(ctx context.Context, t *asynq.Task) error {
	task := &tasks.ClinicToCustomerCDCUpdateTask{}
	if err := json.Unmarshal(t.Payload(), task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	event := task.Event
	if event == nil {
		return fmt.Errorf("nil event in task :%w", asynq.SkipRetry)
	}

	if common.Env.DryRun {
		return nil
	}

	switch strings.ToLower(event.Type) {
	case "insert":
		return dbutils.CreateClinicToCustomerCDC(event.Data, h.dbClient, ctx)
	case "delete":
		return dbutils.DeleteClinicToCustomerCDC(event.Data, h.dbClient, ctx)
	case "update":
		if event.Old.B != 0 {
			return dbutils.UpdateClinicToCustomerCDC(event.Data, event.Old, h.dbClient, ctx)
		} else if event.Old.A != 0 {
			return dbutils.UpdateCustomerToClinicCDC(event.Data, event.Old, h.dbClient, ctx)
		} else {
			return fmt.Errorf("invalid cdc update data %s %w", event.Type, asynq.SkipRetry)
		}
	default:
		return fmt.Errorf("unknown cdc update type %s %w", event.Type, asynq.SkipRetry)
	}
}

func (h *CDCProcessor) HandleClinicToPatientUpdates(ctx context.Context, t *asynq.Task) error {
	task := &tasks.ClinicToPatientCDCUpdateTask{}
	if err := json.Unmarshal(t.Payload(), task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	event := task.Event
	if event == nil {
		return fmt.Errorf("nil event in task :%w", asynq.SkipRetry)
	}

	if common.Env.DryRun {
		return nil
	}

	switch strings.ToLower(event.Type) {
	case "insert":
		return dbutils.CreateClinicToPatientCDC(event.Data, h.dbClient, ctx)
	case "delete":
		return dbutils.DeleteClinicToPatientCDC(event.Data, h.dbClient, ctx)
	case "update":
		return dbutils.UpdateClinicToPatientCDC(event.Data, event.Old, h.dbClient, ctx)
	default:
		return fmt.Errorf("unknown cdc update type %s %w", event.Type, asynq.SkipRetry)
	}
}

func (h *CDCProcessor) HandleClinicToSettingUpdates(ctx context.Context, t *asynq.Task) error {
	task := &tasks.ClinicToSettingCDCUpdateTask{}
	if err := json.Unmarshal(t.Payload(), task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	event := task.Event
	if event == nil {
		return fmt.Errorf("event data is missing in task payload: %w", asynq.SkipRetry)
	}

	if common.Env.DryRun {
		return nil
	}

	switch strings.ToLower(event.Type) {
	case "insert":
		// Handle insertion of a new clinic-to-setting relationship
		return dbutils.AddSettingToClinic(int(event.Data.B), int(event.Data.A), h.dbClient, ctx)

	case "update":
		// Handle updating a clinic's setting
		return dbutils.UpdateSettingToClinic(int(event.Data.B), int(event.Old.B), int(event.Data.A), h.dbClient, ctx)

	case "delete":
		return dbutils.RemoveSettingFromClinic(int(event.Data.B), int(event.Data.A), h.dbClient, ctx)

	default:
		// Log and ignore unknown operation types
		return fmt.Errorf("unknown operation type: %s: %w", event.Type, asynq.SkipRetry)
	}
}
