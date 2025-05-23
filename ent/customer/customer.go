// Code generated by ent, DO NOT EDIT.

package customer

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the customer type in the database.
	Label = "customer"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "customer_id"
	// FieldUserID holds the string denoting the user_id field in the database.
	FieldUserID = "user_id"
	// FieldCustomerType holds the string denoting the customer_type field in the database.
	FieldCustomerType = "customer_type"
	// FieldCustomerFirstName holds the string denoting the customer_first_name field in the database.
	FieldCustomerFirstName = "customer_first_name"
	// FieldCustomerLastName holds the string denoting the customer_last_name field in the database.
	FieldCustomerLastName = "customer_last_name"
	// FieldCustomerMiddleName holds the string denoting the customer_middle_name field in the database.
	FieldCustomerMiddleName = "customer_middle_name"
	// FieldCustomerTypeID holds the string denoting the customer_type_id field in the database.
	FieldCustomerTypeID = "customer_type_id"
	// FieldCustomerSuffix holds the string denoting the customer_suffix field in the database.
	FieldCustomerSuffix = "customer_suffix"
	// FieldCustomerSamplesReceived holds the string denoting the customer_samples_received field in the database.
	FieldCustomerSamplesReceived = "customer_samples_received"
	// FieldCustomerRequestSubmitTime holds the string denoting the customer_request_submit_time field in the database.
	FieldCustomerRequestSubmitTime = "customer_request_submit_time"
	// FieldCustomerSignupTime holds the string denoting the customer_signup_time field in the database.
	FieldCustomerSignupTime = "customer_signup_time"
	// FieldIsActive holds the string denoting the is_active field in the database.
	FieldIsActive = "isActive"
	// FieldSalesID holds the string denoting the sales_id field in the database.
	FieldSalesID = "sales_id"
	// FieldCustomerNpiNumber holds the string denoting the customer_npi_number field in the database.
	FieldCustomerNpiNumber = "customer_npi_number"
	// FieldReferralSource holds the string denoting the referral_source field in the database.
	FieldReferralSource = "referral_source"
	// FieldOrderPlacementAllowed holds the string denoting the order_placement_allowed field in the database.
	FieldOrderPlacementAllowed = "order_placement_allowed"
	// FieldBetaProgramEnabled holds the string denoting the beta_program_enabled field in the database.
	FieldBetaProgramEnabled = "beta_program_enabled"
	// FieldOnboardingQuestionnaireFilledOn holds the string denoting the onboarding_questionnaire_filled_on field in the database.
	FieldOnboardingQuestionnaireFilledOn = "onboarding_questionnaire_filled_on"
	// EdgeSamples holds the string denoting the samples edge name in mutations.
	EdgeSamples = "samples"
	// EdgeCustomerContacts holds the string denoting the customer_contacts edge name in mutations.
	EdgeCustomerContacts = "customer_contacts"
	// EdgeCustomerAddresses holds the string denoting the customer_addresses edge name in mutations.
	EdgeCustomerAddresses = "customer_addresses"
	// EdgeClinics holds the string denoting the clinics edge name in mutations.
	EdgeClinics = "clinics"
	// EdgeSales holds the string denoting the sales edge name in mutations.
	EdgeSales = "sales"
	// EdgeUser holds the string denoting the user edge name in mutations.
	EdgeUser = "user"
	// EdgeOrders holds the string denoting the orders edge name in mutations.
	EdgeOrders = "orders"
	// EdgeCurrentPatients holds the string denoting the current_patients edge name in mutations.
	EdgeCurrentPatients = "current_patients"
	// EdgePatients holds the string denoting the patients edge name in mutations.
	EdgePatients = "patients"
	// EdgeCustomerBetaProgramParticipations holds the string denoting the customer_beta_program_participations edge name in mutations.
	EdgeCustomerBetaProgramParticipations = "customer_beta_program_participations"
	// EdgeCustomerSettingsOnClinics holds the string denoting the customer_settings_on_clinics edge name in mutations.
	EdgeCustomerSettingsOnClinics = "customer_settings_on_clinics"
	// EdgeCustomerAddressesOnClinics holds the string denoting the customer_addresses_on_clinics edge name in mutations.
	EdgeCustomerAddressesOnClinics = "customer_addresses_on_clinics"
	// EdgeCustomerContactsOnClinics holds the string denoting the customer_contacts_on_clinics edge name in mutations.
	EdgeCustomerContactsOnClinics = "customer_contacts_on_clinics"
	// SampleFieldID holds the string denoting the ID field of the Sample.
	SampleFieldID = "sample_id"
	// ContactFieldID holds the string denoting the ID field of the Contact.
	ContactFieldID = "contact_id"
	// AddressFieldID holds the string denoting the ID field of the Address.
	AddressFieldID = "address_id"
	// ClinicFieldID holds the string denoting the ID field of the Clinic.
	ClinicFieldID = "clinic_id"
	// InternalUserFieldID holds the string denoting the ID field of the InternalUser.
	InternalUserFieldID = "internal_user_id"
	// UserFieldID holds the string denoting the ID field of the User.
	UserFieldID = "user_id"
	// OrderInfoFieldID holds the string denoting the ID field of the OrderInfo.
	OrderInfoFieldID = "order_id"
	// PatientFieldID holds the string denoting the ID field of the Patient.
	PatientFieldID = "patient_id"
	// BetaProgramParticipationFieldID holds the string denoting the ID field of the BetaProgramParticipation.
	BetaProgramParticipationFieldID = "id"
	// CustomerSettingOnClinicsFieldID holds the string denoting the ID field of the CustomerSettingOnClinics.
	CustomerSettingOnClinicsFieldID = "id"
	// CustomerAddressOnClinicsFieldID holds the string denoting the ID field of the CustomerAddressOnClinics.
	CustomerAddressOnClinicsFieldID = "id"
	// CustomerContactOnClinicsFieldID holds the string denoting the ID field of the CustomerContactOnClinics.
	CustomerContactOnClinicsFieldID = "id"
	// Table holds the table name of the customer in the database.
	Table = "customer"
	// SamplesTable is the table that holds the samples relation/edge.
	SamplesTable = "sample"
	// SamplesInverseTable is the table name for the Sample entity.
	// It exists in this package in order to avoid circular dependency with the "sample" package.
	SamplesInverseTable = "sample"
	// SamplesColumn is the table column denoting the samples relation/edge.
	SamplesColumn = "customer_id"
	// CustomerContactsTable is the table that holds the customer_contacts relation/edge.
	CustomerContactsTable = "contact"
	// CustomerContactsInverseTable is the table name for the Contact entity.
	// It exists in this package in order to avoid circular dependency with the "contact" package.
	CustomerContactsInverseTable = "contact"
	// CustomerContactsColumn is the table column denoting the customer_contacts relation/edge.
	CustomerContactsColumn = "customer_id"
	// CustomerAddressesTable is the table that holds the customer_addresses relation/edge.
	CustomerAddressesTable = "address"
	// CustomerAddressesInverseTable is the table name for the Address entity.
	// It exists in this package in order to avoid circular dependency with the "address" package.
	CustomerAddressesInverseTable = "address"
	// CustomerAddressesColumn is the table column denoting the customer_addresses relation/edge.
	CustomerAddressesColumn = "customer_id"
	// ClinicsTable is the table that holds the clinics relation/edge. The primary key declared below.
	ClinicsTable = "clinic_customers"
	// ClinicsInverseTable is the table name for the Clinic entity.
	// It exists in this package in order to avoid circular dependency with the "clinic" package.
	ClinicsInverseTable = "clinic"
	// SalesTable is the table that holds the sales relation/edge.
	SalesTable = "customer"
	// SalesInverseTable is the table name for the InternalUser entity.
	// It exists in this package in order to avoid circular dependency with the "internaluser" package.
	SalesInverseTable = "internal_user"
	// SalesColumn is the table column denoting the sales relation/edge.
	SalesColumn = "sales_id"
	// UserTable is the table that holds the user relation/edge.
	UserTable = "customer"
	// UserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserInverseTable = "user"
	// UserColumn is the table column denoting the user relation/edge.
	UserColumn = "user_id"
	// OrdersTable is the table that holds the orders relation/edge.
	OrdersTable = "order_info"
	// OrdersInverseTable is the table name for the OrderInfo entity.
	// It exists in this package in order to avoid circular dependency with the "orderinfo" package.
	OrdersInverseTable = "order_info"
	// OrdersColumn is the table column denoting the orders relation/edge.
	OrdersColumn = "customer_id"
	// CurrentPatientsTable is the table that holds the current_patients relation/edge.
	CurrentPatientsTable = "patient"
	// CurrentPatientsInverseTable is the table name for the Patient entity.
	// It exists in this package in order to avoid circular dependency with the "patient" package.
	CurrentPatientsInverseTable = "patient"
	// CurrentPatientsColumn is the table column denoting the current_patients relation/edge.
	CurrentPatientsColumn = "customer_id"
	// PatientsTable is the table that holds the patients relation/edge. The primary key declared below.
	PatientsTable = "customer_patients"
	// PatientsInverseTable is the table name for the Patient entity.
	// It exists in this package in order to avoid circular dependency with the "patient" package.
	PatientsInverseTable = "patient"
	// CustomerBetaProgramParticipationsTable is the table that holds the customer_beta_program_participations relation/edge.
	CustomerBetaProgramParticipationsTable = "beta_program_participations"
	// CustomerBetaProgramParticipationsInverseTable is the table name for the BetaProgramParticipation entity.
	// It exists in this package in order to avoid circular dependency with the "betaprogramparticipation" package.
	CustomerBetaProgramParticipationsInverseTable = "beta_program_participations"
	// CustomerBetaProgramParticipationsColumn is the table column denoting the customer_beta_program_participations relation/edge.
	CustomerBetaProgramParticipationsColumn = "customer_id"
	// CustomerSettingsOnClinicsTable is the table that holds the customer_settings_on_clinics relation/edge.
	CustomerSettingsOnClinicsTable = "customer_setting_on_clinics"
	// CustomerSettingsOnClinicsInverseTable is the table name for the CustomerSettingOnClinics entity.
	// It exists in this package in order to avoid circular dependency with the "customersettingonclinics" package.
	CustomerSettingsOnClinicsInverseTable = "customer_setting_on_clinics"
	// CustomerSettingsOnClinicsColumn is the table column denoting the customer_settings_on_clinics relation/edge.
	CustomerSettingsOnClinicsColumn = "customer_id"
	// CustomerAddressesOnClinicsTable is the table that holds the customer_addresses_on_clinics relation/edge.
	CustomerAddressesOnClinicsTable = "customer_address_on_clinics"
	// CustomerAddressesOnClinicsInverseTable is the table name for the CustomerAddressOnClinics entity.
	// It exists in this package in order to avoid circular dependency with the "customeraddressonclinics" package.
	CustomerAddressesOnClinicsInverseTable = "customer_address_on_clinics"
	// CustomerAddressesOnClinicsColumn is the table column denoting the customer_addresses_on_clinics relation/edge.
	CustomerAddressesOnClinicsColumn = "customer_id"
	// CustomerContactsOnClinicsTable is the table that holds the customer_contacts_on_clinics relation/edge.
	CustomerContactsOnClinicsTable = "customer_contact_on_clinics"
	// CustomerContactsOnClinicsInverseTable is the table name for the CustomerContactOnClinics entity.
	// It exists in this package in order to avoid circular dependency with the "customercontactonclinics" package.
	CustomerContactsOnClinicsInverseTable = "customer_contact_on_clinics"
	// CustomerContactsOnClinicsColumn is the table column denoting the customer_contacts_on_clinics relation/edge.
	CustomerContactsOnClinicsColumn = "customer_id"
)

