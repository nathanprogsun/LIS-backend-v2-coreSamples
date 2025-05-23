// Code generated by ent, DO NOT EDIT.

package betaprogramparticipation

import (
	"coresamples/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldLTE(FieldID, id))
}

// BetaProgramID applies equality check predicate on the "beta_program_id" field. It's identical to BetaProgramIDEQ.
func BetaProgramID(v int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldEQ(FieldBetaProgramID, v))
}

// CustomerID applies equality check predicate on the "customer_id" field. It's identical to CustomerIDEQ.
func CustomerID(v int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldEQ(FieldCustomerID, v))
}

// ClinicID applies equality check predicate on the "clinic_id" field. It's identical to ClinicIDEQ.
func ClinicID(v int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldEQ(FieldClinicID, v))
}

// IsActive applies equality check predicate on the "is_active" field. It's identical to IsActiveEQ.
func IsActive(v bool) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldEQ(FieldIsActive, v))
}

// HasModifiedStartTime applies equality check predicate on the "has_modified_start_time" field. It's identical to HasModifiedStartTimeEQ.
func HasModifiedStartTime(v bool) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldEQ(FieldHasModifiedStartTime, v))
}

// ModifiedStartTime applies equality check predicate on the "modified_start_time" field. It's identical to ModifiedStartTimeEQ.
func ModifiedStartTime(v time.Time) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldEQ(FieldModifiedStartTime, v))
}

// ModifiedEndTime applies equality check predicate on the "modified_end_time" field. It's identical to ModifiedEndTimeEQ.
func ModifiedEndTime(v time.Time) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldEQ(FieldModifiedEndTime, v))
}

// BetaProgramIDEQ applies the EQ predicate on the "beta_program_id" field.
func BetaProgramIDEQ(v int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldEQ(FieldBetaProgramID, v))
}

// BetaProgramIDNEQ applies the NEQ predicate on the "beta_program_id" field.
func BetaProgramIDNEQ(v int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldNEQ(FieldBetaProgramID, v))
}

// BetaProgramIDIn applies the In predicate on the "beta_program_id" field.
func BetaProgramIDIn(vs ...int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldIn(FieldBetaProgramID, vs...))
}

// BetaProgramIDNotIn applies the NotIn predicate on the "beta_program_id" field.
func BetaProgramIDNotIn(vs ...int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldNotIn(FieldBetaProgramID, vs...))
}

// CustomerIDEQ applies the EQ predicate on the "customer_id" field.
func CustomerIDEQ(v int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldEQ(FieldCustomerID, v))
}

// CustomerIDNEQ applies the NEQ predicate on the "customer_id" field.
func CustomerIDNEQ(v int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldNEQ(FieldCustomerID, v))
}

// CustomerIDIn applies the In predicate on the "customer_id" field.
func CustomerIDIn(vs ...int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldIn(FieldCustomerID, vs...))
}

// CustomerIDNotIn applies the NotIn predicate on the "customer_id" field.
func CustomerIDNotIn(vs ...int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldNotIn(FieldCustomerID, vs...))
}

// ClinicIDEQ applies the EQ predicate on the "clinic_id" field.
func ClinicIDEQ(v int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldEQ(FieldClinicID, v))
}

// ClinicIDNEQ applies the NEQ predicate on the "clinic_id" field.
func ClinicIDNEQ(v int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldNEQ(FieldClinicID, v))
}

// ClinicIDIn applies the In predicate on the "clinic_id" field.
func ClinicIDIn(vs ...int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldIn(FieldClinicID, vs...))
}

// ClinicIDNotIn applies the NotIn predicate on the "clinic_id" field.
func ClinicIDNotIn(vs ...int) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldNotIn(FieldClinicID, vs...))
}

// IsActiveEQ applies the EQ predicate on the "is_active" field.
func IsActiveEQ(v bool) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldEQ(FieldIsActive, v))
}

// IsActiveNEQ applies the NEQ predicate on the "is_active" field.
func IsActiveNEQ(v bool) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldNEQ(FieldIsActive, v))
}

// HasModifiedStartTimeEQ applies the EQ predicate on the "has_modified_start_time" field.
func HasModifiedStartTimeEQ(v bool) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldEQ(FieldHasModifiedStartTime, v))
}

// HasModifiedStartTimeNEQ applies the NEQ predicate on the "has_modified_start_time" field.
func HasModifiedStartTimeNEQ(v bool) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldNEQ(FieldHasModifiedStartTime, v))
}

// ModifiedStartTimeEQ applies the EQ predicate on the "modified_start_time" field.
func ModifiedStartTimeEQ(v time.Time) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldEQ(FieldModifiedStartTime, v))
}

// ModifiedStartTimeNEQ applies the NEQ predicate on the "modified_start_time" field.
func ModifiedStartTimeNEQ(v time.Time) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldNEQ(FieldModifiedStartTime, v))
}

