// Code generated by ent, DO NOT EDIT.

package pendingordercredits

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the pendingordercredits type in the database.
	Label = "pending_order_credits"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldOrderID holds the string denoting the order_id field in the database.
	FieldOrderID = "order_id"
	// FieldCreditID holds the string denoting the credit_id field in the database.
	FieldCreditID = "credit_id"
	// FieldClinicID holds the string denoting the clinic_id field in the database.
	FieldClinicID = "clinic_id"
	// Table holds the table name of the pendingordercredits in the database.
	Table = "pending_order_credits"
)

// Columns holds all SQL columns for pendingordercredits fields.
var Columns = []string{
	FieldID,
	FieldOrderID,
	FieldCreditID,
	FieldClinicID,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the PendingOrderCredits queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByOrderID orders the results by the order_id field.
func ByOrderID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOrderID, opts...).ToFunc()
}

// ByCreditID orders the results by the credit_id field.
func ByCreditID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreditID, opts...).ToFunc()
}

// ByClinicID orders the results by the clinic_id field.
func ByClinicID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldClinicID, opts...).ToFunc()
}
