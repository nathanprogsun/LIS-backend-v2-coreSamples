package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// PendingOrderCredits holds the schema definition for the PendingOrderCredits entity.
type PendingOrderCredits struct {
	ent.Schema
}

// Fields of the PendingOrderCredits.
func (PendingOrderCredits) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("order_id"),
		field.Int64("credit_id"),
		field.Int64("clinic_id"),
	}
}

// Edges of the PendingOrderCredits.
func (PendingOrderCredits) Edges() []ent.Edge {
	return nil
}

func (PendingOrderCredits) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "pending_order_credits"},
	}
}
