// Code generated by ent, DO NOT EDIT.

package ent

import (
	"coresamples/ent/tubetype"
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// TubeType is the model entity for the TubeType schema.
type TubeType struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// TubeName holds the value of the "tube_name" field.
	TubeName string `json:"tube_name,omitempty"`
	// TubeTypeEnum holds the value of the "tube_type_enum" field.
	TubeTypeEnum string `json:"tube_type_enum,omitempty"`
	// TubeTypeSymbol holds the value of the "tube_type_symbol" field.
	TubeTypeSymbol string `json:"tube_type_symbol,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the TubeTypeQuery when eager-loading is set.
	Edges        TubeTypeEdges `json:"edges"`
	selectValues sql.SelectValues
}

// TubeTypeEdges holds the relations/edges for other nodes in the graph.
type TubeTypeEdges struct {
	// Tube holds the value of the tube edge.
	Tube []*Tube `json:"tube,omitempty"`
	// SampleTypes holds the value of the sample_types edge.
	SampleTypes []*SampleType `json:"sample_types,omitempty"`
	// Tests holds the value of the tests edge.
	Tests []*Test `json:"tests,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// TubeOrErr returns the Tube value or an error if the edge
// was not loaded in eager-loading.
func (e TubeTypeEdges) TubeOrErr() ([]*Tube, error) {
	if e.loadedTypes[0] {
		return e.Tube, nil
	}
	return nil, &NotLoadedError{edge: "tube"}
}

// SampleTypesOrErr returns the SampleTypes value or an error if the edge
// was not loaded in eager-loading.
func (e TubeTypeEdges) SampleTypesOrErr() ([]*SampleType, error) {
	if e.loadedTypes[1] {
		return e.SampleTypes, nil
	}
	return nil, &NotLoadedError{edge: "sample_types"}
}

// TestsOrErr returns the Tests value or an error if the edge
// was not loaded in eager-loading.
func (e TubeTypeEdges) TestsOrErr() ([]*Test, error) {
	if e.loadedTypes[2] {
		return e.Tests, nil
	}
	return nil, &NotLoadedError{edge: "tests"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*TubeType) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case tubetype.FieldID:
			values[i] = new(sql.NullInt64)
		case tubetype.FieldTubeName, tubetype.FieldTubeTypeEnum, tubetype.FieldTubeTypeSymbol:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the TubeType fields.
func (tt *TubeType) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case tubetype.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			tt.ID = int(value.Int64)
		case tubetype.FieldTubeName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field tube_name", values[i])
			} else if value.Valid {
				tt.TubeName = value.String
			}
		case tubetype.FieldTubeTypeEnum:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field tube_type_enum", values[i])
			} else if value.Valid {
				tt.TubeTypeEnum = value.String
			}
		case tubetype.FieldTubeTypeSymbol:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field tube_type_symbol", values[i])
			} else if value.Valid {
				tt.TubeTypeSymbol = value.String
			}
		default:
			tt.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the TubeType.
// This includes values selected through modifiers, order, etc.
func (tt *TubeType) Value(name string) (ent.Value, error) {
	return tt.selectValues.Get(name)
}

// QueryTube queries the "tube" edge of the TubeType entity.
func (tt *TubeType) QueryTube() *TubeQuery {
	return NewTubeTypeClient(tt.config).QueryTube(tt)
}

// QuerySampleTypes queries the "sample_types" edge of the TubeType entity.
func (tt *TubeType) QuerySampleTypes() *SampleTypeQuery {
	return NewTubeTypeClient(tt.config).QuerySampleTypes(tt)
}

// QueryTests queries the "tests" edge of the TubeType entity.
func (tt *TubeType) QueryTests() *TestQuery {
	return NewTubeTypeClient(tt.config).QueryTests(tt)
}

// Update returns a builder for updating this TubeType.
// Note that you need to call TubeType.Unwrap() before calling this method if this TubeType
// was returned from a transaction, and the transaction was committed or rolled back.
func (tt *TubeType) Update() *TubeTypeUpdateOne {
	return NewTubeTypeClient(tt.config).UpdateOne(tt)
}

// Unwrap unwraps the TubeType entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (tt *TubeType) Unwrap() *TubeType {
	_tx, ok := tt.config.driver.(*txDriver)
	if !ok {
		panic("ent: TubeType is not a transactional entity")
	}
	tt.config.driver = _tx.drv
	return tt
}

// String implements the fmt.Stringer.
func (tt *TubeType) String() string {
	var builder strings.Builder
	builder.WriteString("TubeType(")
	builder.WriteString(fmt.Sprintf("id=%v, ", tt.ID))
	builder.WriteString("tube_name=")
	builder.WriteString(tt.TubeName)
	builder.WriteString(", ")
	builder.WriteString("tube_type_enum=")
	builder.WriteString(tt.TubeTypeEnum)
	builder.WriteString(", ")
	builder.WriteString("tube_type_symbol=")
	builder.WriteString(tt.TubeTypeSymbol)
	builder.WriteByte(')')
	return builder.String()
}

// TubeTypes is a parsable slice of TubeType.
type TubeTypes []*TubeType
