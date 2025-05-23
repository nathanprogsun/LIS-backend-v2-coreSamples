// Code generated by ent, DO NOT EDIT.

package tuberequirement

import (
	"coresamples/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldLTE(FieldID, id))
}

// SampleID applies equality check predicate on the "sample_id" field. It's identical to SampleIDEQ.
func SampleID(v int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEQ(FieldSampleID, v))
}

// TubeType applies equality check predicate on the "tube_type" field. It's identical to TubeTypeEQ.
func TubeType(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEQ(FieldTubeType, v))
}

// RequiredCount applies equality check predicate on the "required_count" field. It's identical to RequiredCountEQ.
func RequiredCount(v int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEQ(FieldRequiredCount, v))
}

// RequiredCountCreateTime applies equality check predicate on the "required_count_create_time" field. It's identical to RequiredCountCreateTimeEQ.
func RequiredCountCreateTime(v time.Time) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEQ(FieldRequiredCountCreateTime, v))
}

// RequiredBy applies equality check predicate on the "required_by" field. It's identical to RequiredByEQ.
func RequiredBy(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEQ(FieldRequiredBy, v))
}

// ModifiedBy applies equality check predicate on the "modified_by" field. It's identical to ModifiedByEQ.
func ModifiedBy(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEQ(FieldModifiedBy, v))
}

// ModifiedTime applies equality check predicate on the "modified_time" field. It's identical to ModifiedTimeEQ.
func ModifiedTime(v time.Time) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEQ(FieldModifiedTime, v))
}

// SampleIDEQ applies the EQ predicate on the "sample_id" field.
func SampleIDEQ(v int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEQ(FieldSampleID, v))
}

// SampleIDNEQ applies the NEQ predicate on the "sample_id" field.
func SampleIDNEQ(v int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldNEQ(FieldSampleID, v))
}

// SampleIDIn applies the In predicate on the "sample_id" field.
func SampleIDIn(vs ...int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldIn(FieldSampleID, vs...))
}

// SampleIDNotIn applies the NotIn predicate on the "sample_id" field.
func SampleIDNotIn(vs ...int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldNotIn(FieldSampleID, vs...))
}

// SampleIDIsNil applies the IsNil predicate on the "sample_id" field.
func SampleIDIsNil() predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldIsNull(FieldSampleID))
}

// SampleIDNotNil applies the NotNil predicate on the "sample_id" field.
func SampleIDNotNil() predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldNotNull(FieldSampleID))
}

// TubeTypeEQ applies the EQ predicate on the "tube_type" field.
func TubeTypeEQ(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEQ(FieldTubeType, v))
}

// TubeTypeNEQ applies the NEQ predicate on the "tube_type" field.
func TubeTypeNEQ(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldNEQ(FieldTubeType, v))
}

// TubeTypeIn applies the In predicate on the "tube_type" field.
func TubeTypeIn(vs ...string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldIn(FieldTubeType, vs...))
}

// TubeTypeNotIn applies the NotIn predicate on the "tube_type" field.
func TubeTypeNotIn(vs ...string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldNotIn(FieldTubeType, vs...))
}

// TubeTypeGT applies the GT predicate on the "tube_type" field.
func TubeTypeGT(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldGT(FieldTubeType, v))
}

// TubeTypeGTE applies the GTE predicate on the "tube_type" field.
func TubeTypeGTE(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldGTE(FieldTubeType, v))
}

// TubeTypeLT applies the LT predicate on the "tube_type" field.
func TubeTypeLT(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldLT(FieldTubeType, v))
}

// TubeTypeLTE applies the LTE predicate on the "tube_type" field.
func TubeTypeLTE(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldLTE(FieldTubeType, v))
}

// TubeTypeContains applies the Contains predicate on the "tube_type" field.
func TubeTypeContains(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldContains(FieldTubeType, v))
}

// TubeTypeHasPrefix applies the HasPrefix predicate on the "tube_type" field.
func TubeTypeHasPrefix(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldHasPrefix(FieldTubeType, v))
}

// TubeTypeHasSuffix applies the HasSuffix predicate on the "tube_type" field.
func TubeTypeHasSuffix(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldHasSuffix(FieldTubeType, v))
}

// TubeTypeEqualFold applies the EqualFold predicate on the "tube_type" field.
func TubeTypeEqualFold(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEqualFold(FieldTubeType, v))
}

// TubeTypeContainsFold applies the ContainsFold predicate on the "tube_type" field.
func TubeTypeContainsFold(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldContainsFold(FieldTubeType, v))
}

// RequiredCountEQ applies the EQ predicate on the "required_count" field.
func RequiredCountEQ(v int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEQ(FieldRequiredCount, v))
}

