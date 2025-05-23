package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// BetaProgramParticipation holds the schema definition.
type BetaProgramParticipation struct {
	ent.Schema
}

// Fields of the BetaProgramParticipation.
func (BetaProgramParticipation) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Unique().
			Immutable(),
		field.Int("beta_program_id"),
		field.Int("customer_id"),
		field.Int("clinic_id"),
		field.Bool("is_active").
			Default(true),
		field.Bool("has_modified_start_time").
			Default(false),
		field.Time("modified_start_time").
			Nillable().
			Optional(),
		field.Time("modified_end_time").
			Nillable().
			Optional(),
	}
}

// Edges of the BetaProgramParticipation.
func (BetaProgramParticipation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("beta_program", BetaProgram.Type).
			Ref("participations").
			Field("beta_program_id").
			Unique().
			Required(),

		edge.From("customer", Customer.Type).
			Ref("customer_beta_program_participations").
			Field("customer_id").
			Unique().
			Required(),

		edge.From("clinic", Clinic.Type).
			Ref("clinic_beta_program_participations").
			Field("clinic_id").
			Unique().
			Required(),
	}
}

// Indexes to ensure uniqueness on (beta_program_id, customer_id, clinic_id)
func (BetaProgramParticipation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("beta_program_id", "customer_id", "clinic_id").
			Unique(),
	}
}
