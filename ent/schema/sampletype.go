package schema

import (
	"coresamples/model"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// SampleType holds the schema definition for the SampleType entity.
type SampleType struct {
	ent.Schema
}

// Fields of the SampleType.
func (SampleType) Fields() []ent.Field {
	return append([]ent.Field{
		field.Int("id").StorageKey("sample_type_id").StructTag(`json:"sample_type_id"`),
		field.String("sample_type_name"),
		field.String("sample_type_code"),
		//TODO: fix the typo
		field.String("sample_type_enum").StorageKey("sample_type_emun").Optional(),
		field.String("sample_type_enum_old_lis_request").StorageKey("sample_type_emun_old_lis_request").Optional(),
		field.String("sample_type_description"),
		field.String("primary_sample_type_group"),
		field.Bool("is_active").StorageKey("isActive").Default(true)},
		model.CommonFields...)
}

// Edges of the SampleType.
func (SampleType) Edges() []ent.Edge {
	// link samples, test, tube_types, sample_type_group?
	//TODO: import the AB table
	return []ent.Edge{
		edge.To("tube_types", TubeType.Type).StorageKey(
			edge.Table("_sample_type_to_tube_type"),
		),
		edge.To("tests", Test.Type).StorageKey(
			edge.Table("_sample_type_to_test"),
		),
	}
}

func (SampleType) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sample_type"},
	}
}
