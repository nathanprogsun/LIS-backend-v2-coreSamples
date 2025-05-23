package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

// SalesTerritory holds the schema definition for the SalesTerritory entity.
type SalesTerritory struct {
	ent.Schema
}

// Fields of the SalesTerritory.
func (SalesTerritory) Fields() []ent.Field {
	return []ent.Field{
		field.String("sales"),
		field.String("state").Optional(),
		field.Int("zipcode").Optional(),
		field.String("country").Optional(),
		field.Time("updatedAt").
			StorageKey("updated_at").StructTag(`json:"updated_at"`).
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Optional().UpdateDefault(time.Now),
	}
}

// Edges of the SalesTerritory.
func (SalesTerritory) Edges() []ent.Edge {
	return nil
}

func (SalesTerritory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sales_territory"},
	}
}
