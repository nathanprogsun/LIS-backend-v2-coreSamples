// Code generated by ent, DO NOT EDIT.

package ent

import (
	"coresamples/ent/address"
	"coresamples/ent/clinic"
	"coresamples/ent/customer"
	"coresamples/ent/internaluser"
	"coresamples/ent/patient"
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// Address is the model entity for the Address schema.
type Address struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"address_id"`
	// AddressType holds the value of the "address_type" field.
	AddressType string `json:"address_type,omitempty"`
	// StreetAddress holds the value of the "street_address" field.
	StreetAddress string `json:"street_address,omitempty"`
	// AptPo holds the value of the "apt_po" field.
	AptPo string `json:"apt_po,omitempty"`
	// City holds the value of the "city" field.
	City string `json:"city,omitempty"`
	// State holds the value of the "state" field.
	State string `json:"state,omitempty"`
	// Zipcode holds the value of the "zipcode" field.
	Zipcode string `json:"zipcode,omitempty"`
	// Country holds the value of the "country" field.
	Country string `json:"country,omitempty"`
	// AddressConfirmed holds the value of the "address_confirmed" field.
	AddressConfirmed bool `json:"address_confirmed,omitempty"`
	// IsPrimaryAddress holds the value of the "is_primary_address" field.
	IsPrimaryAddress bool `json:"is_primary_address,omitempty"`
	// CustomerID holds the value of the "customer_id" field.
	CustomerID int `json:"customer_id,omitempty"`
	// PatientID holds the value of the "patient_id" field.
	PatientID int `json:"patient_id,omitempty"`
	// ClinicID holds the value of the "clinic_id" field.
	ClinicID int `json:"clinic_id,omitempty"`
	// InternalUserID holds the value of the "internal_user_id" field.
	InternalUserID int `json:"internal_user_id,omitempty"`
	// AddressLevel holds the value of the "address_level" field.
	AddressLevel int `json:"address_level,omitempty"`
	// AddressLevelName holds the value of the "address_level_name" field.
	AddressLevelName string `json:"address_level_name,omitempty"`
	// ApplyToAllGroupMember holds the value of the "apply_to_all_group_member" field.
	ApplyToAllGroupMember bool `json:"applyToAllGroupMember"`
	// GroupAddressID holds the value of the "group_address_id" field.
	GroupAddressID int `json:"group_address_id,omitempty"`
	// IsGroupAddress holds the value of the "is_group_address" field.
	IsGroupAddress bool `json:"isGroupAddress"`
	// UseAsDefaultCreateAddress holds the value of the "use_as_default_create_address" field.
	UseAsDefaultCreateAddress bool `json:"useAsDefaultCreateAddress"`
	// UseGroupAddress holds the value of the "use_group_address" field.
	UseGroupAddress bool `json:"useGroupAddress"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the AddressQuery when eager-loading is set.
	Edges        AddressEdges `json:"edges"`
	selectValues sql.SelectValues
}