// Columns holds all SQL columns for customer fields.
var Columns = []string{
	FieldID,
	FieldUserID,
	FieldCustomerType,
	FieldCustomerFirstName,
	FieldCustomerLastName,
	FieldCustomerMiddleName,
	FieldCustomerTypeID,
	FieldCustomerSuffix,
	FieldCustomerSamplesReceived,
	FieldCustomerRequestSubmitTime,
	FieldCustomerSignupTime,
	FieldIsActive,
	FieldSalesID,
	FieldCustomerNpiNumber,
	FieldReferralSource,
	FieldOrderPlacementAllowed,
	FieldBetaProgramEnabled,
	FieldOnboardingQuestionnaireFilledOn,
}

var (
	// ClinicsPrimaryKey and ClinicsColumn2 are the table columns denoting the
	// primary key for the clinics relation (M2M).
	ClinicsPrimaryKey = []string{"clinic_id", "customer_id"}
	// PatientsPrimaryKey and PatientsColumn2 are the table columns denoting the
	// primary key for the patients relation (M2M).
	PatientsPrimaryKey = []string{"customer_id", "patient_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCustomerType holds the default value on creation for the "customer_type" field.
	DefaultCustomerType string
	// DefaultCustomerSignupTime holds the default value on creation for the "customer_signup_time" field.
	DefaultCustomerSignupTime func() time.Time
	// DefaultIsActive holds the default value on creation for the "is_active" field.
	DefaultIsActive bool
	// DefaultOrderPlacementAllowed holds the default value on creation for the "order_placement_allowed" field.
	DefaultOrderPlacementAllowed bool
	// DefaultBetaProgramEnabled holds the default value on creation for the "beta_program_enabled" field.
	DefaultBetaProgramEnabled bool
)

// OrderOption defines the ordering options for the Customer queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByUserID orders the results by the user_id field.
func ByUserID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUserID, opts...).ToFunc()
}

