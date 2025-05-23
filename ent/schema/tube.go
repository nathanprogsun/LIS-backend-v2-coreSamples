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

// Tube holds the schema definition for the Tube entity.
type Tube struct {
	ent.Schema
}

// Fields of the Tube.
func (Tube) Fields() []ent.Field {
	return []ent.Field{
		field.String("tube_id").Unique(),
		field.Int("sample_id"),
		field.String("tube_storage").Default("N/A"),
		field.Time("tube_receive_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Default(time.Now),
		field.Int("tube_turnaround_time"),
		field.Int("tube_stability").Default(0),
		field.Bool("isActive").StorageKey("isActive").Default(true),
		field.String("issues").Default("N/A"),
		field.Time("tube_collection_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Default(time.Now),
	}
}

// Edges of the Tube.
func (Tube) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tube_type", TubeType.Type),
		edge.To("sample", Sample.Type).Field("sample_id").Unique().Required(),
	}
}

func (Tube) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "tube"},
	}
}
