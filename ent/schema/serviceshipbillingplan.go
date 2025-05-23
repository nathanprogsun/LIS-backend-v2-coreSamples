package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ServiceshipBillingPlan holds the schema definition for the ServiceshipBillingPlan entity.
type ServiceshipBillingPlan struct {
	ent.Schema
}

// Fields of the ServiceshipBillingPlan.
func (ServiceshipBillingPlan) Fields() []ent.Field {
	return []ent.Field{
		field.Float32("fee").
			SchemaType(map[string]string{
				dialect.MySQL: "decimal(7,2)",
			}),
		field.Int32("billing_cycle").Positive(),
		field.Enum("interval").Values("monthly", "daily").Default("monthly"),
		field.Time("effective_time"),
	}
}

// Edges of the ServiceshipBillingPlan.
func (ServiceshipBillingPlan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("account_subscription", AccountSubscription.Type),
		edge.From("serviceship", Serviceship.Type).
			Ref("serviceship_billing_plan").
			Unique().
			Required(), // a billing plan cannot be created without a membership
	}
}

func (ServiceshipBillingPlan) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "serviceship_billing_plan"},
	}
}
