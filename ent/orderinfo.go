// Code generated by ent, DO NOT EDIT.

package ent

import (
	"coresamples/ent/address"
	"coresamples/ent/clinic"
	"coresamples/ent/contact"
	"coresamples/ent/customer"
	"coresamples/ent/orderinfo"
	"coresamples/ent/sample"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// OrderInfo is the model entity for the OrderInfo schema.
type OrderInfo struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"order_id"`
	// OrderTitle holds the value of the "order_title" field.
	OrderTitle string `json:"order_title,omitempty"`
	// OrderType holds the value of the "order_type" field.
	OrderType string `json:"order_type,omitempty"`
	// OrderDescription holds the value of the "order_description" field.
	OrderDescription string `json:"order_description,omitempty"`
	// OrderConfirmationNumber holds the value of the "order_confirmation_number" field.
	OrderConfirmationNumber string `json:"order_confirmation_number,omitempty"`
	// ClinicID holds the value of the "clinic_id" field.
	ClinicID int `json:"clinic_id,omitempty"`
	// CustomerID holds the value of the "customer_id" field.
	CustomerID int `json:"customer_id,omitempty"`
	// OrderCreateTime holds the value of the "order_create_time" field.
	OrderCreateTime time.Time `json:"order_create_time,omitempty"`
	// OrderServiceTime holds the value of the "order_service_time" field.
	OrderServiceTime time.Time `json:"order_service_time,omitempty"`
	// OrderProcessTime holds the value of the "order_process_time" field.
	OrderProcessTime time.Time `json:"order_process_time,omitempty"`
	// OrderRedrawTime holds the value of the "order_redraw_time" field.
	OrderRedrawTime time.Time `json:"order_redraw_time,omitempty"`
	// OrderCancelTime holds the value of the "order_cancel_time" field.
	OrderCancelTime time.Time `json:"order_cancel_time,omitempty"`
	// IsActive holds the value of the "isActive" field.
	IsActive bool `json:"isActive,omitempty"`
	// HasOrderSetting holds the value of the "has_order_setting" field.
	HasOrderSetting bool `json:"has_order_setting,omitempty"`
	// OrderCanceled holds the value of the "order_canceled" field.
	OrderCanceled bool `json:"order_canceled,omitempty"`
	// OrderFlagged holds the value of the "order_flagged" field.
	OrderFlagged bool `json:"order_flagged,omitempty"`
	// OrderStatus holds the value of the "order_status" field.
	OrderStatus string `json:"order_status,omitempty"`
	// OrderMajorStatus holds the value of the "order_major_status" field.
	OrderMajorStatus string `json:"order_major_status,omitempty"`
	// OrderKitStatus holds the value of the "order_kit_status" field.
	OrderKitStatus string `json:"order_kit_status,omitempty"`
	// OrderReportStatus holds the value of the "order_report_status" field.
	OrderReportStatus string `json:"order_report_status,omitempty"`
	// OrderTnpIssueStatus holds the value of the "order_tnp_issue_status" field.
	OrderTnpIssueStatus string `json:"order_tnp_issue_status,omitempty"`
	// OrderBillingIssueStatus holds the value of the "order_billing_issue_status" field.
	OrderBillingIssueStatus string `json:"order_billing_issue_status,omitempty"`
	// OrderMissingInfoIssueStatus holds the value of the "order_missing_info_issue_status" field.
	OrderMissingInfoIssueStatus string `json:"order_missing_info_issue_status,omitempty"`
	// OrderIncompleteQuestionnaireIssueStatus holds the value of the "order_incomplete_questionnaire_issue_status" field.
	OrderIncompleteQuestionnaireIssueStatus string `json:"order_incomplete_questionnaire_issue_status,omitempty"`
	// OrderNyWaiveFormIssueStatus holds the value of the "order_ny_waive_form_issue_status" field.
	OrderNyWaiveFormIssueStatus string `json:"order_ny_waive_form_issue_status,omitempty"`
	// OrderLabIssueStatus holds the value of the "order_lab_issue_status" field.
	OrderLabIssueStatus string `json:"order_lab_issue_status,omitempty"`
	// OrderProcessingTime holds the value of the "order_processing_time" field.
	OrderProcessingTime time.Time `json:"order_processing_time,omitempty"`
	// OrderMinorStatus holds the value of the "order_minor_status" field.
	OrderMinorStatus string `json:"order_minor_status,omitempty"`
	// PatientFirstName holds the value of the "patient_first_name" field.
	PatientFirstName string `json:"patient_first_name,omitempty"`
	// PatientLastName holds the value of the "patient_last_name" field.
	PatientLastName string `json:"patient_last_name,omitempty"`
	// OrderSource holds the value of the "order_source" field.
	OrderSource string `json:"order_source,omitempty"`
	// OrderChargeMethod holds the value of the "order_charge_method" field.
	OrderChargeMethod string `json:"order_charge_method,omitempty"`
	// OrderPlacingType holds the value of the "order_placing_type" field.
	OrderPlacingType string `json:"order_placing_type,omitempty"`
	// BillingOrderID holds the value of the "billing_order_id" field.
	BillingOrderID string `json:"billing_order_id,omitempty"`
	// ContactID holds the value of the "contact_id" field.
	ContactID int `json:"contact_id,omitempty"`
	// AddressID holds the value of the "address_id" field.
	AddressID int `json:"address_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the OrderInfoQuery when eager-loading is set.
	Edges        OrderInfoEdges `json:"edges"`
	selectValues sql.SelectValues
}

