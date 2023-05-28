package problems

import (
	"time"

	"github.com/gerladeno/chat-service/internal/store"
	"github.com/gerladeno/chat-service/internal/types"
)

type Problem struct {
	ID         types.ProblemID
	ChatID     types.ChatID
	ManagerID  types.UserID
	ResolvedAt time.Time
	CreatedAt  time.Time
}

func adaptStoreProblem(p *store.Problem) Problem {
	return Problem{
		ID:         p.ID,
		ChatID:     p.ChatID,
		ManagerID:  p.ManagerID,
		ResolvedAt: p.ResolvedAt,
		CreatedAt:  p.CreatedAt,
	}
}
