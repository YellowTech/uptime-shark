package schema

import (
	"uptime/api/logentry"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Monitor struct {
	ent.Schema
}

// Fields of the User.
func (Monitor) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
            Default(uuid.New).Unique(),
		field.String("name").MaxLen(100).NotEmpty(),

		// check intervall in seconds
		field.Int64("interval").Positive(),
		field.Int64("nextCheck").Default(0),

		field.Bool("status"),
		field.String("statusMessage"),
		field.Bool("inverted").Default(false),
		
		field.JSON("logs", []logentry.LogEntry{}),
		field.Int("nrLogs").Default(30),

		// mode of the monitor
		// is enforced when trying to save a service
		field.String("mode").NotEmpty(),

		field.String("url").MaxLen(1024),
		field.Int("retries"),

		
	}
}

// Edges of the User.
func (Monitor) Edges() []ent.Edge {
	return nil
}
