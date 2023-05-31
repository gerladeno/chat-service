package getchathistory

import (
	"errors"
	"fmt"
	"time"

	"github.com/gerladeno/chat-service/internal/types"
	"github.com/gerladeno/chat-service/internal/validator"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrInvalidCursor  = errors.New("invalid cursor")
)

type Request struct {
	ChatID    types.ChatID `validate:"required"`
	ManagerID types.UserID `validate:"required"`
	PageSize  int          `validate:"omitempty,gte=10,lte=100"`
	Cursor    string       `validate:"omitempty,base64url"`
}

func (r Request) Validate() error {
	if err := r.ChatID.Validate(); err != nil {
		return fmt.Errorf("validate get chat history request chat id: %v", err)
	}
	if err := r.ManagerID.Validate(); err != nil {
		return fmt.Errorf("validate get chats history request managerId: %v", err)
	}
	if r.PageSize == 0 && r.Cursor == "" || r.PageSize != 0 && r.Cursor != "" {
		return ErrInvalidRequest
	}
	return validator.Validator.Struct(r)
}

type Response struct {
	Messages   []Message
	NextCursor string
}

type Message struct {
	ID        types.MessageID
	AuthorID  types.UserID
	Body      string
	CreatedAt time.Time
}