// ByCustomerType orders the results by the customer_type field.
func ByCustomerType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCustomerType, opts...).ToFunc()
}

// ByCustomerFirstName orders the results by the customer_first_name field.
func ByCustomerFirstName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCustomerFirstName, opts...).ToFunc()
}

// ByCustomerLastName orders the results by the customer_last_name field.
func ByCustomerLastName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCustomerLastName, opts...).ToFunc()
}

// ByCustomerMiddleName orders the results by the customer_middle_name field.
func ByCustomerMiddleName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCustomerMiddleName, opts...).ToFunc()
}

// ByCustomerTypeID orders the results by the customer_type_id field.
func ByCustomerTypeID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCustomerTypeID, opts...).ToFunc()
}

// ByCustomerSuffix orders the results by the customer_suffix field.
func ByCustomerSuffix(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCustomerSuffix, opts...).ToFunc()
}

// ByCustomerSamplesReceived orders the results by the customer_samples_received field.
func ByCustomerSamplesReceived(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCustomerSamplesReceived, opts...).ToFunc()
}

// ByCustomerRequestSubmitTime orders the results by the customer_request_submit_time field.
func ByCustomerRequestSubmitTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCustomerRequestSubmitTime, opts...).ToFunc()
}

// ByCustomerSignupTime orders the results by the customer_signup_time field.
func ByCustomerSignupTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCustomerSignupTime, opts...).ToFunc()
}