// ModifiedStartTimeIn applies the In predicate on the "modified_start_time" field.
func ModifiedStartTimeIn(vs ...time.Time) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldIn(FieldModifiedStartTime, vs...))
}

// ModifiedStartTimeNotIn applies the NotIn predicate on the "modified_start_time" field.
func ModifiedStartTimeNotIn(vs ...time.Time) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldNotIn(FieldModifiedStartTime, vs...))
}

// ModifiedStartTimeGT applies the GT predicate on the "modified_start_time" field.
func ModifiedStartTimeGT(v time.Time) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldGT(FieldModifiedStartTime, v))
}

// ModifiedStartTimeGTE applies the GTE predicate on the "modified_start_time" field.
func ModifiedStartTimeGTE(v time.Time) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldGTE(FieldModifiedStartTime, v))
}

// ModifiedStartTimeLT applies the LT predicate on the "modified_start_time" field.
func ModifiedStartTimeLT(v time.Time) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldLT(FieldModifiedStartTime, v))
}

// ModifiedStartTimeLTE applies the LTE predicate on the "modified_start_time" field.
func ModifiedStartTimeLTE(v time.Time) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldLTE(FieldModifiedStartTime, v))
}

// ModifiedStartTimeIsNil applies the IsNil predicate on the "modified_start_time" field.
func ModifiedStartTimeIsNil() predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldIsNull(FieldModifiedStartTime))
}

// ModifiedStartTimeNotNil applies the NotNil predicate on the "modified_start_time" field.
func ModifiedStartTimeNotNil() predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldNotNull(FieldModifiedStartTime))
}

// ModifiedEndTimeEQ applies the EQ predicate on the "modified_end_time" field.
func ModifiedEndTimeEQ(v time.Time) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldEQ(FieldModifiedEndTime, v))
}

// ModifiedEndTimeNEQ applies the NEQ predicate on the "modified_end_time" field.
func ModifiedEndTimeNEQ(v time.Time) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldNEQ(FieldModifiedEndTime, v))
}

// ModifiedEndTimeIn applies the In predicate on the "modified_end_time" field.
func ModifiedEndTimeIn(vs ...time.Time) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldIn(FieldModifiedEndTime, vs...))
}

// ModifiedEndTimeNotIn applies the NotIn predicate on the "modified_end_time" field.
func ModifiedEndTimeNotIn(vs ...time.Time) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldNotIn(FieldModifiedEndTime, vs...))
}

// ModifiedEndTimeGT applies the GT predicate on the "modified_end_time" field.
func ModifiedEndTimeGT(v time.Time) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldGT(FieldModifiedEndTime, v))
}

// ModifiedEndTimeGTE applies the GTE predicate on the "modified_end_time" field.
func ModifiedEndTimeGTE(v time.Time) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldGTE(FieldModifiedEndTime, v))
}

// ModifiedEndTimeLT applies the LT predicate on the "modified_end_time" field.
func ModifiedEndTimeLT(v time.Time) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldLT(FieldModifiedEndTime, v))
}

// ModifiedEndTimeLTE applies the LTE predicate on the "modified_end_time" field.
func ModifiedEndTimeLTE(v time.Time) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldLTE(FieldModifiedEndTime, v))
}

// ModifiedEndTimeIsNil applies the IsNil predicate on the "modified_end_time" field.
func ModifiedEndTimeIsNil() predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldIsNull(FieldModifiedEndTime))
}

// ModifiedEndTimeNotNil applies the NotNil predicate on the "modified_end_time" field.
func ModifiedEndTimeNotNil() predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.FieldNotNull(FieldModifiedEndTime))
}

// HasBetaProgram applies the HasEdge predicate on the "beta_program" edge.
func HasBetaProgram() predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, BetaProgramTable, BetaProgramColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasBetaProgramWith applies the HasEdge predicate on the "beta_program" edge with a given conditions (other predicates).
func HasBetaProgramWith(preds ...predicate.BetaProgram) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(func(s *sql.Selector) {
		step := newBetaProgramStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasCustomer applies the HasEdge predicate on the "customer" edge.
func HasCustomer() predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, CustomerTable, CustomerColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCustomerWith applies the HasEdge predicate on the "customer" edge with a given conditions (other predicates).
func HasCustomerWith(preds ...predicate.Customer) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(func(s *sql.Selector) {
		step := newCustomerStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasClinic applies the HasEdge predicate on the "clinic" edge.
func HasClinic() predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ClinicTable, ClinicColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasClinicWith applies the HasEdge predicate on the "clinic" edge with a given conditions (other predicates).
func HasClinicWith(preds ...predicate.Clinic) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(func(s *sql.Selector) {
		step := newClinicStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.BetaProgramParticipation) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.BetaProgramParticipation) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.BetaProgramParticipation) predicate.BetaProgramParticipation {
	return predicate.BetaProgramParticipation(sql.NotPredicates(p))
}
