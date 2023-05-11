package problems

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/gerladeno/chat-service/internal/store/problem"
	"github.com/gerladeno/chat-service/internal/types"
)

func (r *Repo) CreateIfNotExists(ctx context.Context, chatID types.ChatID) (types.ProblemID, error) {
	problemID, err := r.db.Problem(ctx).Create().
		SetChatID(chatID).
		SetCreatedAt(time.Now()).OnConflict(
		sql.ConflictColumns("chat_id"),
		sql.ConflictWhere(sql.IsNull("resolved_at"))).UpdateChatID().
		ID(ctx)
	if err != nil {
		return types.NewProblemID(), fmt.Errorf("upserting problem: %v", err)
	}
	return problemID, nil
}

func (r *Repo) GetManagerOpenProblemsCount(ctx context.Context, managerID types.UserID) (int, error) {
	count, err := r.db.Problem(ctx).Query().Where(
		problem.ManagerID(managerID),
		problem.ResolvedAtIsNil(),
	).Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("get manager open problems count: %v", err)
	}
	return count, nil
}
