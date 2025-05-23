package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TubeType holds the schema definition for the TubeType entity.
type TubeType struct {
	ent.Schema
}

// Fields of the TubeType.
func (TubeType) Fields() []ent.Field {
	return []ent.Field{
		field.String("tube_name"),
		field.String("tube_type_enum"),
		field.String("tube_type_symbol"),
	}
}

// Edges of the TubeType.
func (TubeType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tube", Tube.Type).Ref("tube_type"),
		edge.From("sample_types", SampleType.Type).Ref("tube_types"),
		edge.To("tests", Test.Type).StorageKey(
			edge.Table("_tube_type_to_test"),
			edge.Columns("tube_type_id", "test_id"),
		),
	}
}

func (TubeType) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "tube_type"},
	}
}
