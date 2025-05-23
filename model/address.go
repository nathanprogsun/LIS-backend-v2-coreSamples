package model

type Address struct {
	AddressID        int32  `json:"address_id,omitempty"`
	AddressType      string `json:"address_type,omitempty"`
	StreetAddress    string `json:"street_address,omitempty"`
	AptPO            string `json:"apt_po,omitempty"`
	City             string `json:"city,omitempty"`
	State            string `json:"state,omitempty"`
	Zipcode          string `json:"zipcode,omitempty"`
	Country          string `json:"country,omitempty"`
	AddressConfirmed bool   `json:"address_confirmed,omitempty"`
	IsPrimaryAddress bool   `json:"is_primary_address,omitempty"`
	AddressLevel     int32  `json:"address_level,omitempty"`
	AddressLevelName string `json:"address_level_name,omitempty"`
}

type CustomerAddressOnClinicsCreation struct {
	CustomerID  int32  `json:"customer_id,omitempty"`
	ClinicID    int32  `json:"clinic_id,omitempty"`
	AddressID   int32  `json:"address_id,omitempty"`
	AddressType string `json:"address_type,omitempty"`
}
