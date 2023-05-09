package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/gerladeno/chat-service/internal/types"
	"time"
)

// jobMaxAttempts is some limit as protection from endless retries of outbox jobs.
const jobMaxAttempts = 30

type Job struct {
	ent.Schema
}

func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.JobID{}).Default(types.NewJobID).Unique().Immutable(),
		field.Text("name").Immutable(),
		field.Text("payload").Immutable(),
		field.Int("attempts").Max(jobMaxAttempts).Default(0),
		field.Time("available_at").Immutable(),
		field.Time("reserved_until").Default(time.Now()),
		newCreatedAtField(),
	}
}

func (Job) Indexes() []ent.Index {
	// FIXME: Расставь индексы на основе запросов в сервисе Outbox.
	return nil
}

type FailedJob struct {
	ent.Schema
}

func (FailedJob) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.FailedJobID{}).Default(types.NewFailedJobID).Unique().Immutable(),
		field.Text("name").Immutable(),
		field.Text("payload").Immutable(),
		field.Text("reason").Immutable(),
		newCreatedAtField(),
	}
}
