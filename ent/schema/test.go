package schema

import (
	"coresamples/model"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Test holds the schema definition for the Test entity.
type Test struct {
	ent.Schema
}

// Fields of the Test.
func (Test) Fields() []ent.Field {
	return append([]ent.Field{
		field.Int("id").StorageKey("test_id").StructTag(`json:"test_id"`),
		field.String("test_name"),
		field.String("test_code"),
		field.String("display_name"),
		field.String("test_description").Nillable(),
		field.String("assay_name").Nillable(),
		field.Bool("isActive").StorageKey("isActive"),
	}, model.CommonFields...)
}

// Edges of the Test.
func (Test) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("test_details", TestDetail.Type).Ref("test").StructTag(`json:"test_details"`),
		edge.From("order_info", OrderInfo.Type).Ref("tests"),
		edge.From("sample_types", SampleType.Type).Ref("tests"),
		//TODO: import the AB table
		edge.From("tube_types", TubeType.Type).Ref("tests"),
	}
}

func (Test) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "test"},
	}
}
