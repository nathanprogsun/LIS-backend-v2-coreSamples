package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// OrderInfo holds the schema definition for the OrderInfo entity.
type OrderInfo struct {
	ent.Schema
}

// Fields of the OrderInfo.
func (OrderInfo) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StorageKey("order_id").StructTag(`json:"order_id"`),
		//TODO: requires some special handling when migrating data since we change from M2M to O2M
		//field.Int("patient_id").Optional(),
		field.String("order_title").Default("Vibrant America Order"),
		field.String("order_type").Default("Initial"),
		field.String("order_description").Default("Normal Order"),
		field.String("order_confirmation_number"),
		field.Int("clinic_id").Optional(),
		field.Int("customer_id").Optional(),
		field.Time("order_create_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Default(time.Now),
		field.Time("order_service_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Default(time.Now).Optional(),
		field.Time("order_process_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Default(time.Now).Optional(),
		field.Time("order_redraw_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Default(time.Now).Optional(),
		field.Time("order_cancel_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Default(time.Now).Optional(),
		field.Bool("isActive").StorageKey("isActive").Default(true),
		field.Bool("has_order_setting").Default(false),
		field.Bool("order_canceled").Default(false),
		field.Bool("order_flagged").Default(false),
		field.String("order_status").Default("order_received").Optional(),
		field.String("order_major_status").Default("scheduled_order").Optional(),
		field.String("order_kit_status").Default("kit_ready_for_shipment").Optional(),
		field.String("order_report_status").Default("report_not_ready").Optional(),
		field.String("order_tnp_issue_status").Default("no_tnp_issue").Optional(),
		field.String("order_billing_issue_status").Default("no_billing_issue").Optional(),
		field.String("order_missing_info_issue_status").Default("no_missing_info_issue").Optional(),
		field.String("order_incomplete_questionnaire_issue_status").Default("no_incomplete_questionnaire_issue").Optional(),
		field.String("order_ny_waive_form_issue_status").Default("no_new_ny_waive_form_issue").Optional(),
		field.String("order_lab_issue_status").Default("no_lab_issue").Optional(),
		field.Time("order_processing_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Optional(),
		field.String("order_minor_status").Optional(),
		field.String("patient_first_name").Optional(),
		field.String("patient_last_name").Optional(),
		field.String("order_source").Default("BillingSpring").Optional(),
		field.String("order_charge_method").Optional(),
		field.String("order_placing_type").Default("Standard").Optional(),
		field.String("billing_order_id").Optional(),
		//TODO: add indexes and edges
		field.Int("contact_id").Optional(),
		field.Int("address_id").Optional(),
	}
}

// Edges of the OrderInfo.
func (OrderInfo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tests", Test.Type).StorageKey(
			edge.Table("_order_info_to_test"),
			edge.Columns("order_id", "test_id"),
			//edge.Symbols("order_id", "test_id"),
		),
		edge.From("order_flags", OrderFlag.Type).Ref("flagged_orders"),
		//TODO: change sample order relation to O2O
		edge.To("sample", Sample.Type).Unique(),
		edge.From("contact", Contact.Type).Ref("orders").Field("contact_id").Unique(),
		edge.From("address", Address.Type).Ref("orders").Field("address_id").Unique(),
		edge.From("clinic", Clinic.Type).Ref("clinic_orders").Field("clinic_id").Unique(),
		edge.From("customer_info", Customer.Type).Ref("orders").Field("customer_id").Unique(),
		//edge.From("patient", Patient.Type).Ref("patient_orders").Unique().Field("patient_id"),
	}
}

func (OrderInfo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "order_info"},
	}
}
