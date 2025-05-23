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

// LabOrderSendHistory holds the schema definition for the LabOrderSendHistory entity.
type LabOrderSendHistory struct {
	ent.Schema
}

// Fields of the LabOrderSendHistory.
func (LabOrderSendHistory) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StorageKey("lab_order_id").StructTag(`json:"lab_order_id"`),
		field.Int("sample_id"),
		field.String("tube_type"),
		field.Time("sendout_time").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Default(time.Now).Optional(),
		field.Bool("is_redraw_order").Optional().Default(false),
		field.Bool("is_lab_special_order").Optional().Default(false),
		field.String("action").Optional(),
		field.Bool("is_resend_blocked").Default(true).Optional(),
	}
}

// Edges of the LabOrderSendHistory.
func (LabOrderSendHistory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("sample", Sample.Type).Field("sample_id").Unique().Required(),
	}
}

func (LabOrderSendHistory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "lab_order_send_history"},
	}
}
