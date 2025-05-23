package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// CountryList holds the schema definition for the CountryList entity.
type CountryList struct {
	ent.Schema
}

// Fields of the CountryList.
func (CountryList) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StorageKey("vibrant_country_id").StructTag(`json:"vibrant_country_id"`),
		field.String("country_name"),
		field.String("alpha_2_code"),
		field.String("alpha_3_code"),
		field.String("country_code_enum"),
		field.String("iso"),
		field.String("country_region"),
		field.String("country_subregion"),
		field.String("country_region_code"),
		field.String("country_sub_region_code"),
	}
}

// Edges of the CountryList.
func (CountryList) Edges() []ent.Edge {
	return nil
}

func (CountryList) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "country_list"},
	}
}
