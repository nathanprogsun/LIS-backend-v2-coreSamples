package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Sample holds the schema definition for the Sample entity.
type SampleIDGenerate struct {
	ent.Schema
}

// Fields of the Sample.
func (SampleIDGenerate) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			StorageKey("sample_id").
			StructTag(`json:"sample_id"`).
			Unique().
			Immutable().
			Min(2000000),
		field.String("barcode").
			Optional().
			Unique(),
	}
}

// Edges of the Sample.
func (SampleIDGenerate) Edges() []ent.Edge {
	//TODO: add customer
	//TODO: add patient
	//TODO: add sample flags
	//TODO: sample types?
	return []ent.Edge{}
}

func (SampleIDGenerate) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sample_id_generate"},
	}
}