// ByIsActive orders the results by the is_active field.
func ByIsActive(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsActive, opts...).ToFunc()
}

// BySalesID orders the results by the sales_id field.
func BySalesID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSalesID, opts...).ToFunc()
}

// ByCustomerNpiNumber orders the results by the customer_npi_number field.
func ByCustomerNpiNumber(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCustomerNpiNumber, opts...).ToFunc()
}

// ByReferralSource orders the results by the referral_source field.
func ByReferralSource(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldReferralSource, opts...).ToFunc()
}

// ByOrderPlacementAllowed orders the results by the order_placement_allowed field.
func ByOrderPlacementAllowed(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOrderPlacementAllowed, opts...).ToFunc()
}

// ByBetaProgramEnabled orders the results by the beta_program_enabled field.
func ByBetaProgramEnabled(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBetaProgramEnabled, opts...).ToFunc()
}

// ByOnboardingQuestionnaireFilledOn orders the results by the onboarding_questionnaire_filled_on field.
func ByOnboardingQuestionnaireFilledOn(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOnboardingQuestionnaireFilledOn, opts...).ToFunc()
}

// BySamplesCount orders the results by samples count.
func BySamplesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newSamplesStep(), opts...)
	}
}

// BySamples orders the results by samples terms.
func BySamples(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newSamplesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByCustomerContactsCount orders the results by customer_contacts count.
func ByCustomerContactsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newCustomerContactsStep(), opts...)
	}
}

// ByCustomerContacts orders the results by customer_contacts terms.
func ByCustomerContacts(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCustomerContactsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByCustomerAddressesCount orders the results by customer_addresses count.
func ByCustomerAddressesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newCustomerAddressesStep(), opts...)
	}
}

// ByCustomerAddresses orders the results by customer_addresses terms.
func ByCustomerAddresses(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCustomerAddressesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByClinicsCount orders the results by clinics count.
func ByClinicsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newClinicsStep(), opts...)
	}
}

// ByClinics orders the results by clinics terms.
func ByClinics(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newClinicsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// BySalesField orders the results by sales field.
func BySalesField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newSalesStep(), sql.OrderByField(field, opts...))
	}
}

