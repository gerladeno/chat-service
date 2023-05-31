package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
	return []ent.Index{
		index.Fields("available_at", "reserved_until"),
	}
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
