package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

// LoginHistory holds the schema definition for the LoginHistory entity.
type LoginHistory struct {
	ent.Schema
}

// Fields of the LoginHistory.
func (LoginHistory) Fields() []ent.Field {
	return []ent.Field{
		field.String("username"),
		field.Time("login_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Default(time.Now),
		field.String("login_ip"),
		field.Bool("login_successfully"),
		field.String("failure_reason").Optional(),
		field.String("login_portal").Optional(),
		field.String("token").MaxLen(2000).Optional(),
	}
}

// Edges of the LoginHistory.
func (LoginHistory) Edges() []ent.Edge {
	return nil
}

func (LoginHistory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "login_history"},
	}
}
