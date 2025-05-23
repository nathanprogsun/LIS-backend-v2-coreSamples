package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Sample holds the schema definition for the Sample entity.
type Sample struct {
	ent.Schema
}

// Fields of the Sample.
func (Sample) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StorageKey("sample_id").StructTag(`json:"sample_id"`),
		field.String("accession_id"),
		field.String("sample_storage").Optional().Default("N/A"),
		field.Int("tube_count").Optional().Default(0),
		field.Int("order_id").Optional(),
		field.Int("patient_id").Optional(),
		field.String("sample_order_method").Optional(),
		field.Time("sample_collection_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Optional(),
		field.Time("sample_received_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Optional(),
		field.String("sample_description").Optional().Default("N/A"),
		field.Int("delayed_hours"),
		field.Time("sample_report_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Optional(),
		field.Time("internal_received_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Optional(),
		field.String("sample_report_type").Optional(),
		field.Int("customer_id").Optional(),
		field.Float("fasting_hours").Optional(),
		field.String("fasting_status").Optional(),
	}
}

// Edges of the Sample.
func (Sample) Edges() []ent.Edge {
	//TODO: add customer
	//TODO: add patient
	//TODO: add sample flags
	//TODO: sample types?
	return []ent.Edge{
		// For O2O relation in ent, table with the fk should use the from edge
		edge.From("order", OrderInfo.Type).Field("order_id").Unique().Ref("sample"),
		edge.From("tubes", Tube.Type).Ref("sample"),
		edge.From("sample_receive_records", TubeReceive.Type).Ref("sample"),
		edge.From("sample_required_tubes", TubeRequirement.Type).Ref("sample"),
		edge.From("patient", Patient.Type).Field("patient_id").Unique().Ref("samples"),
		edge.From("customer", Customer.Type).Field("customer_id").Ref("samples").Unique(),
	}
}

func (Sample) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sample"},
	}
}