// RequiredCountNEQ applies the NEQ predicate on the "required_count" field.
func RequiredCountNEQ(v int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldNEQ(FieldRequiredCount, v))
}

// RequiredCountIn applies the In predicate on the "required_count" field.
func RequiredCountIn(vs ...int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldIn(FieldRequiredCount, vs...))
}

// RequiredCountNotIn applies the NotIn predicate on the "required_count" field.
func RequiredCountNotIn(vs ...int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldNotIn(FieldRequiredCount, vs...))
}

// RequiredCountGT applies the GT predicate on the "required_count" field.
func RequiredCountGT(v int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldGT(FieldRequiredCount, v))
}

// RequiredCountGTE applies the GTE predicate on the "required_count" field.
func RequiredCountGTE(v int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldGTE(FieldRequiredCount, v))
}

// RequiredCountLT applies the LT predicate on the "required_count" field.
func RequiredCountLT(v int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldLT(FieldRequiredCount, v))
}

// RequiredCountLTE applies the LTE predicate on the "required_count" field.
func RequiredCountLTE(v int) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldLTE(FieldRequiredCount, v))
}

// RequiredCountCreateTimeEQ applies the EQ predicate on the "required_count_create_time" field.
func RequiredCountCreateTimeEQ(v time.Time) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEQ(FieldRequiredCountCreateTime, v))
}

// RequiredCountCreateTimeNEQ applies the NEQ predicate on the "required_count_create_time" field.
func RequiredCountCreateTimeNEQ(v time.Time) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldNEQ(FieldRequiredCountCreateTime, v))
}

// RequiredCountCreateTimeIn applies the In predicate on the "required_count_create_time" field.
func RequiredCountCreateTimeIn(vs ...time.Time) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldIn(FieldRequiredCountCreateTime, vs...))
}

// RequiredCountCreateTimeNotIn applies the NotIn predicate on the "required_count_create_time" field.
func RequiredCountCreateTimeNotIn(vs ...time.Time) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldNotIn(FieldRequiredCountCreateTime, vs...))
}

// RequiredCountCreateTimeGT applies the GT predicate on the "required_count_create_time" field.
func RequiredCountCreateTimeGT(v time.Time) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldGT(FieldRequiredCountCreateTime, v))
}

// RequiredCountCreateTimeGTE applies the GTE predicate on the "required_count_create_time" field.
func RequiredCountCreateTimeGTE(v time.Time) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldGTE(FieldRequiredCountCreateTime, v))
}

// RequiredCountCreateTimeLT applies the LT predicate on the "required_count_create_time" field.
func RequiredCountCreateTimeLT(v time.Time) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldLT(FieldRequiredCountCreateTime, v))
}

// RequiredCountCreateTimeLTE applies the LTE predicate on the "required_count_create_time" field.
func RequiredCountCreateTimeLTE(v time.Time) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldLTE(FieldRequiredCountCreateTime, v))
}

// RequiredByEQ applies the EQ predicate on the "required_by" field.
func RequiredByEQ(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEQ(FieldRequiredBy, v))
}

// RequiredByNEQ applies the NEQ predicate on the "required_by" field.
func RequiredByNEQ(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldNEQ(FieldRequiredBy, v))
}

// RequiredByIn applies the In predicate on the "required_by" field.
func RequiredByIn(vs ...string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldIn(FieldRequiredBy, vs...))
}

// RequiredByNotIn applies the NotIn predicate on the "required_by" field.
func RequiredByNotIn(vs ...string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldNotIn(FieldRequiredBy, vs...))
}

// RequiredByGT applies the GT predicate on the "required_by" field.
func RequiredByGT(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldGT(FieldRequiredBy, v))
}

// RequiredByGTE applies the GTE predicate on the "required_by" field.
func RequiredByGTE(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldGTE(FieldRequiredBy, v))
}

// RequiredByLT applies the LT predicate on the "required_by" field.
func RequiredByLT(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldLT(FieldRequiredBy, v))
}

// RequiredByLTE applies the LTE predicate on the "required_by" field.
func RequiredByLTE(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldLTE(FieldRequiredBy, v))
}

// RequiredByContains applies the Contains predicate on the "required_by" field.
func RequiredByContains(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldContains(FieldRequiredBy, v))
}

// RequiredByHasPrefix applies the HasPrefix predicate on the "required_by" field.
func RequiredByHasPrefix(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldHasPrefix(FieldRequiredBy, v))
}

// RequiredByHasSuffix applies the HasSuffix predicate on the "required_by" field.
func RequiredByHasSuffix(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldHasSuffix(FieldRequiredBy, v))
}

// RequiredByEqualFold applies the EqualFold predicate on the "required_by" field.
func RequiredByEqualFold(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEqualFold(FieldRequiredBy, v))
}

