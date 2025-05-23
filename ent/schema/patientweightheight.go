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

// PatientWeightHeight holds the schema definition for the PatientWeightHeight entity.
type PatientWeightHeight struct {
	ent.Schema
}

// Fields of the PatientWeightHeight.
func (PatientWeightHeight) Fields() []ent.Field {
	return []ent.Field{
		field.Int("patient_id"),
		field.String("weight"),
		field.String("height"),
		field.String("weight_unit"),
		field.String("height_unit"),
		field.Time("created_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Default(time.Now),
		field.Time("updated_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Optional().UpdateDefault(time.Now),
	}
}

// Edges of the PatientWeightHeight.
func (PatientWeightHeight) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("patient", Patient.Type).
			Ref("patient_weight_height_history").Field("patient_id").Required().Unique(),
	}
}

func (PatientWeightHeight) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "patient_weight_height"},
	}
}
