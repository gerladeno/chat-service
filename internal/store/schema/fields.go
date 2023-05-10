package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

func newCreatedAtField() ent.Field {
	return field.Time("created_at").
		Default(time.Now).
		Immutable()
}
