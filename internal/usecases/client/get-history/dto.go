package gethistory

import (
	"time"

	"github.com/gerladeno/chat-service/internal/types"
	"github.com/gerladeno/chat-service/internal/validator"
)

type Request struct {
	ID       types.RequestID `validate:"required"`
	ClientID types.UserID    `validate:"required"`
	PageSize int             `validate:"omitempty,gte=10,lte=100"`
	Cursor   string          `validate:"omitempty,base64url"`
}

func (r Request) Validate() error {
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
	ID         types.MessageID
	AuthorID   types.UserID
	Body       string
	CreatedAt  time.Time
	IsReceived bool
	IsBlocked  bool
	IsService  bool
}
