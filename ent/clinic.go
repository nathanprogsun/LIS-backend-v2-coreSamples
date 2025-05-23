// Code generated by ent, DO NOT EDIT.

package ent

import (
	"coresamples/ent/clinic"
	"coresamples/ent/user"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// Clinic is the model entity for the Clinic schema.
type Clinic struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"clinic_id"`
	// ClinicName holds the value of the "clinic_name" field.
	ClinicName string `json:"clinic_name,omitempty"`
	// UserID holds the value of the "user_id" field.
	UserID int `json:"user_id,omitempty"`
	// IsActive holds the value of the "is_active" field.
	IsActive bool `json:"isActive"`
	// ClinicAccountID holds the value of the "clinic_account_id" field.
	ClinicAccountID int `json:"clinic_account_id,omitempty"`
	// ClinicNameOldSystem holds the value of the "clinic_name_old_system" field.
	ClinicNameOldSystem string `json:"clinic_name_old_system,omitempty"`
	// ClinicSignupTime holds the value of the "clinic_signup_time" field.
	ClinicSignupTime time.Time `json:"clinic_signup_time,omitempty"`
	// ClinicUpdatedTime holds the value of the "clinic_updated_time" field.
	ClinicUpdatedTime time.Time `json:"clinic_updated_time,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ClinicQuery when eager-loading is set.
	Edges        ClinicEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ClinicEdges holds the relations/edges for other nodes in the graph.
type ClinicEdges struct {
	// User holds the value of the user edge.
	User *User `json:"user,omitempty"`
	// ClinicContacts holds the value of the clinic_contacts edge.
	ClinicContacts []*Contact `json:"clinic_contacts,omitempty"`
	// ClinicAddresses holds the value of the clinic_addresses edge.
	ClinicAddresses []*Address `json:"clinic_addresses,omitempty"`
	// Customers holds the value of the customers edge.
	Customers []*Customer `json:"customers,omitempty"`
	// ClinicSettings holds the value of the clinic_settings edge.
	ClinicSettings []*Setting `json:"clinic_settings,omitempty"`
	// ClinicOrders holds the value of the clinic_orders edge.
	ClinicOrders []*OrderInfo `json:"clinic_orders,omitempty"`
	// ClinicPatients holds the value of the clinic_patients edge.
	ClinicPatients []*Patient `json:"clinic_patients,omitempty"`
	// ClinicBetaProgramParticipations holds the value of the clinic_beta_program_participations edge.
	ClinicBetaProgramParticipations []*BetaProgramParticipation `json:"clinic_beta_program_participations,omitempty"`
	// ClinicCustomerSettings holds the value of the clinic_customer_settings edge.
	ClinicCustomerSettings []*CustomerSettingOnClinics `json:"clinic_customer_settings,omitempty"`
	// ClinicCustomerAddresses holds the value of the clinic_customer_addresses edge.
	ClinicCustomerAddresses []*CustomerAddressOnClinics `json:"clinic_customer_addresses,omitempty"`
	// ClinicCustomerContacts holds the value of the clinic_customer_contacts edge.
	ClinicCustomerContacts []*CustomerContactOnClinics `json:"clinic_customer_contacts,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [11]bool
}

// UserOrErr returns the User value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ClinicEdges) UserOrErr() (*User, error) {
	if e.loadedTypes[0] {
		if e.User == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: user.Label}
		}
		return e.User, nil
	}
	return nil, &NotLoadedError{edge: "user"}
}

// ClinicContactsOrErr returns the ClinicContacts value or an error if the edge
// was not loaded in eager-loading.
func (e ClinicEdges) ClinicContactsOrErr() ([]*Contact, error) {
	if e.loadedTypes[1] {
		return e.ClinicContacts, nil
	}
	return nil, &NotLoadedError{edge: "clinic_contacts"}
}

// ClinicAddressesOrErr returns the ClinicAddresses value or an error if the edge
// was not loaded in eager-loading.
func (e ClinicEdges) ClinicAddressesOrErr() ([]*Address, error) {
	if e.loadedTypes[2] {
		return e.ClinicAddresses, nil
	}
	return nil, &NotLoadedError{edge: "clinic_addresses"}
}

// CustomersOrErr returns the Customers value or an error if the edge
// was not loaded in eager-loading.
func (e ClinicEdges) CustomersOrErr() ([]*Customer, error) {
	if e.loadedTypes[3] {
		return e.Customers, nil
	}
	return nil, &NotLoadedError{edge: "customers"}
}

// ClinicSettingsOrErr returns the ClinicSettings value or an error if the edge
// was not loaded in eager-loading.
func (e ClinicEdges) ClinicSettingsOrErr() ([]*Setting, error) {
	if e.loadedTypes[4] {
		return e.ClinicSettings, nil
	}
	return nil, &NotLoadedError{edge: "clinic_settings"}
}

// ClinicOrdersOrErr returns the ClinicOrders value or an error if the edge
// was not loaded in eager-loading.
func (e ClinicEdges) ClinicOrdersOrErr() ([]*OrderInfo, error) {
	if e.loadedTypes[5] {
		return e.ClinicOrders, nil
	}
	return nil, &NotLoadedError{edge: "clinic_orders"}
}

// ClinicPatientsOrErr returns the ClinicPatients value or an error if the edge
// was not loaded in eager-loading.
func (e ClinicEdges) ClinicPatientsOrErr() ([]*Patient, error) {
	if e.loadedTypes[6] {
		return e.ClinicPatients, nil
	}
	return nil, &NotLoadedError{edge: "clinic_patients"}
}

// ClinicBetaProgramParticipationsOrErr returns the ClinicBetaProgramParticipations value or an error if the edge
// was not loaded in eager-loading.
func (e ClinicEdges) ClinicBetaProgramParticipationsOrErr() ([]*BetaProgramParticipation, error) {
	if e.loadedTypes[7] {
		return e.ClinicBetaProgramParticipations, nil
	}
	return nil, &NotLoadedError{edge: "clinic_beta_program_participations"}
}

// ClinicCustomerSettingsOrErr returns the ClinicCustomerSettings value or an error if the edge
// was not loaded in eager-loading.
func (e ClinicEdges) ClinicCustomerSettingsOrErr() ([]*CustomerSettingOnClinics, error) {
	if e.loadedTypes[8] {
		return e.ClinicCustomerSettings, nil
	}
	return nil, &NotLoadedError{edge: "clinic_customer_settings"}
}

// ClinicCustomerAddressesOrErr returns the ClinicCustomerAddresses value or an error if the edge
// was not loaded in eager-loading.
func (e ClinicEdges) ClinicCustomerAddressesOrErr() ([]*CustomerAddressOnClinics, error) {
	if e.loadedTypes[9] {
		return e.ClinicCustomerAddresses, nil
	}
	return nil, &NotLoadedError{edge: "clinic_customer_addresses"}
}

// ClinicCustomerContactsOrErr returns the ClinicCustomerContacts value or an error if the edge
// was not loaded in eager-loading.
func (e ClinicEdges) ClinicCustomerContactsOrErr() ([]*CustomerContactOnClinics, error) {
	if e.loadedTypes[10] {
		return e.ClinicCustomerContacts, nil
	}
	return nil, &NotLoadedError{edge: "clinic_customer_contacts"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Clinic) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case clinic.FieldIsActive:
			values[i] = new(sql.NullBool)
		case clinic.FieldID, clinic.FieldUserID, clinic.FieldClinicAccountID:
			values[i] = new(sql.NullInt64)
		case clinic.FieldClinicName, clinic.FieldClinicNameOldSystem:
			values[i] = new(sql.NullString)
		case clinic.FieldClinicSignupTime, clinic.FieldClinicUpdatedTime:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Clinic fields.
func (c *Clinic) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case clinic.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			c.ID = int(value.Int64)
		case clinic.FieldClinicName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field clinic_name", values[i])
			} else if value.Valid {
				c.ClinicName = value.String
			}
		case clinic.FieldUserID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field user_id", values[i])
			} else if value.Valid {
				c.UserID = int(value.Int64)
			}
		case clinic.FieldIsActive:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_active", values[i])
			} else if value.Valid {
				c.IsActive = value.Bool
			}
		case clinic.FieldClinicAccountID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field clinic_account_id", values[i])
			} else if value.Valid {
				c.ClinicAccountID = int(value.Int64)
			}
		case clinic.FieldClinicNameOldSystem:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field clinic_name_old_system", values[i])
			} else if value.Valid {
				c.ClinicNameOldSystem = value.String
			}
		case clinic.FieldClinicSignupTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field clinic_signup_time", values[i])
			} else if value.Valid {
				c.ClinicSignupTime = value.Time
			}
		case clinic.FieldClinicUpdatedTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field clinic_updated_time", values[i])
			} else if value.Valid {
				c.ClinicUpdatedTime = value.Time
			}
		default:
			c.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Clinic.
// This includes values selected through modifiers, order, etc.
func (c *Clinic) Value(name string) (ent.Value, error) {
	return c.selectValues.Get(name)
}

// QueryUser queries the "user" edge of the Clinic entity.
func (c *Clinic) QueryUser() *UserQuery {
	return NewClinicClient(c.config).QueryUser(c)
}

// QueryClinicContacts queries the "clinic_contacts" edge of the Clinic entity.
func (c *Clinic) QueryClinicContacts() *ContactQuery {
	return NewClinicClient(c.config).QueryClinicContacts(c)
}

// QueryClinicAddresses queries the "clinic_addresses" edge of the Clinic entity.
func (c *Clinic) QueryClinicAddresses() *AddressQuery {
	return NewClinicClient(c.config).QueryClinicAddresses(c)
}

// QueryCustomers queries the "customers" edge of the Clinic entity.
func (c *Clinic) QueryCustomers() *CustomerQuery {
	return NewClinicClient(c.config).QueryCustomers(c)
}

// QueryClinicSettings queries the "clinic_settings" edge of the Clinic entity.
func (c *Clinic) QueryClinicSettings() *SettingQuery {
	return NewClinicClient(c.config).QueryClinicSettings(c)
}

// QueryClinicOrders queries the "clinic_orders" edge of the Clinic entity.
func (c *Clinic) QueryClinicOrders() *OrderInfoQuery {
	return NewClinicClient(c.config).QueryClinicOrders(c)
}

// QueryClinicPatients queries the "clinic_patients" edge of the Clinic entity.
func (c *Clinic) QueryClinicPatients() *PatientQuery {
	return NewClinicClient(c.config).QueryClinicPatients(c)
}

// QueryClinicBetaProgramParticipations queries the "clinic_beta_program_participations" edge of the Clinic entity.
func (c *Clinic) QueryClinicBetaProgramParticipations() *BetaProgramParticipationQuery {
	return NewClinicClient(c.config).QueryClinicBetaProgramParticipations(c)
}

// QueryClinicCustomerSettings queries the "clinic_customer_settings" edge of the Clinic entity.
func (c *Clinic) QueryClinicCustomerSettings() *CustomerSettingOnClinicsQuery {
	return NewClinicClient(c.config).QueryClinicCustomerSettings(c)
}

// QueryClinicCustomerAddresses queries the "clinic_customer_addresses" edge of the Clinic entity.
func (c *Clinic) QueryClinicCustomerAddresses() *CustomerAddressOnClinicsQuery {
	return NewClinicClient(c.config).QueryClinicCustomerAddresses(c)
}

// QueryClinicCustomerContacts queries the "clinic_customer_contacts" edge of the Clinic entity.
func (c *Clinic) QueryClinicCustomerContacts() *CustomerContactOnClinicsQuery {
	return NewClinicClient(c.config).QueryClinicCustomerContacts(c)
}

// Update returns a builder for updating this Clinic.
// Note that you need to call Clinic.Unwrap() before calling this method if this Clinic
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Clinic) Update() *ClinicUpdateOne {
	return NewClinicClient(c.config).UpdateOne(c)
}

// Unwrap unwraps the Clinic entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Clinic) Unwrap() *Clinic {
	_tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Clinic is not a transactional entity")
	}
	c.config.driver = _tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Clinic) String() string {
	var builder strings.Builder
	builder.WriteString("Clinic(")
	builder.WriteString(fmt.Sprintf("id=%v, ", c.ID))
	builder.WriteString("clinic_name=")
	builder.WriteString(c.ClinicName)
	builder.WriteString(", ")
	builder.WriteString("user_id=")
	builder.WriteString(fmt.Sprintf("%v", c.UserID))
	builder.WriteString(", ")
	builder.WriteString("is_active=")
	builder.WriteString(fmt.Sprintf("%v", c.IsActive))
	builder.WriteString(", ")
	builder.WriteString("clinic_account_id=")
	builder.WriteString(fmt.Sprintf("%v", c.ClinicAccountID))
	builder.WriteString(", ")
	builder.WriteString("clinic_name_old_system=")
	builder.WriteString(c.ClinicNameOldSystem)
	builder.WriteString(", ")
	builder.WriteString("clinic_signup_time=")
	builder.WriteString(c.ClinicSignupTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("clinic_updated_time=")
	builder.WriteString(c.ClinicUpdatedTime.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// Clinics is a parsable slice of Clinic.
type Clinics []*Clinic
