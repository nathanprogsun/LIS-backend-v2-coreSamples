// Code generated by ent, DO NOT EDIT.

package tubetype

import (
	"coresamples/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.TubeType {
	return predicate.TubeType(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.TubeType {
	return predicate.TubeType(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.TubeType {
	return predicate.TubeType(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.TubeType {
	return predicate.TubeType(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.TubeType {
	return predicate.TubeType(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.TubeType {
	return predicate.TubeType(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.TubeType {
	return predicate.TubeType(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.TubeType {
	return predicate.TubeType(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.TubeType {
	return predicate.TubeType(sql.FieldLTE(FieldID, id))
}

// TubeName applies equality check predicate on the "tube_name" field. It's identical to TubeNameEQ.
func TubeName(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldEQ(FieldTubeName, v))
}

// TubeTypeEnum applies equality check predicate on the "tube_type_enum" field. It's identical to TubeTypeEnumEQ.
func TubeTypeEnum(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldEQ(FieldTubeTypeEnum, v))
}

// TubeTypeSymbol applies equality check predicate on the "tube_type_symbol" field. It's identical to TubeTypeSymbolEQ.
func TubeTypeSymbol(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldEQ(FieldTubeTypeSymbol, v))
}

// TubeNameEQ applies the EQ predicate on the "tube_name" field.
func TubeNameEQ(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldEQ(FieldTubeName, v))
}

// TubeNameNEQ applies the NEQ predicate on the "tube_name" field.
func TubeNameNEQ(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldNEQ(FieldTubeName, v))
}

// TubeNameIn applies the In predicate on the "tube_name" field.
func TubeNameIn(vs ...string) predicate.TubeType {
	return predicate.TubeType(sql.FieldIn(FieldTubeName, vs...))
}

// TubeNameNotIn applies the NotIn predicate on the "tube_name" field.
func TubeNameNotIn(vs ...string) predicate.TubeType {
	return predicate.TubeType(sql.FieldNotIn(FieldTubeName, vs...))
}

// TubeNameGT applies the GT predicate on the "tube_name" field.
func TubeNameGT(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldGT(FieldTubeName, v))
}

// TubeNameGTE applies the GTE predicate on the "tube_name" field.
func TubeNameGTE(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldGTE(FieldTubeName, v))
}

// TubeNameLT applies the LT predicate on the "tube_name" field.
func TubeNameLT(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldLT(FieldTubeName, v))
}

// TubeNameLTE applies the LTE predicate on the "tube_name" field.
func TubeNameLTE(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldLTE(FieldTubeName, v))
}

// TubeNameContains applies the Contains predicate on the "tube_name" field.
func TubeNameContains(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldContains(FieldTubeName, v))
}

// TubeNameHasPrefix applies the HasPrefix predicate on the "tube_name" field.
func TubeNameHasPrefix(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldHasPrefix(FieldTubeName, v))
}

// TubeNameHasSuffix applies the HasSuffix predicate on the "tube_name" field.
func TubeNameHasSuffix(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldHasSuffix(FieldTubeName, v))
}

// TubeNameEqualFold applies the EqualFold predicate on the "tube_name" field.
func TubeNameEqualFold(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldEqualFold(FieldTubeName, v))
}

// TubeNameContainsFold applies the ContainsFold predicate on the "tube_name" field.
func TubeNameContainsFold(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldContainsFold(FieldTubeName, v))
}

// TubeTypeEnumEQ applies the EQ predicate on the "tube_type_enum" field.
func TubeTypeEnumEQ(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldEQ(FieldTubeTypeEnum, v))
}

// TubeTypeEnumNEQ applies the NEQ predicate on the "tube_type_enum" field.
func TubeTypeEnumNEQ(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldNEQ(FieldTubeTypeEnum, v))
}

// TubeTypeEnumIn applies the In predicate on the "tube_type_enum" field.
func TubeTypeEnumIn(vs ...string) predicate.TubeType {
	return predicate.TubeType(sql.FieldIn(FieldTubeTypeEnum, vs...))
}

// TubeTypeEnumNotIn applies the NotIn predicate on the "tube_type_enum" field.
func TubeTypeEnumNotIn(vs ...string) predicate.TubeType {
	return predicate.TubeType(sql.FieldNotIn(FieldTubeTypeEnum, vs...))
}

// TubeTypeEnumGT applies the GT predicate on the "tube_type_enum" field.
func TubeTypeEnumGT(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldGT(FieldTubeTypeEnum, v))
}

// TubeTypeEnumGTE applies the GTE predicate on the "tube_type_enum" field.
func TubeTypeEnumGTE(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldGTE(FieldTubeTypeEnum, v))
}

// TubeTypeEnumLT applies the LT predicate on the "tube_type_enum" field.
func TubeTypeEnumLT(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldLT(FieldTubeTypeEnum, v))
}

// TubeTypeEnumLTE applies the LTE predicate on the "tube_type_enum" field.
func TubeTypeEnumLTE(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldLTE(FieldTubeTypeEnum, v))
}

// TubeTypeEnumContains applies the Contains predicate on the "tube_type_enum" field.
func TubeTypeEnumContains(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldContains(FieldTubeTypeEnum, v))
}

