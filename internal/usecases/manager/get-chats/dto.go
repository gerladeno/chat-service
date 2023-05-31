package getchats

import (
	"fmt"

	"github.com/gerladeno/chat-service/internal/types"
)

type Request struct {
	ID        types.RequestID
	ManagerID types.UserID
}

func (r Request) Validate() error {
	if err := r.ID.Validate(); err != nil {
		return fmt.Errorf("validate get chats request id: %v", err)
	}
	if err := r.ManagerID.Validate(); err != nil {
		return fmt.Errorf("validate get chats request managerId: %v", err)
	}
	return nil
}

type Response struct {
	Chats []Chat
}

type Chat struct {
	ID       types.ChatID
	ClientID types.UserID
}
