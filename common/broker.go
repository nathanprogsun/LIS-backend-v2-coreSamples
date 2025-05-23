package common

import (
	"encoding/json"

	capi "github.com/hashicorp/consul/api"
)

type KafkaConfiguration struct {
	Address                    []string `json:"address,omitempty"`
	TopicGeneralEvent          string   `json:"topic_general_event,omitempty"`
	TopicPostOrder             string   `json:"topic_post_order,omitempty"`
	TopicCancelOrder           string   `json:"topic_cancel_order,omitempty"`
	TopicRedrawOrder           string   `json:"topic_redraw_info,omitempty"`
	TopicTransactionShipping   string   `json:"topic_transaction_shipping,omitempty"`
	TopicSampleTest            string   `json:"topic_sample_test,omitempty"`
	TopicEmail                 string   `json:"topic_email,omitempty"`
	TopicEditOrder             string   `json:"topic_edit_order,omitempty"`
	TopicHubspot               string   `json:"topic_hubspot,omitempty"`
	GroupIDPostOrder           string   `json:"group_id_post_order,omitempty"`
	GroupIDGeneralEvent        string   `json:"group_id_general_event,omitempty"`
	GroupIDCancelOrder         string   `json:"group_id_cancel_order,omitempty"`
	GroupIDRedrawOrder         string   `json:"group_id_redraw_info,omitempty"`
	GroupIDTransactionShipping string   `json:"group_id_transaction_shipping,omitempty"`
	GroupIDEditOrder           string   `json:"group_id_edit_order,omitempty"`
	GroupIDHubspotEvent        string   `json:"group_id_hubspot_event"`
	// CDC updates
	GroupIDLISCoreCDC                       string `json:"group_id_lis_core_cdc"`
	TopicLISCoreCDCAddress                  string `json:"topic_lis_core_cdc_address,omitempty"`
	TopicLISCoreCDCClinic                   string `json:"topic_lis_core_cdc_clinic,omitempty"`
	TopicLISCoreCDCContact                  string `json:"topic_lis_core_cdc_contact,omitempty"`
	TopicLISCoreCDCCustomer                 string `json:"topic_lis_core_cdc_customer,omitempty"`
	TopicLISCoreCDCInternalUser             string `json:"topic_lis_core_cdc_internal_user,omitempty"`
	TopicLISCoreCDCPatient                  string `json:"topic_lis_core_cdc_patient,omitempty"`
	TopicLISCoreCDCSetting                  string `json:"topic_lis_core_cdc_setting,omitempty"`
	TopicLISCoreCDCUser                     string `json:"topic_lis_core_cdc_user,omitempty"`
	TopicLISCoreCDCCustomerToPatient        string `json:"topic_lis_core_cdc__customertopatient,omitempty"`
	TopicLISCoreCDCCustomerSettingOnClinics string `json:"topic_lis_core_cdc_customersettingonclinics,omitempty"`
	TopicLISCoreCDCClinicToCustomer         string `json:"topic_lis_core_cdc__clinictocustomer,omitempty"`
	TopicLISCoreCDCClinicToPatient          string `json:"topic_lis_core_cdc__clinictopatient,omitempty"`
	TopicLISCoreCDCClinicToSetting          string `json:"topic_lis_core_cdc__clinictosetting,omitempty"`
}

var LocalKafkaConfigs = &KafkaConfiguration{}

func GetLocalKafkaConfigsFromConsul(client *capi.Client, prefix string, key string) {
	val, _, err := client.KV().Get(prefix+"/"+key, nil)
	if err != nil {
		Fatal(err)
	}
	err = json.Unmarshal(val.Value, LocalKafkaConfigs)
	if err != nil {
		Fatal(err)
	}
}
