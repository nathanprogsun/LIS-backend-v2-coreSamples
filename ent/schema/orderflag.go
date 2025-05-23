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

// OrderFlag holds the schema definition for the OrderFlag entity.
type OrderFlag struct {
	ent.Schema
}

// Fields of the OrderFlag.
func (OrderFlag) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StorageKey("order_flag_id").StructTag(`json:"order_flag_id"`),
		field.String("order_flag_name").Unique(),
		field.String("order_flag_description").Optional(),
		field.String("order_flag_display_name").Optional(),
		field.Bool("order_flag_is_active").Default(true),
		field.Time("order_flag_created_at").
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}).Default(time.Now),
		field.String("order_flag_color").Optional().Default("#90EE90"),
		field.String("order_flaged_by").Optional().Default("System"),
		field.Bool("order_flag_allow_duplicates_under_same_category").Default(false),
		field.String("order_flag_category").Optional(),
		field.Int("order_flag_level").Default(0),
	}
}

// Edges of the OrderFlag.
func (OrderFlag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("flagged_orders", OrderInfo.Type).StorageKey(
			edge.Table("_order_flag_to_order"),
			edge.Columns("order_flag_id", "order_id"),
			//edge.Symbols("order_flag_id", "order_id"), symbols creates a foreign key of given name
		),
	}
}

func (OrderFlag) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "order_flag"},
	}
}
