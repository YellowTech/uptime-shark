package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Notification struct {
	ent.Schema
}

func (Notification) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
            Default(uuid.New).Unique(),
		field.String("name").MaxLen(100).NotEmpty(),

		field.JSON("settings", []string{}),

		field.Bool("active"),
	}
}

func (Notification) Edges() []ent.Edge {
	return nil
}
