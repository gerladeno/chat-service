package managerload

import (
	"context"
	"fmt"

	"github.com/gerladeno/chat-service/internal/types"
)

func (s *Service) CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error) {
	count, err := s.problemsRepo.GetManagerOpenProblemsCount(ctx, managerID)
	if err != nil {
		return false, fmt.Errorf("getting manager's problems count: %v", err)
	}
	return count < s.maxProblemsAtTime, nil
}
