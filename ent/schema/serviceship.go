package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Serviceship holds the schema definition for the Serviceship entity.
type Serviceship struct {
	ent.Schema
}

// Fields of the Serviceship.
func (Serviceship) Fields() []ent.Field {
	return []ent.Field{
		field.String("tag").NotEmpty(),
		field.Enum("type").Values("membership"),
	}
}

// Edges of the Serviceship.
func (Serviceship) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("serviceship_billing_plan", ServiceshipBillingPlan.Type),
		edge.To("account_subscription", AccountSubscription.Type),
	}
}

func (Serviceship) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "serviceship"},
	}
}
