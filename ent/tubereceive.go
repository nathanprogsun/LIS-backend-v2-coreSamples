// Code generated by ent, DO NOT EDIT.

package ent

import (
	"coresamples/ent/sample"
	"coresamples/ent/tubereceive"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// TubeReceive is the model entity for the TubeReceive schema.
type TubeReceive struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// SampleID holds the value of the "sample_id" field.
	SampleID int `json:"sample_id,omitempty"`
	// TubeType holds the value of the "tube_type" field.
	TubeType string `json:"tube_type,omitempty"`
	// ReceivedCount holds the value of the "received_count" field.
	ReceivedCount int `json:"received_count,omitempty"`
	// ReceivedTime holds the value of the "received_time" field.
	ReceivedTime time.Time `json:"received_time,omitempty"`
	// ReceivedBy holds the value of the "received_by" field.
	ReceivedBy string `json:"received_by,omitempty"`
	// ModifiedBy holds the value of the "modified_by" field.
	ModifiedBy string `json:"modified_by,omitempty"`
	// ModifiedTime holds the value of the "modified_time" field.
	ModifiedTime time.Time `json:"modified_time,omitempty"`
	// CollectionTime holds the value of the "collection_time" field.
	CollectionTime time.Time `json:"collection_time,omitempty"`
	// IsRedraw holds the value of the "is_redraw" field.
	IsRedraw bool `json:"is_redraw,omitempty"`
	// IsRerun holds the value of the "is_rerun" field.
	IsRerun bool `json:"is_rerun,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the TubeReceiveQuery when eager-loading is set.
	Edges        TubeReceiveEdges `json:"edges"`
	selectValues sql.SelectValues
}

// TubeReceiveEdges holds the relations/edges for other nodes in the graph.
type TubeReceiveEdges struct {
	// Sample holds the value of the sample edge.
	Sample *Sample `json:"sample,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// SampleOrErr returns the Sample value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e TubeReceiveEdges) SampleOrErr() (*Sample, error) {
	if e.loadedTypes[0] {
		if e.Sample == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: sample.Label}
		}
		return e.Sample, nil
	}
	return nil, &NotLoadedError{edge: "sample"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*TubeReceive) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case tubereceive.FieldIsRedraw, tubereceive.FieldIsRerun:
			values[i] = new(sql.NullBool)
		case tubereceive.FieldID, tubereceive.FieldSampleID, tubereceive.FieldReceivedCount:
			values[i] = new(sql.NullInt64)
		case tubereceive.FieldTubeType, tubereceive.FieldReceivedBy, tubereceive.FieldModifiedBy:
			values[i] = new(sql.NullString)
		case tubereceive.FieldReceivedTime, tubereceive.FieldModifiedTime, tubereceive.FieldCollectionTime:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the TubeReceive fields.
func (tr *TubeReceive) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case tubereceive.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			tr.ID = int(value.Int64)
		case tubereceive.FieldSampleID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field sample_id", values[i])
			} else if value.Valid {
				tr.SampleID = int(value.Int64)
			}
		case tubereceive.FieldTubeType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field tube_type", values[i])
			} else if value.Valid {
				tr.TubeType = value.String
			}
		case tubereceive.FieldReceivedCount:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field received_count", values[i])
			} else if value.Valid {
				tr.ReceivedCount = int(value.Int64)
			}
		case tubereceive.FieldReceivedTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field received_time", values[i])
			} else if value.Valid {
				tr.ReceivedTime = value.Time
			}
		case tubereceive.FieldReceivedBy:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field received_by", values[i])
			} else if value.Valid {
				tr.ReceivedBy = value.String
			}
		case tubereceive.FieldModifiedBy:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field modified_by", values[i])
			} else if value.Valid {
				tr.ModifiedBy = value.String
			}
		case tubereceive.FieldModifiedTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field modified_time", values[i])
			} else if value.Valid {
				tr.ModifiedTime = value.Time
			}
		case tubereceive.FieldCollectionTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field collection_time", values[i])
			} else if value.Valid {
				tr.CollectionTime = value.Time
			}
		case tubereceive.FieldIsRedraw:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_redraw", values[i])
			} else if value.Valid {
				tr.IsRedraw = value.Bool
			}
		case tubereceive.FieldIsRerun:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_rerun", values[i])
			} else if value.Valid {
				tr.IsRerun = value.Bool
			}
		default:
			tr.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the TubeReceive.
// This includes values selected through modifiers, order, etc.
func (tr *TubeReceive) Value(name string) (ent.Value, error) {
	return tr.selectValues.Get(name)
}

// QuerySample queries the "sample" edge of the TubeReceive entity.
func (tr *TubeReceive) QuerySample() *SampleQuery {
	return NewTubeReceiveClient(tr.config).QuerySample(tr)
}

// Update returns a builder for updating this TubeReceive.
// Note that you need to call TubeReceive.Unwrap() before calling this method if this TubeReceive
// was returned from a transaction, and the transaction was committed or rolled back.
func (tr *TubeReceive) Update() *TubeReceiveUpdateOne {
	return NewTubeReceiveClient(tr.config).UpdateOne(tr)
}

// Unwrap unwraps the TubeReceive entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (tr *TubeReceive) Unwrap() *TubeReceive {
	_tx, ok := tr.config.driver.(*txDriver)
	if !ok {
		panic("ent: TubeReceive is not a transactional entity")
	}
	tr.config.driver = _tx.drv
	return tr
}

// String implements the fmt.Stringer.
func (tr *TubeReceive) String() string {
	var builder strings.Builder
	builder.WriteString("TubeReceive(")
	builder.WriteString(fmt.Sprintf("id=%v, ", tr.ID))
	builder.WriteString("sample_id=")
	builder.WriteString(fmt.Sprintf("%v", tr.SampleID))
	builder.WriteString(", ")
	builder.WriteString("tube_type=")
	builder.WriteString(tr.TubeType)
	builder.WriteString(", ")
	builder.WriteString("received_count=")
	builder.WriteString(fmt.Sprintf("%v", tr.ReceivedCount))
	builder.WriteString(", ")
	builder.WriteString("received_time=")
	builder.WriteString(tr.ReceivedTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("received_by=")
	builder.WriteString(tr.ReceivedBy)
	builder.WriteString(", ")
	builder.WriteString("modified_by=")
	builder.WriteString(tr.ModifiedBy)
	builder.WriteString(", ")
	builder.WriteString("modified_time=")
	builder.WriteString(tr.ModifiedTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("collection_time=")
	builder.WriteString(tr.CollectionTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("is_redraw=")
	builder.WriteString(fmt.Sprintf("%v", tr.IsRedraw))
	builder.WriteString(", ")
	builder.WriteString("is_rerun=")
	builder.WriteString(fmt.Sprintf("%v", tr.IsRerun))
	builder.WriteByte(')')
	return builder.String()
}

// TubeReceives is a parsable slice of TubeReceive.
type TubeReceives []*TubeReceive