// RequiredByContainsFold applies the ContainsFold predicate on the "required_by" field.
func RequiredByContainsFold(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldContainsFold(FieldRequiredBy, v))
}

// ModifiedByEQ applies the EQ predicate on the "modified_by" field.
func ModifiedByEQ(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEQ(FieldModifiedBy, v))
}

// ModifiedByNEQ applies the NEQ predicate on the "modified_by" field.
func ModifiedByNEQ(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldNEQ(FieldModifiedBy, v))
}

// ModifiedByIn applies the In predicate on the "modified_by" field.
func ModifiedByIn(vs ...string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldIn(FieldModifiedBy, vs...))
}

// ModifiedByNotIn applies the NotIn predicate on the "modified_by" field.
func ModifiedByNotIn(vs ...string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldNotIn(FieldModifiedBy, vs...))
}

// ModifiedByGT applies the GT predicate on the "modified_by" field.
func ModifiedByGT(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldGT(FieldModifiedBy, v))
}

// ModifiedByGTE applies the GTE predicate on the "modified_by" field.
func ModifiedByGTE(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldGTE(FieldModifiedBy, v))
}

// ModifiedByLT applies the LT predicate on the "modified_by" field.
func ModifiedByLT(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldLT(FieldModifiedBy, v))
}

// ModifiedByLTE applies the LTE predicate on the "modified_by" field.
func ModifiedByLTE(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldLTE(FieldModifiedBy, v))
}

// ModifiedByContains applies the Contains predicate on the "modified_by" field.
func ModifiedByContains(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldContains(FieldModifiedBy, v))
}

// ModifiedByHasPrefix applies the HasPrefix predicate on the "modified_by" field.
func ModifiedByHasPrefix(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldHasPrefix(FieldModifiedBy, v))
}

// ModifiedByHasSuffix applies the HasSuffix predicate on the "modified_by" field.
func ModifiedByHasSuffix(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldHasSuffix(FieldModifiedBy, v))
}

// ModifiedByIsNil applies the IsNil predicate on the "modified_by" field.
func ModifiedByIsNil() predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldIsNull(FieldModifiedBy))
}

// ModifiedByNotNil applies the NotNil predicate on the "modified_by" field.
func ModifiedByNotNil() predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldNotNull(FieldModifiedBy))
}

// ModifiedByEqualFold applies the EqualFold predicate on the "modified_by" field.
func ModifiedByEqualFold(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEqualFold(FieldModifiedBy, v))
}

// ModifiedByContainsFold applies the ContainsFold predicate on the "modified_by" field.
func ModifiedByContainsFold(v string) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldContainsFold(FieldModifiedBy, v))
}

// ModifiedTimeEQ applies the EQ predicate on the "modified_time" field.
func ModifiedTimeEQ(v time.Time) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldEQ(FieldModifiedTime, v))
}

// ModifiedTimeNEQ applies the NEQ predicate on the "modified_time" field.
func ModifiedTimeNEQ(v time.Time) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldNEQ(FieldModifiedTime, v))
}

// ModifiedTimeIn applies the In predicate on the "modified_time" field.
func ModifiedTimeIn(vs ...time.Time) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldIn(FieldModifiedTime, vs...))
}

// ModifiedTimeNotIn applies the NotIn predicate on the "modified_time" field.
func ModifiedTimeNotIn(vs ...time.Time) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldNotIn(FieldModifiedTime, vs...))
}

// ModifiedTimeGT applies the GT predicate on the "modified_time" field.
func ModifiedTimeGT(v time.Time) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldGT(FieldModifiedTime, v))
}

// ModifiedTimeGTE applies the GTE predicate on the "modified_time" field.
func ModifiedTimeGTE(v time.Time) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldGTE(FieldModifiedTime, v))
}

// ModifiedTimeLT applies the LT predicate on the "modified_time" field.
func ModifiedTimeLT(v time.Time) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldLT(FieldModifiedTime, v))
}

// ModifiedTimeLTE applies the LTE predicate on the "modified_time" field.
func ModifiedTimeLTE(v time.Time) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.FieldLTE(FieldModifiedTime, v))
}

// HasSample applies the HasEdge predicate on the "sample" edge.
func HasSample() predicate.TubeRequirement {
	return predicate.TubeRequirement(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, SampleTable, SampleColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSampleWith applies the HasEdge predicate on the "sample" edge with a given conditions (other predicates).
func HasSampleWith(preds ...predicate.Sample) predicate.TubeRequirement {
	return predicate.TubeRequirement(func(s *sql.Selector) {
		step := newSampleStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.TubeRequirement) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.TubeRequirement) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.TubeRequirement) predicate.TubeRequirement {
	return predicate.TubeRequirement(sql.NotPredicates(p))
}
