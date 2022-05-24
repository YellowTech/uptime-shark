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

func (Monitor) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
            Default(uuid.New).Unique(),
		field.String("name").MaxLen(100).NotEmpty(),

		// check intervall in seconds
		field.Int64("interval").Positive(),
		field.Int64("nextCheck").Default(0),
		
		// wether the monitor is enabled or not
		field.Bool("enabled"),
		
		// is the monitor up or down
		field.Bool("status").Default(true),
		field.String("statusMessage"),

		// wether reachable is bad
		field.Bool("inverted").Default(false),
		
		field.JSON("logs", []logentry.LogEntry{}),
		field.Int("nrLogs").Default(30),

		// mode of the monitor
		// is enforced when trying to save a service TM
		field.String("mode").NotEmpty(),

		field.String("url").MaxLen(1024),
		field.Int("retries"),
	}
}

func (Monitor) Edges() []ent.Edge {
	return nil
}
