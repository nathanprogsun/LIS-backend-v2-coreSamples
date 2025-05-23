package model

type Contact struct {
	ContactID          int32  `json:"contact_id,omitempty"`
	ContactDescription string `json:"contact_description,omitempty"`
	ContactDetails     string `json:"contact_details,omitempty"`
	ContactType        string `json:"contact_type,omitempty"`
	IsPrimaryContact   bool   `json:"is_primary_contact,omitempty"`
	ContactLevel       int32  `json:"contact_level,omitempty"`
	ContactLevelName   string `json:"contact_level_name,omitempty"`
}

type CustomerContactOnClinicsCreation struct {
	CustomerID  int32  `json:"customer_id,omitempty"`
	ClinicID    int32  `json:"clinic_id,omitempty"`
	ContactID   int32  `json:"contact_id,omitempty"`
	ContactType string `json:"contact_type,omitempty"`
}
