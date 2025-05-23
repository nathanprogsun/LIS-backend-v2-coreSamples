// Code generated by ent, DO NOT EDIT.

package ent

import (
	"coresamples/ent/address"
	"coresamples/ent/clinic"
	"coresamples/ent/customer"
	"coresamples/ent/customeraddressonclinics"
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// CustomerAddressOnClinics is the model entity for the CustomerAddressOnClinics schema.
type CustomerAddressOnClinics struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CustomerID holds the value of the "customer_id" field.
	CustomerID int `json:"customer_id,omitempty"`
	// ClinicID holds the value of the "clinic_id" field.
	ClinicID int `json:"clinic_id,omitempty"`
	// AddressID holds the value of the "address_id" field.
	AddressID int `json:"address_id,omitempty"`
	// AddressType holds the value of the "address_type" field.
	AddressType string `json:"address_type,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the CustomerAddressOnClinicsQuery when eager-loading is set.
	Edges        CustomerAddressOnClinicsEdges `json:"edges"`
	selectValues sql.SelectValues
}

// CustomerAddressOnClinicsEdges holds the relations/edges for other nodes in the graph.
type CustomerAddressOnClinicsEdges struct {
	// Customer holds the value of the customer edge.
	Customer *Customer `json:"customer,omitempty"`
	// Clinic holds the value of the clinic edge.
	Clinic *Clinic `json:"clinic,omitempty"`
	// Address holds the value of the address edge.
	Address *Address `json:"address,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// CustomerOrErr returns the Customer value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e CustomerAddressOnClinicsEdges) CustomerOrErr() (*Customer, error) {
	if e.loadedTypes[0] {
		if e.Customer == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: customer.Label}
		}
		return e.Customer, nil
	}
	return nil, &NotLoadedError{edge: "customer"}
}

// ClinicOrErr returns the Clinic value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e CustomerAddressOnClinicsEdges) ClinicOrErr() (*Clinic, error) {
	if e.loadedTypes[1] {
		if e.Clinic == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: clinic.Label}
		}
		return e.Clinic, nil
	}
	return nil, &NotLoadedError{edge: "clinic"}
}

// AddressOrErr returns the Address value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e CustomerAddressOnClinicsEdges) AddressOrErr() (*Address, error) {
	if e.loadedTypes[2] {
		if e.Address == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: address.Label}
		}
		return e.Address, nil
	}
	return nil, &NotLoadedError{edge: "address"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*CustomerAddressOnClinics) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case customeraddressonclinics.FieldID, customeraddressonclinics.FieldCustomerID, customeraddressonclinics.FieldClinicID, customeraddressonclinics.FieldAddressID:
			values[i] = new(sql.NullInt64)
		case customeraddressonclinics.FieldAddressType:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the CustomerAddressOnClinics fields.
func (caoc *CustomerAddressOnClinics) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case customeraddressonclinics.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			caoc.ID = int(value.Int64)
		case customeraddressonclinics.FieldCustomerID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field customer_id", values[i])
			} else if value.Valid {
				caoc.CustomerID = int(value.Int64)
			}
		case customeraddressonclinics.FieldClinicID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field clinic_id", values[i])
			} else if value.Valid {
				caoc.ClinicID = int(value.Int64)
			}
		case customeraddressonclinics.FieldAddressID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field address_id", values[i])
			} else if value.Valid {
				caoc.AddressID = int(value.Int64)
			}
		case customeraddressonclinics.FieldAddressType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field address_type", values[i])
			} else if value.Valid {
				caoc.AddressType = value.String
			}
		default:
			caoc.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the CustomerAddressOnClinics.
// This includes values selected through modifiers, order, etc.
func (caoc *CustomerAddressOnClinics) Value(name string) (ent.Value, error) {
	return caoc.selectValues.Get(name)
}

// QueryCustomer queries the "customer" edge of the CustomerAddressOnClinics entity.
func (caoc *CustomerAddressOnClinics) QueryCustomer() *CustomerQuery {
	return NewCustomerAddressOnClinicsClient(caoc.config).QueryCustomer(caoc)
}

// QueryClinic queries the "clinic" edge of the CustomerAddressOnClinics entity.
func (caoc *CustomerAddressOnClinics) QueryClinic() *ClinicQuery {
	return NewCustomerAddressOnClinicsClient(caoc.config).QueryClinic(caoc)
}

// QueryAddress queries the "address" edge of the CustomerAddressOnClinics entity.
func (caoc *CustomerAddressOnClinics) QueryAddress() *AddressQuery {
	return NewCustomerAddressOnClinicsClient(caoc.config).QueryAddress(caoc)
}

// Update returns a builder for updating this CustomerAddressOnClinics.
// Note that you need to call CustomerAddressOnClinics.Unwrap() before calling this method if this CustomerAddressOnClinics
// was returned from a transaction, and the transaction was committed or rolled back.
func (caoc *CustomerAddressOnClinics) Update() *CustomerAddressOnClinicsUpdateOne {
	return NewCustomerAddressOnClinicsClient(caoc.config).UpdateOne(caoc)
}

// Unwrap unwraps the CustomerAddressOnClinics entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (caoc *CustomerAddressOnClinics) Unwrap() *CustomerAddressOnClinics {
	_tx, ok := caoc.config.driver.(*txDriver)
	if !ok {
		panic("ent: CustomerAddressOnClinics is not a transactional entity")
	}
	caoc.config.driver = _tx.drv
	return caoc
}

// String implements the fmt.Stringer.
func (caoc *CustomerAddressOnClinics) String() string {
	var builder strings.Builder
	builder.WriteString("CustomerAddressOnClinics(")
	builder.WriteString(fmt.Sprintf("id=%v, ", caoc.ID))
	builder.WriteString("customer_id=")
	builder.WriteString(fmt.Sprintf("%v", caoc.CustomerID))
	builder.WriteString(", ")
	builder.WriteString("clinic_id=")
	builder.WriteString(fmt.Sprintf("%v", caoc.ClinicID))
	builder.WriteString(", ")
	builder.WriteString("address_id=")
	builder.WriteString(fmt.Sprintf("%v", caoc.AddressID))
	builder.WriteString(", ")
	builder.WriteString("address_type=")
	builder.WriteString(caoc.AddressType)
	builder.WriteByte(')')
	return builder.String()
}

// CustomerAddressOnClinicsSlice is a parsable slice of CustomerAddressOnClinics.
type CustomerAddressOnClinicsSlice []*CustomerAddressOnClinics
