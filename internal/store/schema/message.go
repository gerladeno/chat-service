package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/gerladeno/chat-service/internal/types"
)

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

const messageBodyMaxLength = 3000

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.MessageID{}).Default(types.NewMessageID).Unique(),
		field.UUID("author_id", types.UserID{}).Optional(),
		field.UUID("chat_id", types.ChatID{}),
		field.UUID("initial_request_id", types.RequestID{}).Immutable(),
		field.UUID("problem_id", types.ProblemID{}),
		field.Bool("is_visible_for_client").Default(false),
		field.Bool("is_visible_for_manager").Default(false),
		field.Text("body").NotEmpty().MaxLen(messageBodyMaxLength).Immutable(),
		field.Time("checked_at").Optional(),
		field.Bool("is_blocked").Default(false),
		field.Bool("is_service").Default(false).Immutable(),
		newCreatedAtField(),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("problem", Problem.Type).Ref("messages").Unique().Required().Field("problem_id"),
		edge.From("chat", Chat.Type).Ref("messages").Unique().Required().Field("chat_id"),
	}
}

func (Message) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("chat_id"),
		index.Fields("created_at", "is_visible_for_client"),
		index.Fields("initial_request_id").Annotations(entsql.IndexWhere("NOT is_service")).Unique(),
	}
}
