package messagesrepo

import (
	"time"

	"github.com/gerladeno/chat-service/internal/store"
	"github.com/gerladeno/chat-service/internal/types"
)

type Message struct {
	ID                  types.MessageID
	RequestID           types.RequestID
	ChatID              types.ChatID
	AuthorID            types.UserID
	Body                string
	CreatedAt           time.Time
	IsVisibleForClient  bool
	IsVisibleForManager bool
	IsBlocked           bool
	IsService           bool
}

func adaptStoreMessage(m *store.Message) Message {
	return Message{
		ID:                  m.ID,
		RequestID:           m.InitialRequestID,
		ChatID:              m.ChatID,
		AuthorID:            m.AuthorID,
		Body:                m.Body,
		CreatedAt:           m.CreatedAt,
		IsVisibleForClient:  m.IsVisibleForClient,
		IsVisibleForManager: m.IsVisibleForManager,
		IsBlocked:           m.IsBlocked,
		IsService:           m.IsService,
	}
}
