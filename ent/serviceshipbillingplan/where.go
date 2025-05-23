// Code generated by ent, DO NOT EDIT.

package serviceshipbillingplan

import (
	"coresamples/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldLTE(FieldID, id))
}

// Fee applies equality check predicate on the "fee" field. It's identical to FeeEQ.
func Fee(v float32) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldEQ(FieldFee, v))
}

// BillingCycle applies equality check predicate on the "billing_cycle" field. It's identical to BillingCycleEQ.
func BillingCycle(v int32) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldEQ(FieldBillingCycle, v))
}

// EffectiveTime applies equality check predicate on the "effective_time" field. It's identical to EffectiveTimeEQ.
func EffectiveTime(v time.Time) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldEQ(FieldEffectiveTime, v))
}

// FeeEQ applies the EQ predicate on the "fee" field.
func FeeEQ(v float32) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldEQ(FieldFee, v))
}

// FeeNEQ applies the NEQ predicate on the "fee" field.
func FeeNEQ(v float32) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldNEQ(FieldFee, v))
}

// FeeIn applies the In predicate on the "fee" field.
func FeeIn(vs ...float32) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldIn(FieldFee, vs...))
}

// FeeNotIn applies the NotIn predicate on the "fee" field.
func FeeNotIn(vs ...float32) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldNotIn(FieldFee, vs...))
}

// FeeGT applies the GT predicate on the "fee" field.
func FeeGT(v float32) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldGT(FieldFee, v))
}

// FeeGTE applies the GTE predicate on the "fee" field.
func FeeGTE(v float32) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldGTE(FieldFee, v))
}

// FeeLT applies the LT predicate on the "fee" field.
func FeeLT(v float32) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldLT(FieldFee, v))
}

// FeeLTE applies the LTE predicate on the "fee" field.
func FeeLTE(v float32) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldLTE(FieldFee, v))
}

// BillingCycleEQ applies the EQ predicate on the "billing_cycle" field.
func BillingCycleEQ(v int32) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldEQ(FieldBillingCycle, v))
}

// BillingCycleNEQ applies the NEQ predicate on the "billing_cycle" field.
func BillingCycleNEQ(v int32) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldNEQ(FieldBillingCycle, v))
}

// BillingCycleIn applies the In predicate on the "billing_cycle" field.
func BillingCycleIn(vs ...int32) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldIn(FieldBillingCycle, vs...))
}

// BillingCycleNotIn applies the NotIn predicate on the "billing_cycle" field.
func BillingCycleNotIn(vs ...int32) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldNotIn(FieldBillingCycle, vs...))
}

// BillingCycleGT applies the GT predicate on the "billing_cycle" field.
func BillingCycleGT(v int32) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldGT(FieldBillingCycle, v))
}

// BillingCycleGTE applies the GTE predicate on the "billing_cycle" field.
func BillingCycleGTE(v int32) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldGTE(FieldBillingCycle, v))
}

// BillingCycleLT applies the LT predicate on the "billing_cycle" field.
func BillingCycleLT(v int32) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldLT(FieldBillingCycle, v))
}

// BillingCycleLTE applies the LTE predicate on the "billing_cycle" field.
func BillingCycleLTE(v int32) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldLTE(FieldBillingCycle, v))
}

// IntervalEQ applies the EQ predicate on the "interval" field.
func IntervalEQ(v Interval) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldEQ(FieldInterval, v))
}

// IntervalNEQ applies the NEQ predicate on the "interval" field.
func IntervalNEQ(v Interval) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldNEQ(FieldInterval, v))
}

// IntervalIn applies the In predicate on the "interval" field.
func IntervalIn(vs ...Interval) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldIn(FieldInterval, vs...))
}

// IntervalNotIn applies the NotIn predicate on the "interval" field.
func IntervalNotIn(vs ...Interval) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldNotIn(FieldInterval, vs...))
}

// EffectiveTimeEQ applies the EQ predicate on the "effective_time" field.
func EffectiveTimeEQ(v time.Time) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldEQ(FieldEffectiveTime, v))
}

// EffectiveTimeNEQ applies the NEQ predicate on the "effective_time" field.
func EffectiveTimeNEQ(v time.Time) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldNEQ(FieldEffectiveTime, v))
}

// EffectiveTimeIn applies the In predicate on the "effective_time" field.
func EffectiveTimeIn(vs ...time.Time) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldIn(FieldEffectiveTime, vs...))
}

// EffectiveTimeNotIn applies the NotIn predicate on the "effective_time" field.
func EffectiveTimeNotIn(vs ...time.Time) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldNotIn(FieldEffectiveTime, vs...))
}

// EffectiveTimeGT applies the GT predicate on the "effective_time" field.
func EffectiveTimeGT(v time.Time) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldGT(FieldEffectiveTime, v))
}

// EffectiveTimeGTE applies the GTE predicate on the "effective_time" field.
func EffectiveTimeGTE(v time.Time) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldGTE(FieldEffectiveTime, v))
}

// EffectiveTimeLT applies the LT predicate on the "effective_time" field.
func EffectiveTimeLT(v time.Time) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldLT(FieldEffectiveTime, v))
}

// EffectiveTimeLTE applies the LTE predicate on the "effective_time" field.
func EffectiveTimeLTE(v time.Time) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.FieldLTE(FieldEffectiveTime, v))
}

// HasAccountSubscription applies the HasEdge predicate on the "account_subscription" edge.
func HasAccountSubscription() predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, AccountSubscriptionTable, AccountSubscriptionColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasAccountSubscriptionWith applies the HasEdge predicate on the "account_subscription" edge with a given conditions (other predicates).
func HasAccountSubscriptionWith(preds ...predicate.AccountSubscription) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(func(s *sql.Selector) {
		step := newAccountSubscriptionStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasServiceship applies the HasEdge predicate on the "serviceship" edge.
func HasServiceship() predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ServiceshipTable, ServiceshipColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasServiceshipWith applies the HasEdge predicate on the "serviceship" edge with a given conditions (other predicates).
func HasServiceshipWith(preds ...predicate.Serviceship) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(func(s *sql.Selector) {
		step := newServiceshipStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ServiceshipBillingPlan) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ServiceshipBillingPlan) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.ServiceshipBillingPlan) predicate.ServiceshipBillingPlan {
	return predicate.ServiceshipBillingPlan(sql.NotPredicates(p))
}
