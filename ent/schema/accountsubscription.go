package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// AccountSubscription holds the schema definition for the AccountSubscription entity.
type AccountSubscription struct {
	ent.Schema
}

// Fields of the AccountSubscription.
func (AccountSubscription) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("account_id").NonNegative(),
		field.Enum("account_type").Values("clinic"),
		field.String("subscriber_name"),
		field.String("email"),
		field.Time("start_time"),
		field.Time("end_time").Default(time.Time{}),
	}
}

// Edges of the AccountSubscription.
func (AccountSubscription) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("serviceship_billing_plan", ServiceshipBillingPlan.Type).
			Ref("account_subscription").
			Unique().
			Required(),
		edge.From("serviceship", Serviceship.Type).
			Ref("account_subscription").
			Unique().
			Required(),
	}
}

func (AccountSubscription) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "account_subscription"},
	}
}
