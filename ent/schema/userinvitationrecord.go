package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// UserInvitationRecord holds the schema definition for the UserInvitationRecord entity.
type UserInvitationRecord struct {
	ent.Schema
}

// Fields of the UserInvitationRecord.
func (UserInvitationRecord) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Unique().
			Immutable().
			Positive(),
		field.Int("customer_id"),
		field.String("invitation_link").
			MaxLen(3000),
	}
}
