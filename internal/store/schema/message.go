package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/gerladeno/chat-service/internal/types"
	"time"
)

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.MessageID{}).Default(types.NewMessageID).Unique(),
		field.UUID("author_id", types.UserID{}),
		field.Bool("is_visible_for_client"),
		field.Bool("is_visible_for_manager"),
		field.Text("body"),
		field.Time("checked_at").Optional(),
		field.Bool("is_blocked"),
		field.Bool("is_service"),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("problem", Problem.Type).Ref("messages").Unique(),
		edge.From("chat", Chat.Type).Ref("messages").Unique(),
	}
}
