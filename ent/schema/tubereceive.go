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

// TubeReceive holds the schema definition for the TubeReceive entity.
type TubeReceive struct {
	ent.Schema
}

// Fields of the TubeReceive.
func (TubeReceive) Fields() []ent.Field {
	return []ent.Field{
		field.Int("sample_id"),
		field.String("tube_type"),
		field.Int("received_count"),
		field.Time("received_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Default(time.Now),
		field.String("received_by"),
		field.String("modified_by").Optional(),
		field.Time("modified_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Optional().UpdateDefault(time.Now),
		field.Time("collection_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Default(time.Now),
		field.Bool("is_redraw").Default(false),
		field.Bool("is_rerun").Default(false),
	}
}

// Edges of the TubeReceive.
func (TubeReceive) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("sample", Sample.Type).Field("sample_id").Unique().Required(),
	}
}

func (TubeReceive) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "tube_receive"},
	}
}
