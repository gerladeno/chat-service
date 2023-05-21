package canreceiveproblems

import (
	"fmt"

	"github.com/gerladeno/chat-service/internal/types"
)

type Request struct {
	ID        types.RequestID `validate:"required"`
	ManagerID types.UserID    `validate:"required"`
}

func (r Request) Validate() error {
	if err := r.ID.Validate(); err != nil {
		return fmt.Errorf("can receive problems handler request id validation: %v", err)
	}
	if err := r.ManagerID.Validate(); err != nil {
		return fmt.Errorf("can receive problems handler request manager id validation: %v", err)
	}
	return nil
}

type Response struct {
	Result bool
}
