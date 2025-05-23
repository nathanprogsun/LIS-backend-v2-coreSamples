package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// TubeRequirement holds the schema definition for the TubeRequirement entity.
type TubeRequirement struct {
	ent.Schema
}

// Fields of the TubeRequirement.
func (TubeRequirement) Fields() []ent.Field {
	return []ent.Field{
		field.Int("sample_id").Optional(),
		field.String("tube_type"),
		field.Int("required_count"),
		field.Time("required_count_create_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Default(time.Now),
		field.String("required_by"),
		field.String("modified_by").Optional(),
		field.Time("modified_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the TubeRequirement.
func (TubeRequirement) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("sample", Sample.Type).Field("sample_id").Unique(),
	}
}

func (TubeRequirement) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "tube_requirement"},
	}
}