// AddressEdges holds the relations/edges for other nodes in the graph.
type AddressEdges struct {
	// Clinic holds the value of the clinic edge.
	Clinic *Clinic `json:"clinic,omitempty"`
	// Customer holds the value of the customer edge.
	Customer *Customer `json:"customer,omitempty"`
	// CustomerClinicMappings holds the value of the customer_clinic_mappings edge.
	CustomerClinicMappings []*CustomerAddressOnClinics `json:"customer_clinic_mappings,omitempty"`
	// MemberAddresses holds the value of the member_addresses edge.
	MemberAddresses []*Address `json:"member_addresses,omitempty"`
	// GroupAddress holds the value of the group_address edge.
	GroupAddress *Address `json:"group_address,omitempty"`
	// InternalUser holds the value of the internal_user edge.
	InternalUser *InternalUser `json:"internal_user,omitempty"`
	// Patient holds the value of the patient edge.
	Patient *Patient `json:"patient,omitempty"`
	// Orders holds the value of the orders edge.
	Orders []*OrderInfo `json:"orders,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [8]bool
}

// ClinicOrErr returns the Clinic value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AddressEdges) ClinicOrErr() (*Clinic, error) {
	if e.loadedTypes[0] {
		if e.Clinic == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: clinic.Label}
		}
		return e.Clinic, nil
	}
	return nil, &NotLoadedError{edge: "clinic"}
}

// CustomerOrErr returns the Customer value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AddressEdges) CustomerOrErr() (*Customer, error) {
	if e.loadedTypes[1] {
		if e.Customer == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: customer.Label}
		}
		return e.Customer, nil
	}
	return nil, &NotLoadedError{edge: "customer"}
}

// CustomerClinicMappingsOrErr returns the CustomerClinicMappings value or an error if the edge
// was not loaded in eager-loading.
func (e AddressEdges) CustomerClinicMappingsOrErr() ([]*CustomerAddressOnClinics, error) {
	if e.loadedTypes[2] {
		return e.CustomerClinicMappings, nil
	}
	return nil, &NotLoadedError{edge: "customer_clinic_mappings"}
}

// MemberAddressesOrErr returns the MemberAddresses value or an error if the edge
// was not loaded in eager-loading.
func (e AddressEdges) MemberAddressesOrErr() ([]*Address, error) {
	if e.loadedTypes[3] {
		return e.MemberAddresses, nil
	}
	return nil, &NotLoadedError{edge: "member_addresses"}
}

// GroupAddressOrErr returns the GroupAddress value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AddressEdges) GroupAddressOrErr() (*Address, error) {
	if e.loadedTypes[4] {
		if e.GroupAddress == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: address.Label}
		}
		return e.GroupAddress, nil
	}
	return nil, &NotLoadedError{edge: "group_address"}
}

// InternalUserOrErr returns the InternalUser value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AddressEdges) InternalUserOrErr() (*InternalUser, error) {
	if e.loadedTypes[5] {
		if e.InternalUser == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: internaluser.Label}
		}
		return e.InternalUser, nil
	}
	return nil, &NotLoadedError{edge: "internal_user"}
}

// PatientOrErr returns the Patient value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AddressEdges) PatientOrErr() (*Patient, error) {
	if e.loadedTypes[6] {
		if e.Patient == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: patient.Label}
		}
		return e.Patient, nil
	}
	return nil, &NotLoadedError{edge: "patient"}
}

// OrdersOrErr returns the Orders value or an error if the edge
// was not loaded in eager-loading.
func (e AddressEdges) OrdersOrErr() ([]*OrderInfo, error) {
	if e.loadedTypes[7] {
		return e.Orders, nil
	}
	return nil, &NotLoadedError{edge: "orders"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Address) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case address.FieldAddressConfirmed, address.FieldIsPrimaryAddress, address.FieldApplyToAllGroupMember, address.FieldIsGroupAddress, address.FieldUseAsDefaultCreateAddress, address.FieldUseGroupAddress:
			values[i] = new(sql.NullBool)
		case address.FieldID, address.FieldCustomerID, address.FieldPatientID, address.FieldClinicID, address.FieldInternalUserID, address.FieldAddressLevel, address.FieldGroupAddressID:
			values[i] = new(sql.NullInt64)
		case address.FieldAddressType, address.FieldStreetAddress, address.FieldAptPo, address.FieldCity, address.FieldState, address.FieldZipcode, address.FieldCountry, address.FieldAddressLevelName:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Address fields.
func (a *Address) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case address.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			a.ID = int(value.Int64)
		case address.FieldAddressType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field address_type", values[i])
			} else if value.Valid {
				a.AddressType = value.String
			}
		case address.FieldStreetAddress:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field street_address", values[i])
			} else if value.Valid {
				a.StreetAddress = value.String
			}
		case address.FieldAptPo:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field apt_po", values[i])
			} else if value.Valid {
				a.AptPo = value.String
			}
		case address.FieldCity:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field city", values[i])
			} else if value.Valid {
				a.City = value.String
			}
		case address.FieldState:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field state", values[i])
			} else if value.Valid {
				a.State = value.String
			}
		case address.FieldZipcode:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field zipcode", values[i])
			} else if value.Valid {
				a.Zipcode = value.String
			}
		case address.FieldCountry:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field country", values[i])
			} else if value.Valid {
				a.Country = value.String
			}
		case address.FieldAddressConfirmed:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field address_confirmed", values[i])
			} else if value.Valid {
				a.AddressConfirmed = value.Bool
			}
		case address.FieldIsPrimaryAddress:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_primary_address", values[i])
			} else if value.Valid {
				a.IsPrimaryAddress = value.Bool
			}
		case address.FieldCustomerID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field customer_id", values[i])
			} else if value.Valid {
				a.CustomerID = int(value.Int64)
			}
		case address.FieldPatientID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field patient_id", values[i])
			} else if value.Valid {
				a.PatientID = int(value.Int64)
			}
		case address.FieldClinicID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field clinic_id", values[i])
			} else if value.Valid {
				a.ClinicID = int(value.Int64)
			}
		case address.FieldInternalUserID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field internal_user_id", values[i])
			} else if value.Valid {
				a.InternalUserID = int(value.Int64)
			}
		case address.FieldAddressLevel:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field address_level", values[i])
			} else if value.Valid {
				a.AddressLevel = int(value.Int64)
			}
		case address.FieldAddressLevelName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field address_level_name", values[i])
			} else if value.Valid {
				a.AddressLevelName = value.String
			}
		case address.FieldApplyToAllGroupMember:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field apply_to_all_group_member", values[i])
			} else if value.Valid {
				a.ApplyToAllGroupMember = value.Bool
			}
		case address.FieldGroupAddressID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field group_address_id", values[i])
			} else if value.Valid {
				a.GroupAddressID = int(value.Int64)
			}
		case address.FieldIsGroupAddress:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_group_address", values[i])
			} else if value.Valid {
				a.IsGroupAddress = value.Bool
			}
		case address.FieldUseAsDefaultCreateAddress:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field use_as_default_create_address", values[i])
			} else if value.Valid {
				a.UseAsDefaultCreateAddress = value.Bool
			}
		case address.FieldUseGroupAddress:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field use_group_address", values[i])
			} else if value.Valid {
				a.UseGroupAddress = value.Bool
			}
		default:
			a.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Address.
// This includes values selected through modifiers, order, etc.
func (a *Address) Value(name string) (ent.Value, error) {
	return a.selectValues.Get(name)
}

// QueryClinic queries the "clinic" edge of the Address entity.
func (a *Address) QueryClinic() *ClinicQuery {
	return NewAddressClient(a.config).QueryClinic(a)
}

// QueryCustomer queries the "customer" edge of the Address entity.
func (a *Address) QueryCustomer() *CustomerQuery {
	return NewAddressClient(a.config).QueryCustomer(a)
}

// QueryCustomerClinicMappings queries the "customer_clinic_mappings" edge of the Address entity.
func (a *Address) QueryCustomerClinicMappings() *CustomerAddressOnClinicsQuery {
	return NewAddressClient(a.config).QueryCustomerClinicMappings(a)
}

// QueryMemberAddresses queries the "member_addresses" edge of the Address entity.
func (a *Address) QueryMemberAddresses() *AddressQuery {
	return NewAddressClient(a.config).QueryMemberAddresses(a)
}

// QueryGroupAddress queries the "group_address" edge of the Address entity.
func (a *Address) QueryGroupAddress() *AddressQuery {
	return NewAddressClient(a.config).QueryGroupAddress(a)
}

// QueryInternalUser queries the "internal_user" edge of the Address entity.
func (a *Address) QueryInternalUser() *InternalUserQuery {
	return NewAddressClient(a.config).QueryInternalUser(a)
}

// QueryPatient queries the "patient" edge of the Address entity.
func (a *Address) QueryPatient() *PatientQuery {
	return NewAddressClient(a.config).QueryPatient(a)
}

// QueryOrders queries the "orders" edge of the Address entity.
func (a *Address) QueryOrders() *OrderInfoQuery {
	return NewAddressClient(a.config).QueryOrders(a)
}

// Update returns a builder for updating this Address.
// Note that you need to call Address.Unwrap() before calling this method if this Address
// was returned from a transaction, and the transaction was committed or rolled back.
func (a *Address) Update() *AddressUpdateOne {
	return NewAddressClient(a.config).UpdateOne(a)
}

// Unwrap unwraps the Address entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (a *Address) Unwrap() *Address {
	_tx, ok := a.config.driver.(*txDriver)
	if !ok {
		panic("ent: Address is not a transactional entity")
	}
	a.config.driver = _tx.drv
	return a
}

// String implements the fmt.Stringer.
func (a *Address) String() string {
	var builder strings.Builder
	builder.WriteString("Address(")
	builder.WriteString(fmt.Sprintf("id=%v, ", a.ID))
	builder.WriteString("address_type=")
	builder.WriteString(a.AddressType)
	builder.WriteString(", ")
	builder.WriteString("street_address=")
	builder.WriteString(a.StreetAddress)
	builder.WriteString(", ")
	builder.WriteString("apt_po=")
	builder.WriteString(a.AptPo)
	builder.WriteString(", ")
	builder.WriteString("city=")
	builder.WriteString(a.City)
	builder.WriteString(", ")
	builder.WriteString("state=")
	builder.WriteString(a.State)
	builder.WriteString(", ")
	builder.WriteString("zipcode=")
	builder.WriteString(a.Zipcode)
	builder.WriteString(", ")
	builder.WriteString("country=")
	builder.WriteString(a.Country)
	builder.WriteString(", ")
	builder.WriteString("address_confirmed=")
	builder.WriteString(fmt.Sprintf("%v", a.AddressConfirmed))
	builder.WriteString(", ")
	builder.WriteString("is_primary_address=")
	builder.WriteString(fmt.Sprintf("%v", a.IsPrimaryAddress))
	builder.WriteString(", ")
	builder.WriteString("customer_id=")
	builder.WriteString(fmt.Sprintf("%v", a.CustomerID))
	builder.WriteString(", ")
	builder.WriteString("patient_id=")
	builder.WriteString(fmt.Sprintf("%v", a.PatientID))
	builder.WriteString(", ")
	builder.WriteString("clinic_id=")
	builder.WriteString(fmt.Sprintf("%v", a.ClinicID))
	builder.WriteString(", ")
	builder.WriteString("internal_user_id=")
	builder.WriteString(fmt.Sprintf("%v", a.InternalUserID))
	builder.WriteString(", ")
	builder.WriteString("address_level=")
	builder.WriteString(fmt.Sprintf("%v", a.AddressLevel))
	builder.WriteString(", ")
	builder.WriteString("address_level_name=")
	builder.WriteString(a.AddressLevelName)
	builder.WriteString(", ")
	builder.WriteString("apply_to_all_group_member=")
	builder.WriteString(fmt.Sprintf("%v", a.ApplyToAllGroupMember))
	builder.WriteString(", ")
	builder.WriteString("group_address_id=")
	builder.WriteString(fmt.Sprintf("%v", a.GroupAddressID))
	builder.WriteString(", ")
	builder.WriteString("is_group_address=")
	builder.WriteString(fmt.Sprintf("%v", a.IsGroupAddress))
	builder.WriteString(", ")
	builder.WriteString("use_as_default_create_address=")
	builder.WriteString(fmt.Sprintf("%v", a.UseAsDefaultCreateAddress))
	builder.WriteString(", ")
	builder.WriteString("use_group_address=")
	builder.WriteString(fmt.Sprintf("%v", a.UseGroupAddress))
	builder.WriteByte(')')
	return builder.String()
}

// Addresses is a parsable slice of Address.
type Addresses []*Address
