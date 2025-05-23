package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TestDetail holds the schema definition for the TestDetail entity.
type TestDetail struct {
	ent.Schema
}

//  test_detail_name   String
//  test_details_value String
//  test               test      @relation(fields: [test_id], references: [test_id])
//  isActive           Boolean   @default(true)
//  validate_until     DateTime?

// Fields of the TestDetail.
func (TestDetail) Fields() []ent.Field {
	return []ent.Field{
		field.Int("test_id"),
		field.String("test_detail_name"),
		field.String("test_details_value"),
		field.Bool("isActive").StorageKey("isActive").Default(true),
	}
}

// Edges of the TestDetail.
func (TestDetail) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("test", Test.Type).Field("test_id").Unique().Required(),
	}
}

func (TestDetail) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "test_detail"},
	}
}
