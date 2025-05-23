// Code generated by ent, DO NOT EDIT.

package serviceship

import (
	"coresamples/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldLTE(FieldID, id))
}

// Tag applies equality check predicate on the "tag" field. It's identical to TagEQ.
func Tag(v string) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldEQ(FieldTag, v))
}

// TagEQ applies the EQ predicate on the "tag" field.
func TagEQ(v string) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldEQ(FieldTag, v))
}

// TagNEQ applies the NEQ predicate on the "tag" field.
func TagNEQ(v string) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldNEQ(FieldTag, v))
}

// TagIn applies the In predicate on the "tag" field.
func TagIn(vs ...string) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldIn(FieldTag, vs...))
}

// TagNotIn applies the NotIn predicate on the "tag" field.
func TagNotIn(vs ...string) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldNotIn(FieldTag, vs...))
}

// TagGT applies the GT predicate on the "tag" field.
func TagGT(v string) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldGT(FieldTag, v))
}

// TagGTE applies the GTE predicate on the "tag" field.
func TagGTE(v string) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldGTE(FieldTag, v))
}

// TagLT applies the LT predicate on the "tag" field.
func TagLT(v string) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldLT(FieldTag, v))
}

// TagLTE applies the LTE predicate on the "tag" field.
func TagLTE(v string) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldLTE(FieldTag, v))
}

// TagContains applies the Contains predicate on the "tag" field.
func TagContains(v string) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldContains(FieldTag, v))
}

// TagHasPrefix applies the HasPrefix predicate on the "tag" field.
func TagHasPrefix(v string) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldHasPrefix(FieldTag, v))
}

// TagHasSuffix applies the HasSuffix predicate on the "tag" field.
func TagHasSuffix(v string) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldHasSuffix(FieldTag, v))
}

// TagEqualFold applies the EqualFold predicate on the "tag" field.
func TagEqualFold(v string) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldEqualFold(FieldTag, v))
}

// TagContainsFold applies the ContainsFold predicate on the "tag" field.
func TagContainsFold(v string) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldContainsFold(FieldTag, v))
}

// TypeEQ applies the EQ predicate on the "type" field.
func TypeEQ(v Type) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldEQ(FieldType, v))
}

// TypeNEQ applies the NEQ predicate on the "type" field.
func TypeNEQ(v Type) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldNEQ(FieldType, v))
}

// TypeIn applies the In predicate on the "type" field.
func TypeIn(vs ...Type) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldIn(FieldType, vs...))
}

// TypeNotIn applies the NotIn predicate on the "type" field.
func TypeNotIn(vs ...Type) predicate.Serviceship {
	return predicate.Serviceship(sql.FieldNotIn(FieldType, vs...))
}

// HasServiceshipBillingPlan applies the HasEdge predicate on the "serviceship_billing_plan" edge.
func HasServiceshipBillingPlan() predicate.Serviceship {
	return predicate.Serviceship(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ServiceshipBillingPlanTable, ServiceshipBillingPlanColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasServiceshipBillingPlanWith applies the HasEdge predicate on the "serviceship_billing_plan" edge with a given conditions (other predicates).
func HasServiceshipBillingPlanWith(preds ...predicate.ServiceshipBillingPlan) predicate.Serviceship {
	return predicate.Serviceship(func(s *sql.Selector) {
		step := newServiceshipBillingPlanStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasAccountSubscription applies the HasEdge predicate on the "account_subscription" edge.
func HasAccountSubscription() predicate.Serviceship {
	return predicate.Serviceship(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, AccountSubscriptionTable, AccountSubscriptionColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasAccountSubscriptionWith applies the HasEdge predicate on the "account_subscription" edge with a given conditions (other predicates).
func HasAccountSubscriptionWith(preds ...predicate.AccountSubscription) predicate.Serviceship {
	return predicate.Serviceship(func(s *sql.Selector) {
		step := newAccountSubscriptionStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Serviceship) predicate.Serviceship {
	return predicate.Serviceship(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Serviceship) predicate.Serviceship {
	return predicate.Serviceship(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Serviceship) predicate.Serviceship {
	return predicate.Serviceship(sql.NotPredicates(p))
}
