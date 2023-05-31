package closechat

import (
	"github.com/gerladeno/chat-service/internal/types"
	"github.com/gerladeno/chat-service/internal/validator"
)

type Request struct {
	ID        types.RequestID `validate:"required"`
	ChatID    types.ChatID    `validate:"required"`
	ManagerID types.UserID    `validate:"required"`
}

func (r Request) Validate() error {
	return validator.Validator.Struct(r)
}
