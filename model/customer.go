package model

type NPIApiResponse struct {
	ResultCount int              `json:"result_count,omitempty"`
	Results     []NPIResult      `json:"results,omitempty"`
	Errors      []NPIErrorDetail `json:"Errors,omitempty"`
}

type NPIErrorDetail struct {
	Description string `json:"description"`
	Field       string `json:"field"`
	Number      string `json:"number"`
}

type NPIResult struct {
	CreatedEpoch      string        `json:"created_epoch"`
	EnumerationType   string        `json:"enumeration_type"`
	LastUpdatedEpoch  string        `json:"last_updated_epoch"`
	Number            string        `json:"number"`
	Addresses         []NPIAddress  `json:"addresses"`
	PracticeLocations []interface{} `json:"practiceLocations"` // Replace with actual type if needed
	Basic             NPIBasic      `json:"basic"`
	Taxonomies        []NPITaxonomy `json:"taxonomies"`
	Identifiers       []interface{} `json:"identifiers"` // Replace if needed
	Endpoints         []interface{} `json:"endpoints"`   // Replace if needed
	OtherNames        []interface{} `json:"other_names"` // Replace if needed
}

type NPIAddress struct {
	CountryCode     string `json:"country_code"`
	CountryName     string `json:"country_name"`
	AddressPurpose  string `json:"address_purpose"`
	AddressType     string `json:"address_type"`
	Address1        string `json:"address_1"`
	Address2        string `json:"address_2"`
	City            string `json:"city"`
	State           string `json:"state"`
	PostalCode      string `json:"postal_code"`
	TelephoneNumber string `json:"telephone_number"`
	FaxNumber       string `json:"fax_number"`
}

type NPIBasic struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	MiddleName      string `json:"middle_name"`
	Credential      string `json:"credential"`
	SoleProprietor  string `json:"sole_proprietor"`
	Gender          string `json:"gender"`
	EnumerationDate string `json:"enumeration_date"`
	LastUpdated     string `json:"last_updated"`
	Status          string `json:"status"`
	NamePrefix      string `json:"name_prefix"`
	NameSuffix      string `json:"name_suffix"`
}

type NPITaxonomy struct {
	Code          string `json:"code"`
	TaxonomyGroup string `json:"taxonomy_group"`
	Desc          string `json:"desc"`
	State         string `json:"state"`
	License       string `json:"license"`
	Primary       bool   `json:"primary"`
}

type FullCustomer struct {
	CustomerID                int32             `json:"customer_id,omitempty"`
	UserID                    int32             `json:"user_id,omitempty"`
	CustomerFirstName         string            `json:"customer_first_name,omitempty"`
	CustomerLastName          string            `json:"customer_last_name,omitempty"`
	CustomerMiddleName        string            `json:"customer_middle_name,omitempty"`
	CustomerTypeID            string            `json:"customer_type_id,omitempty"`
	CustomerSuffix            string            `json:"customer_suffix,omitempty"`
	CustomerSamplesReceived   string            `json:"customer_samples_received,omitempty"`
	CustomerRequestSubmitTime string            `json:"customer_request_submit_time,omitempty"`
	PaymentMethod             string            `json:"payment_method,omitempty"`
	IsActive                  bool              `json:"isActive,omitempty"`
	Clinics                   []*CustomerClinic `json:"clinics,omitempty"`
	CustomerNPINumber         string            `json:"customer_npi_number,omitempty"`
	SalesID                   int32             `json:"sales_id,omitempty"`
	CustomerSignupTime        string            `json:"customer_signup_time,omitempty"`
}

type CustomerClinic struct {
	ClinicID          int32      `json:"clinic_id,omitempty"`
	ClinicName        string     `json:"clinic_name,omitempty"`
	UserID            int32      `json:"user_id,omitempty"`
	ClinicType        string     `json:"clinic_type,omitempty"`
	IsActive          bool       `json:"isActive,omitempty"`
	ClinicNPINumber   string     `json:"clinic_npi_number,omitempty"`
	ClinicAccountID   int32      `json:"clinic_account_id,omitempty"`
	CustomerAddresses []*Address `json:"customer_addresses,omitempty"`
	CustomerContacts  []*Contact `json:"customer_contacts,omitempty"`
}

type CustomerSales struct {
	CustomerId         int32         `json:"customer_id,omitempty"`
	CustomerFirstName  string        `json:"customer_first_name,omitempty"`
	CustomerLastName   string        `json:"customer_last_name,omitempty"`
	CustomerMiddleName string        `json:"customer_middle_name,omitempty"`
	InternalUser       *InternalUser `json:"internal_user,omitempty"`
}

type InternalUser struct {
	InternalUserRoleId     int32  `json:"internal_user_role_id,omitempty"`
	InternalUserFirstname  string `json:"internal_user_firstname,omitempty"`
	InternalUserLastname   string `json:"internal_user_lastname,omitempty"`
	InternalUserMiddlename string `json:"internal_user_middlename,omitempty"`
	InternalUserEmail      string `json:"internal_user_email,omitempty"`
	InternalUserPhone      string `json:"internal_user_phone,omitempty"`
}

type FuzzyClientObject struct {
	ClientId   int64  `json:"client_id,omitempty"`
	ClientName string `json:"client_name,omitempty"`
}

type CustomerClinicData struct {
	CustomerId   int32  `json:"customer_id,omitempty"`
	CustomerName string `json:"customer_name,omitempty"`
	ClinicName   string `pjson:"clinic_name,omitempty"`
}

type AddCustomerWithNPINumberRequest struct {
	// Customer Basic Info
	CustomerFirstName         string `json:"customer_first_name,omitempty"`
	CustomerLastName          string `json:"customer_last_name,omitempty"`
	CustomerNPINumber         string `json:"customer_npi_number,omitempty"`
	CustomerLoginEmail        string `json:"customer_login_email,omitempty"`
	CustomerNotificationEmail string `json:"customer_notification_email,omitempty"`
	CustomerPhone             string `json:"customer_phone,omitempty"`

	// Customer Address
	CustomerAddressLine1 string `json:"customer_address_line_1,omitempty"`
	CustomerAddressLine2 string `json:"customer_address_line_2,omitempty"`
	CustomerCity         string `json:"customer_city,omitempty"`
	CustomerState        string `json:"customer_state,omitempty"`
	CustomerZipcode      string `json:"customer_zipcode,omitempty"`
	CustomerCountry      string `json:"customer_country,omitempty"`

	// Customer Role in Clinic
	// CustomerRole string `json:"customer_role,omitempty"`
	ClinicID string `json:"clinic_id,omitempty"`

	// Invitation info
	InvitedFromCustomer    string   `json:"invited_from_customer,omitempty"`
	CustomerInvitationLink string   `json:"customer_invitation_link,omitempty"`
	CustomerSuffix         string   `json:"customer_suffix,omitempty"`
	CustomerRoles          []string `json:"customer_roles,omitempty"`
}

type AddCustomerWithNPINumberResponse struct {
	Status       string `json:"status,omitempty"`
	CustomerID   int32  `json:"customer_id,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type CustomerBetaPrograms struct {
	CustomerID   int32    `json:"customer_id,omitempty"`
	ClinicID     int32    `json:"clinic_id,omitempty"`
	BetaPrograms []string `json:"beta_programs,omitempty"`
}
