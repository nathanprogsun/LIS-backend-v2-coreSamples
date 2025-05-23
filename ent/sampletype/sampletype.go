// Code generated by ent, DO NOT EDIT.

package sampletype

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the sampletype type in the database.
	Label = "sample_type"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "sample_type_id"
	// FieldSampleTypeName holds the string denoting the sample_type_name field in the database.
	FieldSampleTypeName = "sample_type_name"
	// FieldSampleTypeCode holds the string denoting the sample_type_code field in the database.
	FieldSampleTypeCode = "sample_type_code"
	// FieldSampleTypeEnum holds the string denoting the sample_type_enum field in the database.
	FieldSampleTypeEnum = "sample_type_emun"
	// FieldSampleTypeEnumOldLisRequest holds the string denoting the sample_type_enum_old_lis_request field in the database.
	FieldSampleTypeEnumOldLisRequest = "sample_type_emun_old_lis_request"
	// FieldSampleTypeDescription holds the string denoting the sample_type_description field in the database.
	FieldSampleTypeDescription = "sample_type_description"
	// FieldPrimarySampleTypeGroup holds the string denoting the primary_sample_type_group field in the database.
	FieldPrimarySampleTypeGroup = "primary_sample_type_group"
	// FieldIsActive holds the string denoting the is_active field in the database.
	FieldIsActive = "isActive"
	// FieldCreatedTime holds the string denoting the created_time field in the database.
	FieldCreatedTime = "created_time"
	// FieldUpdatedTime holds the string denoting the updated_time field in the database.
	FieldUpdatedTime = "updated_time"
	// EdgeTubeTypes holds the string denoting the tube_types edge name in mutations.
	EdgeTubeTypes = "tube_types"
	// EdgeTests holds the string denoting the tests edge name in mutations.
	EdgeTests = "tests"
	// TubeTypeFieldID holds the string denoting the ID field of the TubeType.
	TubeTypeFieldID = "id"
	// TestFieldID holds the string denoting the ID field of the Test.
	TestFieldID = "test_id"
	// Table holds the table name of the sampletype in the database.
	Table = "sample_type"
	// TubeTypesTable is the table that holds the tube_types relation/edge. The primary key declared below.
	TubeTypesTable = "_sample_type_to_tube_type"
	// TubeTypesInverseTable is the table name for the TubeType entity.
	// It exists in this package in order to avoid circular dependency with the "tubetype" package.
	TubeTypesInverseTable = "tube_type"
	// TestsTable is the table that holds the tests relation/edge. The primary key declared below.
	TestsTable = "_sample_type_to_test"
	// TestsInverseTable is the table name for the Test entity.
	// It exists in this package in order to avoid circular dependency with the "test" package.
	TestsInverseTable = "test"
)

// Columns holds all SQL columns for sampletype fields.
var Columns = []string{
	FieldID,
	FieldSampleTypeName,
	FieldSampleTypeCode,
	FieldSampleTypeEnum,
	FieldSampleTypeEnumOldLisRequest,
	FieldSampleTypeDescription,
	FieldPrimarySampleTypeGroup,
	FieldIsActive,
	FieldCreatedTime,
	FieldUpdatedTime,
}

var (
	// TubeTypesPrimaryKey and TubeTypesColumn2 are the table columns denoting the
	// primary key for the tube_types relation (M2M).
	TubeTypesPrimaryKey = []string{"sample_type_id", "tube_type_id"}
	// TestsPrimaryKey and TestsColumn2 are the table columns denoting the
	// primary key for the tests relation (M2M).
	TestsPrimaryKey = []string{"sample_type_id", "test_id"}
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
	// DefaultIsActive holds the default value on creation for the "is_active" field.
	DefaultIsActive bool
	// DefaultCreatedTime holds the default value on creation for the "created_time" field.
	DefaultCreatedTime func() time.Time
	// UpdateDefaultUpdatedTime holds the default value on update for the "updated_time" field.
	UpdateDefaultUpdatedTime func() time.Time
)

// OrderOption defines the ordering options for the SampleType queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// BySampleTypeName orders the results by the sample_type_name field.
func BySampleTypeName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSampleTypeName, opts...).ToFunc()
}

// BySampleTypeCode orders the results by the sample_type_code field.
func BySampleTypeCode(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSampleTypeCode, opts...).ToFunc()
}

// BySampleTypeEnum orders the results by the sample_type_enum field.
func BySampleTypeEnum(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSampleTypeEnum, opts...).ToFunc()
}

// BySampleTypeEnumOldLisRequest orders the results by the sample_type_enum_old_lis_request field.
func BySampleTypeEnumOldLisRequest(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSampleTypeEnumOldLisRequest, opts...).ToFunc()
}

// BySampleTypeDescription orders the results by the sample_type_description field.
func BySampleTypeDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSampleTypeDescription, opts...).ToFunc()
}

// ByPrimarySampleTypeGroup orders the results by the primary_sample_type_group field.
func ByPrimarySampleTypeGroup(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPrimarySampleTypeGroup, opts...).ToFunc()
}

// ByIsActive orders the results by the is_active field.
func ByIsActive(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsActive, opts...).ToFunc()
}

// ByCreatedTime orders the results by the created_time field.
func ByCreatedTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedTime, opts...).ToFunc()
}

// ByUpdatedTime orders the results by the updated_time field.
func ByUpdatedTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedTime, opts...).ToFunc()
}

// ByTubeTypesCount orders the results by tube_types count.
func ByTubeTypesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newTubeTypesStep(), opts...)
	}
}

// ByTubeTypes orders the results by tube_types terms.
func ByTubeTypes(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newTubeTypesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByTestsCount orders the results by tests count.
func ByTestsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newTestsStep(), opts...)
	}
}

// ByTests orders the results by tests terms.
func ByTests(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newTestsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newTubeTypesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(TubeTypesInverseTable, TubeTypeFieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, TubeTypesTable, TubeTypesPrimaryKey...),
	)
}
func newTestsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(TestsInverseTable, TestFieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, TestsTable, TestsPrimaryKey...),
	)
}
