package schema

import (
	"coresamples/model"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"errors"
	"unicode/utf8"
)

// TestList holds the schema definition for the TestList entity.
type TestList struct {
	ent.Schema
}

// Fields of the TestList.
func (TestList) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StorageKey("test_id"),
		field.Enum("test_instrument").NamedValues(model.GetTestInstrumentNamedValues()...).
			Annotations(
				entsql.Annotation{
					Charset:   "utf8",
					Collation: "utf8_general_ci",
				},
			).Optional().Default("N/A"),
		field.Enum("tube_type").NamedValues(model.GetTubeTypeNamedValues()...).Annotations(
			entsql.Annotation{
				Charset:   "utf8",
				Collation: "utf8_general_ci",
			},
		),
		field.String("DI_group_name").
			StorageKey("DI_group_name").
			Annotations(entsql.Annotation{
				Size: 45,
			}, entsql.Annotation{
				Charset:   "latin1",
				Collation: "latin1_swedish_ci",
			}).
			Validate(MaxRuneCount(45)).
			Optional(),
		field.Float("volume_required").
			SchemaType(map[string]string{
				dialect.MySQL: "decimal(7,3)",
			}),
		field.Bool("blood_type"),
	}
}

func (TestList) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "test_list"},
	}
}

// Edges of the TestList.
func (TestList) Edges() []ent.Edge {
	return nil
}

// MaxRuneCount validates the rune length of a string by using the unicode/utf8 package.
func MaxRuneCount(maxLen int) func(s string) error {
	return func(s string) error {
		if utf8.RuneCountInString(s) > maxLen {
			return errors.New("value is more than the max length")
		}
		return nil
	}
}