// ByUserField orders the results by user field.
func ByUserField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newUserStep(), sql.OrderByField(field, opts...))
	}
}

// ByOrdersCount orders the results by orders count.
func ByOrdersCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newOrdersStep(), opts...)
	}
}

// ByOrders orders the results by orders terms.
func ByOrders(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOrdersStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByCurrentPatientsCount orders the results by current_patients count.
func ByCurrentPatientsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newCurrentPatientsStep(), opts...)
	}
}

// ByCurrentPatients orders the results by current_patients terms.
func ByCurrentPatients(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCurrentPatientsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByPatientsCount orders the results by patients count.
func ByPatientsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newPatientsStep(), opts...)
	}
}

// ByPatients orders the results by patients terms.
func ByPatients(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPatientsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByCustomerBetaProgramParticipationsCount orders the results by customer_beta_program_participations count.
func ByCustomerBetaProgramParticipationsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newCustomerBetaProgramParticipationsStep(), opts...)
	}
}

// ByCustomerBetaProgramParticipations orders the results by customer_beta_program_participations terms.
func ByCustomerBetaProgramParticipations(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCustomerBetaProgramParticipationsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByCustomerSettingsOnClinicsCount orders the results by customer_settings_on_clinics count.
func ByCustomerSettingsOnClinicsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newCustomerSettingsOnClinicsStep(), opts...)
	}
}

// ByCustomerSettingsOnClinics orders the results by customer_settings_on_clinics terms.
func ByCustomerSettingsOnClinics(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCustomerSettingsOnClinicsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByCustomerAddressesOnClinicsCount orders the results by customer_addresses_on_clinics count.
func ByCustomerAddressesOnClinicsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newCustomerAddressesOnClinicsStep(), opts...)
	}
}

// ByCustomerAddressesOnClinics orders the results by customer_addresses_on_clinics terms.
func ByCustomerAddressesOnClinics(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCustomerAddressesOnClinicsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByCustomerContactsOnClinicsCount orders the results by customer_contacts_on_clinics count.
func ByCustomerContactsOnClinicsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newCustomerContactsOnClinicsStep(), opts...)
	}
}

// ByCustomerContactsOnClinics orders the results by customer_contacts_on_clinics terms.
func ByCustomerContactsOnClinics(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCustomerContactsOnClinicsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newSamplesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(SamplesInverseTable, SampleFieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, SamplesTable, SamplesColumn),
	)
}
func newCustomerContactsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CustomerContactsInverseTable, ContactFieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, CustomerContactsTable, CustomerContactsColumn),
	)
}
func newCustomerAddressesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CustomerAddressesInverseTable, AddressFieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, CustomerAddressesTable, CustomerAddressesColumn),
	)
}
func newClinicsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ClinicsInverseTable, ClinicFieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, ClinicsTable, ClinicsPrimaryKey...),
	)
}
func newSalesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(SalesInverseTable, InternalUserFieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, SalesTable, SalesColumn),
	)
}
func newUserStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(UserInverseTable, UserFieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
	)
}
func newOrdersStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OrdersInverseTable, OrderInfoFieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, OrdersTable, OrdersColumn),
	)
}
func newCurrentPatientsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CurrentPatientsInverseTable, PatientFieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, CurrentPatientsTable, CurrentPatientsColumn),
	)
}
func newPatientsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PatientsInverseTable, PatientFieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, PatientsTable, PatientsPrimaryKey...),
	)
}
func newCustomerBetaProgramParticipationsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CustomerBetaProgramParticipationsInverseTable, BetaProgramParticipationFieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, CustomerBetaProgramParticipationsTable, CustomerBetaProgramParticipationsColumn),
	)
}
func newCustomerSettingsOnClinicsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CustomerSettingsOnClinicsInverseTable, CustomerSettingOnClinicsFieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, CustomerSettingsOnClinicsTable, CustomerSettingsOnClinicsColumn),
	)
}
func newCustomerAddressesOnClinicsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CustomerAddressesOnClinicsInverseTable, CustomerAddressOnClinicsFieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, CustomerAddressesOnClinicsTable, CustomerAddressesOnClinicsColumn),
	)
}
func newCustomerContactsOnClinicsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CustomerContactsOnClinicsInverseTable, CustomerContactOnClinicsFieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, CustomerContactsOnClinicsTable, CustomerContactsOnClinicsColumn),
	)
}
