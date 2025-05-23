package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// BetaProgram holds the schema definition for the BetaProgram entity.
type BetaProgram struct {
	ent.Schema
}

// Fields of the BetaProgram.
func (BetaProgram) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StorageKey("beta_program_id").StructTag(`json:"beta_program_id"`),
		field.String("beta_program_name"),
		field.String("beta_program_description"),
		field.Bool("is_active").
			Default(true),
		field.Time("beta_program_start_time").
			Default(time.Now),
		field.Time("beta_program_end_time").
			Nillable().
			Optional(),
		field.Time("updated_time").
			UpdateDefault(time.Now),
		field.Time("beta_program_added_on").
			Nillable().
			Optional().
			Default(time.Now),
		field.Bool("allow_self_signup").
			Nillable().
			Optional().
			Default(true),
	}
}

// Edges of the BetaProgram.
func (BetaProgram) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("participations", BetaProgramParticipation.Type),
	}
}
