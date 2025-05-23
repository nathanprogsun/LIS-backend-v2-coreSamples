package tasks

import (
	pb "coresamples/proto"
	"encoding/json"
	"github.com/hibiken/asynq"
)

const (
	TypeAddressCDCUpdate                  = "lis_core_cdc:address"
	TypeClinicCDCUpdate                   = "lis_core_cdc:clinic"
	TypeContactCDCUpdate                  = "lis_core_cdc:contact"
	TypeCustomerCDCUpdate                 = "lis_core_cdc:customer"
	TypeInternalUserCDCUpdate             = "lis_core_cdc:internal_user"
	TypePatientCDCUpdate                  = "lis_core_cdc:patient"
	TypeSettingCDCUpdate                  = "lis_core_cdc:setting"
	TypeUserCDCUpdate                     = "lis_core_cdc:user"
	TypeCustomerToPatientCDCUpdate        = "lis_core_cdc:customer_to_patient"
	TypeCustomerSettingOnClinicsCDCUpdate = "lis_core_cdc:customer_setting_on_clinics"
	TypeClinicToCustomerCDCUpdate         = "lis_core_cdc:clinic_to_customer"
	TypeClinicToPatientCDCUpdate          = "lis_core_cdc:clinic_to_patient"
	TypeClinicToSettingCDCUpdate          = "lis_core_cdc:clinic_to_setting"
)

type AddressCDCUpdateTask struct {
	Event *pb.AddressCDCUpdate
}

type ClinicCDCUpdateTask struct {
	Event *pb.ClinicCDCUpdate
}

type ContactCDCUpdateTask struct {
	Event *pb.ContactCDCUpdate
}

type CustomerCDCUpdateTask struct {
	Event *pb.CustomerCDCUpdate
}

type InternalUserCDCUpdateTask struct {
	Event *pb.InternalUserCDCUpdate
}

type PatientCDCUpdateTask struct {
	Event *pb.PatientCDCUpdate
}

type SettingCDCUpdateTask struct {
	Event *pb.SettingCDCUpdate
}

type UserCDCUpdateTask struct {
	Event *pb.UserCDCUpdate
}

type CustomerToPatientCDCUpdateTask struct {
	Event *pb.CustomerToPatientCDCUpdate
}

type CustomerSettingOnClinicsCDCUpdateTask struct {
	Event *pb.CustomerSettingOnClinicsCDCUpdate
}

type ClinicToCustomerCDCUpdateTask struct {
	Event *pb.ClinicToCustomerCDCUpdate
}

type ClinicToPatientCDCUpdateTask struct {
	Event *pb.ClinicToPatientCDCUpdate
}

type ClinicToSettingCDCUpdateTask struct {
	Event *pb.ClinicToSettingCDCUpdate
}

func NewAddressCDCUpdateTask(task *AddressCDCUpdateTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeAddressCDCUpdate, payload), nil
}

func NewClinicCDCUpdateTask(task *ClinicCDCUpdateTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeClinicCDCUpdate, payload), nil
}

func NewContactCDCUpdateTask(task *ContactCDCUpdateTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeContactCDCUpdate, payload), nil
}

func NewCustomerCDCUpdateTask(task *CustomerCDCUpdateTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeCustomerCDCUpdate, payload), nil
}

func NewInternalUserCDCUpdateTask(task *InternalUserCDCUpdateTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeInternalUserCDCUpdate, payload), nil
}

func NewPatientCDCUpdateTask(task *PatientCDCUpdateTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypePatientCDCUpdate, payload), nil
}

func NewSettingCDCUpdateTask(task *SettingCDCUpdateTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSettingCDCUpdate, payload), nil
}

func NewUserCDCUpdateTask(task *UserCDCUpdateTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeUserCDCUpdate, payload), nil
}

func NewCustomerToPatientCDCUpdateTask(task *CustomerToPatientCDCUpdateTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeCustomerToPatientCDCUpdate, payload), nil
}

func NewCustomerSettingOnClinicsCDCUpdateTask(task *CustomerSettingOnClinicsCDCUpdateTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeCustomerSettingOnClinicsCDCUpdate, payload), nil
}

func NewClinicToCustomerCDCUpdateTask(task *ClinicToCustomerCDCUpdateTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeClinicToCustomerCDCUpdate, payload), nil
}

func NewClinicToPatientCDCUpdateTask(task *ClinicToPatientCDCUpdateTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeClinicToPatientCDCUpdate, payload), nil
}

func NewClinicToSettingCDCUpdateTask(task *ClinicToSettingCDCUpdateTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeClinicToSettingCDCUpdate, payload), nil
}