// TubeTypeEnumHasPrefix applies the HasPrefix predicate on the "tube_type_enum" field.
func TubeTypeEnumHasPrefix(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldHasPrefix(FieldTubeTypeEnum, v))
}

// TubeTypeEnumHasSuffix applies the HasSuffix predicate on the "tube_type_enum" field.
func TubeTypeEnumHasSuffix(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldHasSuffix(FieldTubeTypeEnum, v))
}

// TubeTypeEnumEqualFold applies the EqualFold predicate on the "tube_type_enum" field.
func TubeTypeEnumEqualFold(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldEqualFold(FieldTubeTypeEnum, v))
}

// TubeTypeEnumContainsFold applies the ContainsFold predicate on the "tube_type_enum" field.
func TubeTypeEnumContainsFold(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldContainsFold(FieldTubeTypeEnum, v))
}

// TubeTypeSymbolEQ applies the EQ predicate on the "tube_type_symbol" field.
func TubeTypeSymbolEQ(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldEQ(FieldTubeTypeSymbol, v))
}

// TubeTypeSymbolNEQ applies the NEQ predicate on the "tube_type_symbol" field.
func TubeTypeSymbolNEQ(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldNEQ(FieldTubeTypeSymbol, v))
}

// TubeTypeSymbolIn applies the In predicate on the "tube_type_symbol" field.
func TubeTypeSymbolIn(vs ...string) predicate.TubeType {
	return predicate.TubeType(sql.FieldIn(FieldTubeTypeSymbol, vs...))
}

// TubeTypeSymbolNotIn applies the NotIn predicate on the "tube_type_symbol" field.
func TubeTypeSymbolNotIn(vs ...string) predicate.TubeType {
	return predicate.TubeType(sql.FieldNotIn(FieldTubeTypeSymbol, vs...))
}

// TubeTypeSymbolGT applies the GT predicate on the "tube_type_symbol" field.
func TubeTypeSymbolGT(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldGT(FieldTubeTypeSymbol, v))
}

// TubeTypeSymbolGTE applies the GTE predicate on the "tube_type_symbol" field.
func TubeTypeSymbolGTE(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldGTE(FieldTubeTypeSymbol, v))
}

// TubeTypeSymbolLT applies the LT predicate on the "tube_type_symbol" field.
func TubeTypeSymbolLT(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldLT(FieldTubeTypeSymbol, v))
}

// TubeTypeSymbolLTE applies the LTE predicate on the "tube_type_symbol" field.
func TubeTypeSymbolLTE(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldLTE(FieldTubeTypeSymbol, v))
}

// TubeTypeSymbolContains applies the Contains predicate on the "tube_type_symbol" field.
func TubeTypeSymbolContains(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldContains(FieldTubeTypeSymbol, v))
}

// TubeTypeSymbolHasPrefix applies the HasPrefix predicate on the "tube_type_symbol" field.
func TubeTypeSymbolHasPrefix(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldHasPrefix(FieldTubeTypeSymbol, v))
}

// TubeTypeSymbolHasSuffix applies the HasSuffix predicate on the "tube_type_symbol" field.
func TubeTypeSymbolHasSuffix(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldHasSuffix(FieldTubeTypeSymbol, v))
}

// TubeTypeSymbolEqualFold applies the EqualFold predicate on the "tube_type_symbol" field.
func TubeTypeSymbolEqualFold(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldEqualFold(FieldTubeTypeSymbol, v))
}

// TubeTypeSymbolContainsFold applies the ContainsFold predicate on the "tube_type_symbol" field.
func TubeTypeSymbolContainsFold(v string) predicate.TubeType {
	return predicate.TubeType(sql.FieldContainsFold(FieldTubeTypeSymbol, v))
}

// HasTube applies the HasEdge predicate on the "tube" edge.
func HasTube() predicate.TubeType {
	return predicate.TubeType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, TubeTable, TubePrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTubeWith applies the HasEdge predicate on the "tube" edge with a given conditions (other predicates).
func HasTubeWith(preds ...predicate.Tube) predicate.TubeType {
	return predicate.TubeType(func(s *sql.Selector) {
		step := newTubeStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasSampleTypes applies the HasEdge predicate on the "sample_types" edge.
func HasSampleTypes() predicate.TubeType {
	return predicate.TubeType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, SampleTypesTable, SampleTypesPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSampleTypesWith applies the HasEdge predicate on the "sample_types" edge with a given conditions (other predicates).
func HasSampleTypesWith(preds ...predicate.SampleType) predicate.TubeType {
	return predicate.TubeType(func(s *sql.Selector) {
		step := newSampleTypesStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasTests applies the HasEdge predicate on the "tests" edge.
func HasTests() predicate.TubeType {
	return predicate.TubeType(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, TestsTable, TestsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTestsWith applies the HasEdge predicate on the "tests" edge with a given conditions (other predicates).
func HasTestsWith(preds ...predicate.Test) predicate.TubeType {
	return predicate.TubeType(func(s *sql.Selector) {
		step := newTestsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.TubeType) predicate.TubeType {
	return predicate.TubeType(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.TubeType) predicate.TubeType {
	return predicate.TubeType(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.TubeType) predicate.TubeType {
	return predicate.TubeType(sql.NotPredicates(p))
}
