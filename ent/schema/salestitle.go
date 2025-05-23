package schema

import (
	"coresamples/model"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// SalesTitle holds the schema definition for the SalesTitle entity.
type SalesTitle struct {
	ent.Schema
}

// Fields of the SalesTitle.
func (SalesTitle) Fields() []ent.Field {
	return append([]ent.Field{
		field.String("title_name").Unique(),
		field.Int("order").Unique(),
	}, model.CommonFields...)
}

// Edges of the SalesTitle.
func (SalesTitle) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("sales", SalesTeam.Type).Ref("title"),
	}
}

func (SalesTitle) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sales_title"},
	}
}