// OrderInfoEdges holds the relations/edges for other nodes in the graph.
type OrderInfoEdges struct {
	// Tests holds the value of the tests edge.
	Tests []*Test `json:"tests,omitempty"`
	// OrderFlags holds the value of the order_flags edge.
	OrderFlags []*OrderFlag `json:"order_flags,omitempty"`
	// Sample holds the value of the sample edge.
	Sample *Sample `json:"sample,omitempty"`
	// Contact holds the value of the contact edge.
	Contact *Contact `json:"contact,omitempty"`
	// Address holds the value of the address edge.
	Address *Address `json:"address,omitempty"`
	// Clinic holds the value of the clinic edge.
	Clinic *Clinic `json:"clinic,omitempty"`
	// CustomerInfo holds the value of the customer_info edge.
	CustomerInfo *Customer `json:"customer_info,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [7]bool
}

// TestsOrErr returns the Tests value or an error if the edge
// was not loaded in eager-loading.
func (e OrderInfoEdges) TestsOrErr() ([]*Test, error) {
	if e.loadedTypes[0] {
		return e.Tests, nil
	}
	return nil, &NotLoadedError{edge: "tests"}
}

// OrderFlagsOrErr returns the OrderFlags value or an error if the edge
// was not loaded in eager-loading.
func (e OrderInfoEdges) OrderFlagsOrErr() ([]*OrderFlag, error) {
	if e.loadedTypes[1] {
		return e.OrderFlags, nil
	}
	return nil, &NotLoadedError{edge: "order_flags"}
}

// SampleOrErr returns the Sample value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e OrderInfoEdges) SampleOrErr() (*Sample, error) {
	if e.loadedTypes[2] {
		if e.Sample == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: sample.Label}
		}
		return e.Sample, nil
	}
	return nil, &NotLoadedError{edge: "sample"}
}

// ContactOrErr returns the Contact value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e OrderInfoEdges) ContactOrErr() (*Contact, error) {
	if e.loadedTypes[3] {
		if e.Contact == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: contact.Label}
		}
		return e.Contact, nil
	}
	return nil, &NotLoadedError{edge: "contact"}
}

// AddressOrErr returns the Address value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e OrderInfoEdges) AddressOrErr() (*Address, error) {
	if e.loadedTypes[4] {
		if e.Address == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: address.Label}
		}
		return e.Address, nil
	}
	return nil, &NotLoadedError{edge: "address"}
}

// ClinicOrErr returns the Clinic value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e OrderInfoEdges) ClinicOrErr() (*Clinic, error) {
	if e.loadedTypes[5] {
		if e.Clinic == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: clinic.Label}
		}
		return e.Clinic, nil
	}
	return nil, &NotLoadedError{edge: "clinic"}
}

// CustomerInfoOrErr returns the CustomerInfo value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e OrderInfoEdges) CustomerInfoOrErr() (*Customer, error) {
	if e.loadedTypes[6] {
		if e.CustomerInfo == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: customer.Label}
		}
		return e.CustomerInfo, nil
	}
	return nil, &NotLoadedError{edge: "customer_info"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*OrderInfo) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case orderinfo.FieldIsActive, orderinfo.FieldHasOrderSetting, orderinfo.FieldOrderCanceled, orderinfo.FieldOrderFlagged:
			values[i] = new(sql.NullBool)
		case orderinfo.FieldID, orderinfo.FieldClinicID, orderinfo.FieldCustomerID, orderinfo.FieldContactID, orderinfo.FieldAddressID:
			values[i] = new(sql.NullInt64)
		case orderinfo.FieldOrderTitle, orderinfo.FieldOrderType, orderinfo.FieldOrderDescription, orderinfo.FieldOrderConfirmationNumber, orderinfo.FieldOrderStatus, orderinfo.FieldOrderMajorStatus, orderinfo.FieldOrderKitStatus, orderinfo.FieldOrderReportStatus, orderinfo.FieldOrderTnpIssueStatus, orderinfo.FieldOrderBillingIssueStatus, orderinfo.FieldOrderMissingInfoIssueStatus, orderinfo.FieldOrderIncompleteQuestionnaireIssueStatus, orderinfo.FieldOrderNyWaiveFormIssueStatus, orderinfo.FieldOrderLabIssueStatus, orderinfo.FieldOrderMinorStatus, orderinfo.FieldPatientFirstName, orderinfo.FieldPatientLastName, orderinfo.FieldOrderSource, orderinfo.FieldOrderChargeMethod, orderinfo.FieldOrderPlacingType, orderinfo.FieldBillingOrderID:
			values[i] = new(sql.NullString)
		case orderinfo.FieldOrderCreateTime, orderinfo.FieldOrderServiceTime, orderinfo.FieldOrderProcessTime, orderinfo.FieldOrderRedrawTime, orderinfo.FieldOrderCancelTime, orderinfo.FieldOrderProcessingTime:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the OrderInfo fields.
func (oi *OrderInfo) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case orderinfo.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			oi.ID = int(value.Int64)
		case orderinfo.FieldOrderTitle:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_title", values[i])
			} else if value.Valid {
				oi.OrderTitle = value.String
			}
		case orderinfo.FieldOrderType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_type", values[i])
			} else if value.Valid {
				oi.OrderType = value.String
			}
		case orderinfo.FieldOrderDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_description", values[i])
			} else if value.Valid {
				oi.OrderDescription = value.String
			}
		case orderinfo.FieldOrderConfirmationNumber:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_confirmation_number", values[i])
			} else if value.Valid {
				oi.OrderConfirmationNumber = value.String
			}
		case orderinfo.FieldClinicID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field clinic_id", values[i])
			} else if value.Valid {
				oi.ClinicID = int(value.Int64)
			}
		case orderinfo.FieldCustomerID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field customer_id", values[i])
			} else if value.Valid {
				oi.CustomerID = int(value.Int64)
			}
		case orderinfo.FieldOrderCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field order_create_time", values[i])
			} else if value.Valid {
				oi.OrderCreateTime = value.Time
			}
		case orderinfo.FieldOrderServiceTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field order_service_time", values[i])
			} else if value.Valid {
				oi.OrderServiceTime = value.Time
			}
		case orderinfo.FieldOrderProcessTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field order_process_time", values[i])
			} else if value.Valid {
				oi.OrderProcessTime = value.Time
			}
		case orderinfo.FieldOrderRedrawTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field order_redraw_time", values[i])
			} else if value.Valid {
				oi.OrderRedrawTime = value.Time
			}
		case orderinfo.FieldOrderCancelTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field order_cancel_time", values[i])
			} else if value.Valid {
				oi.OrderCancelTime = value.Time
			}
		case orderinfo.FieldIsActive:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field isActive", values[i])
			} else if value.Valid {
				oi.IsActive = value.Bool
			}
		case orderinfo.FieldHasOrderSetting:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field has_order_setting", values[i])
			} else if value.Valid {
				oi.HasOrderSetting = value.Bool
			}
		case orderinfo.FieldOrderCanceled:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field order_canceled", values[i])
			} else if value.Valid {
				oi.OrderCanceled = value.Bool
			}
		case orderinfo.FieldOrderFlagged:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field order_flagged", values[i])
			} else if value.Valid {
				oi.OrderFlagged = value.Bool
			}
		case orderinfo.FieldOrderStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_status", values[i])
			} else if value.Valid {
				oi.OrderStatus = value.String
			}
		case orderinfo.FieldOrderMajorStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_major_status", values[i])
			} else if value.Valid {
				oi.OrderMajorStatus = value.String
			}
		case orderinfo.FieldOrderKitStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_kit_status", values[i])
			} else if value.Valid {
				oi.OrderKitStatus = value.String
			}
		case orderinfo.FieldOrderReportStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_report_status", values[i])
			} else if value.Valid {
				oi.OrderReportStatus = value.String
			}
		case orderinfo.FieldOrderTnpIssueStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_tnp_issue_status", values[i])
			} else if value.Valid {
				oi.OrderTnpIssueStatus = value.String
			}
		case orderinfo.FieldOrderBillingIssueStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_billing_issue_status", values[i])
			} else if value.Valid {
				oi.OrderBillingIssueStatus = value.String
			}
		case orderinfo.FieldOrderMissingInfoIssueStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_missing_info_issue_status", values[i])
			} else if value.Valid {
				oi.OrderMissingInfoIssueStatus = value.String
			}
		case orderinfo.FieldOrderIncompleteQuestionnaireIssueStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_incomplete_questionnaire_issue_status", values[i])
			} else if value.Valid {
				oi.OrderIncompleteQuestionnaireIssueStatus = value.String
			}
		case orderinfo.FieldOrderNyWaiveFormIssueStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_ny_waive_form_issue_status", values[i])
			} else if value.Valid {
				oi.OrderNyWaiveFormIssueStatus = value.String
			}
		case orderinfo.FieldOrderLabIssueStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_lab_issue_status", values[i])
			} else if value.Valid {
				oi.OrderLabIssueStatus = value.String
			}
		case orderinfo.FieldOrderProcessingTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field order_processing_time", values[i])
			} else if value.Valid {
				oi.OrderProcessingTime = value.Time
			}
		case orderinfo.FieldOrderMinorStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_minor_status", values[i])
			} else if value.Valid {
				oi.OrderMinorStatus = value.String
			}
		case orderinfo.FieldPatientFirstName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field patient_first_name", values[i])
			} else if value.Valid {
				oi.PatientFirstName = value.String
			}
		case orderinfo.FieldPatientLastName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field patient_last_name", values[i])
			} else if value.Valid {
				oi.PatientLastName = value.String
			}
		case orderinfo.FieldOrderSource:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_source", values[i])
			} else if value.Valid {
				oi.OrderSource = value.String
			}
		case orderinfo.FieldOrderChargeMethod:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_charge_method", values[i])
			} else if value.Valid {
				oi.OrderChargeMethod = value.String
			}
		case orderinfo.FieldOrderPlacingType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_placing_type", values[i])
			} else if value.Valid {
				oi.OrderPlacingType = value.String
			}
		case orderinfo.FieldBillingOrderID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field billing_order_id", values[i])
			} else if value.Valid {
				oi.BillingOrderID = value.String
			}
		case orderinfo.FieldContactID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field contact_id", values[i])
			} else if value.Valid {
				oi.ContactID = int(value.Int64)
			}
		case orderinfo.FieldAddressID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field address_id", values[i])
			} else if value.Valid {
				oi.AddressID = int(value.Int64)
			}
		default:
			oi.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the OrderInfo.
// This includes values selected through modifiers, order, etc.
func (oi *OrderInfo) Value(name string) (ent.Value, error) {
	return oi.selectValues.Get(name)
}

// QueryTests queries the "tests" edge of the OrderInfo entity.
func (oi *OrderInfo) QueryTests() *TestQuery {
	return NewOrderInfoClient(oi.config).QueryTests(oi)
}

// QueryOrderFlags queries the "order_flags" edge of the OrderInfo entity.
func (oi *OrderInfo) QueryOrderFlags() *OrderFlagQuery {
	return NewOrderInfoClient(oi.config).QueryOrderFlags(oi)
}

// QuerySample queries the "sample" edge of the OrderInfo entity.
func (oi *OrderInfo) QuerySample() *SampleQuery {
	return NewOrderInfoClient(oi.config).QuerySample(oi)
}

// QueryContact queries the "contact" edge of the OrderInfo entity.
func (oi *OrderInfo) QueryContact() *ContactQuery {
	return NewOrderInfoClient(oi.config).QueryContact(oi)
}

// QueryAddress queries the "address" edge of the OrderInfo entity.
func (oi *OrderInfo) QueryAddress() *AddressQuery {
	return NewOrderInfoClient(oi.config).QueryAddress(oi)
}

// QueryClinic queries the "clinic" edge of the OrderInfo entity.
func (oi *OrderInfo) QueryClinic() *ClinicQuery {
	return NewOrderInfoClient(oi.config).QueryClinic(oi)
}

// QueryCustomerInfo queries the "customer_info" edge of the OrderInfo entity.
func (oi *OrderInfo) QueryCustomerInfo() *CustomerQuery {
	return NewOrderInfoClient(oi.config).QueryCustomerInfo(oi)
}

// Update returns a builder for updating this OrderInfo.
// Note that you need to call OrderInfo.Unwrap() before calling this method if this OrderInfo
// was returned from a transaction, and the transaction was committed or rolled back.
func (oi *OrderInfo) Update() *OrderInfoUpdateOne {
	return NewOrderInfoClient(oi.config).UpdateOne(oi)
}

// Unwrap unwraps the OrderInfo entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (oi *OrderInfo) Unwrap() *OrderInfo {
	_tx, ok := oi.config.driver.(*txDriver)
	if !ok {
		panic("ent: OrderInfo is not a transactional entity")
	}
	oi.config.driver = _tx.drv
	return oi
}

// String implements the fmt.Stringer.
func (oi *OrderInfo) String() string {
	var builder strings.Builder
	builder.WriteString("OrderInfo(")
	builder.WriteString(fmt.Sprintf("id=%v, ", oi.ID))
	builder.WriteString("order_title=")
	builder.WriteString(oi.OrderTitle)
	builder.WriteString(", ")
	builder.WriteString("order_type=")
	builder.WriteString(oi.OrderType)
	builder.WriteString(", ")
	builder.WriteString("order_description=")
	builder.WriteString(oi.OrderDescription)
	builder.WriteString(", ")
	builder.WriteString("order_confirmation_number=")
	builder.WriteString(oi.OrderConfirmationNumber)
	builder.WriteString(", ")
	builder.WriteString("clinic_id=")
	builder.WriteString(fmt.Sprintf("%v", oi.ClinicID))
	builder.WriteString(", ")
	builder.WriteString("customer_id=")
	builder.WriteString(fmt.Sprintf("%v", oi.CustomerID))
	builder.WriteString(", ")
	builder.WriteString("order_create_time=")
	builder.WriteString(oi.OrderCreateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("order_service_time=")
	builder.WriteString(oi.OrderServiceTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("order_process_time=")
	builder.WriteString(oi.OrderProcessTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("order_redraw_time=")
	builder.WriteString(oi.OrderRedrawTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("order_cancel_time=")
	builder.WriteString(oi.OrderCancelTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("isActive=")
	builder.WriteString(fmt.Sprintf("%v", oi.IsActive))
	builder.WriteString(", ")
	builder.WriteString("has_order_setting=")
	builder.WriteString(fmt.Sprintf("%v", oi.HasOrderSetting))
	builder.WriteString(", ")
	builder.WriteString("order_canceled=")
	builder.WriteString(fmt.Sprintf("%v", oi.OrderCanceled))
	builder.WriteString(", ")
	builder.WriteString("order_flagged=")
	builder.WriteString(fmt.Sprintf("%v", oi.OrderFlagged))
	builder.WriteString(", ")
	builder.WriteString("order_status=")
	builder.WriteString(oi.OrderStatus)
	builder.WriteString(", ")
	builder.WriteString("order_major_status=")
	builder.WriteString(oi.OrderMajorStatus)
	builder.WriteString(", ")
	builder.WriteString("order_kit_status=")
	builder.WriteString(oi.OrderKitStatus)
	builder.WriteString(", ")
	builder.WriteString("order_report_status=")
	builder.WriteString(oi.OrderReportStatus)
	builder.WriteString(", ")
	builder.WriteString("order_tnp_issue_status=")
	builder.WriteString(oi.OrderTnpIssueStatus)
	builder.WriteString(", ")
	builder.WriteString("order_billing_issue_status=")
	builder.WriteString(oi.OrderBillingIssueStatus)
	builder.WriteString(", ")
	builder.WriteString("order_missing_info_issue_status=")
	builder.WriteString(oi.OrderMissingInfoIssueStatus)
	builder.WriteString(", ")
	builder.WriteString("order_incomplete_questionnaire_issue_status=")
	builder.WriteString(oi.OrderIncompleteQuestionnaireIssueStatus)
	builder.WriteString(", ")
	builder.WriteString("order_ny_waive_form_issue_status=")
	builder.WriteString(oi.OrderNyWaiveFormIssueStatus)
	builder.WriteString(", ")
	builder.WriteString("order_lab_issue_status=")
	builder.WriteString(oi.OrderLabIssueStatus)
	builder.WriteString(", ")
	builder.WriteString("order_processing_time=")
	builder.WriteString(oi.OrderProcessingTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("order_minor_status=")
	builder.WriteString(oi.OrderMinorStatus)
	builder.WriteString(", ")
	builder.WriteString("patient_first_name=")
	builder.WriteString(oi.PatientFirstName)
	builder.WriteString(", ")
	builder.WriteString("patient_last_name=")
	builder.WriteString(oi.PatientLastName)
	builder.WriteString(", ")
	builder.WriteString("order_source=")
	builder.WriteString(oi.OrderSource)
	builder.WriteString(", ")
	builder.WriteString("order_charge_method=")
	builder.WriteString(oi.OrderChargeMethod)
	builder.WriteString(", ")
	builder.WriteString("order_placing_type=")
	builder.WriteString(oi.OrderPlacingType)
	builder.WriteString(", ")
	builder.WriteString("billing_order_id=")
	builder.WriteString(oi.BillingOrderID)
	builder.WriteString(", ")
	builder.WriteString("contact_id=")
	builder.WriteString(fmt.Sprintf("%v", oi.ContactID))
	builder.WriteString(", ")
	builder.WriteString("address_id=")
	builder.WriteString(fmt.Sprintf("%v", oi.AddressID))
	builder.WriteByte(')')
	return builder.String()
}

// OrderInfos is a parsable slice of OrderInfo.
type OrderInfos []*OrderInfo
