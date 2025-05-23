package schema

import (
	"coresamples/model"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// SalesTeam holds the schema definition for the SalesTeam entity.
type SalesTeam struct {
	ent.Schema
}

// Fields of the SalesTeam.
func (SalesTeam) Fields() []ent.Field {
	return append([]ent.Field{
		field.Int("internal_user_id"),
		field.Int("supervisor_id").Optional(),
		field.Int("title_id").Optional(),
	}, model.CommonFields...)
}

// Edges of the SalesTeam.
func (SalesTeam) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("supervisor", SalesTeam.Type).
			Field("supervisor_id").
			Unique().
			From("subordinates"),
		edge.To("internal_user", InternalUser.Type).
			Field("internal_user_id").
			Unique().
			Required(),
		edge.To("title", SalesTitle.Type).
			Field("title_id").
			Unique(),
	}
}

func (SalesTeam) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sales_team"},
	}
}
