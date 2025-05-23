package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Zipcode holds the schema definition for the Zipcode entity.
type Zipcode struct {
	ent.Schema
}

// Fields of the Zipcode.
func (Zipcode) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StorageKey("zipcode").StructTag(`json:"zipcode"`),
		field.String("ZipCodeType").StorageKey("ZipCodeType"),
		field.String("City").StorageKey("City"),
		field.String("State").StorageKey("State"),
		field.String("LocationType").StorageKey("LocationType"),
		field.Float("Lat").
			SchemaType(map[string]string{
				dialect.MySQL: "decimal(8,2)",
			}).StorageKey("Lat").Optional(),
		field.Float("Long").SchemaType(map[string]string{
			dialect.MySQL: "decimal(8,2)",
		}).StorageKey("Long").Optional(),
	}
}

// Edges of the Zipcode.
func (Zipcode) Edges() []ent.Edge {
	return nil
}

func (Zipcode) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "zipcode"},
	}
}
