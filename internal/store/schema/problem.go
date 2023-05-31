package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/gerladeno/chat-service/internal/types"
)

// Problem holds the schema definition for the Problem entity.
type Problem struct {
	ent.Schema
}

// Fields of the Problem.
func (Problem) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.ProblemID{}).Default(types.NewProblemID).Unique().Immutable(),
		field.UUID("chat_id", types.ChatID{}),
		field.UUID("manager_id", types.UserID{}).Optional(),
		field.Time("resolved_at").Optional(),
		field.UUID("resolved_request_id", types.RequestID{}).Optional(),
		newCreatedAtField(),
	}
}

// Edges of the Problem.
func (Problem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("messages", Message.Type),
		edge.From("chat", Chat.Type).Ref("problems").Unique().Required().Field("chat_id"),
	}
}

func (Problem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id", "manager_id").Annotations(entsql.IndexWhere("resolved_at IS NULL")).Unique(),
		index.Fields("chat_id"),
	}
}
