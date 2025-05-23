package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// TubeInstructions holds the schema definition for the TubeInstructions entity.
type TubeInstructions struct {
	ent.Schema
}

// Fields of the TubeInstructions.
func (TubeInstructions) Fields() []ent.Field {
	return []ent.Field{
		field.String("tube_name_enum"),
		field.Int("sort_order"),
		field.String("tube_instructions").MaxLen(1024),
		field.String("tube_name"),
		field.String("shipping_box"),
		field.String("transfer_tubes_to_send"),
		field.Bool("blood_type").Default(false),
	}
}

// Edges of the TubeInstructions.
func (TubeInstructions) Edges() []ent.Edge {
	return nil
}

func (TubeInstructions) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "tube_instructions"},
	}
}
