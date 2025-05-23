package model

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"time"
)

var CommonFields = []ent.Field{
	field.Time("created_time").
		SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}).Default(time.Now),
	field.Time("updated_time").
		SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}).Optional().UpdateDefault(time.Now),
}
