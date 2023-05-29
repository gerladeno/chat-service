package problems

import (
	"context"
	"errors"
	"fmt"

	"github.com/gerladeno/chat-service/internal/store"
	"github.com/gerladeno/chat-service/internal/store/chat"
	"github.com/gerladeno/chat-service/internal/store/message"
	"github.com/gerladeno/chat-service/internal/store/problem"
	"github.com/gerladeno/chat-service/internal/types"
)

var ErrProblemNotFound = errors.New("problem not found")

func (r *Repo) CreateIfNotExists(ctx context.Context, chatID types.ChatID) (types.ProblemID, error) {
	pID, err := r.db.Problem(ctx).Query().
		Unique(false).
		Where(
			problem.HasChatWith(chat.ID(chatID)),
			problem.ResolvedAtIsNil(),
		).
		FirstID(ctx)
	if nil == err {
		return pID, nil
	}
	if !store.IsNotFound(err) {
		return types.ProblemIDNil, fmt.Errorf("select existent problem: %v", err)
	}

	p, err := r.db.Problem(ctx).Create().
		SetChatID(chatID).
		Save(ctx)
	if err != nil {
		return types.ProblemIDNil, fmt.Errorf("create new problem: %v", err)
	}

	return p.ID, nil
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

func (r *Repo) GetProblemsWithoutManager(ctx context.Context) ([]Problem, error) {
	problems, err := r.db.Problem(ctx).Query().Where(
		problem.ManagerIDIsNil(),
		problem.HasMessagesWith(message.IsVisibleForManager(true)),
	).Order(
		problem.ByCreatedAt(),
	).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all problems without assigned manager: %v", err)
	}
	result := make([]Problem, 0, len(problems))
	for _, p := range problems {
		result = append(result, adaptStoreProblem(p))
	}
	return result, nil
}

func (r *Repo) AssignManager(ctx context.Context, problemID types.ProblemID, managerID types.UserID) error {
	_, err := r.db.Problem(ctx).UpdateOneID(problemID).SetManagerID(managerID).Save(ctx)
	switch {
	case store.IsNotFound(err):
		return ErrProblemNotFound
	case err != nil:
		return fmt.Errorf("assign manager for problem: %v", err)
	}
	return nil
}

func (r *Repo) GetRequestID(ctx context.Context, problemID types.ProblemID) (types.RequestID, error) {
	msg, err := r.db.Message(ctx).Query().Where(
		message.ProblemID(problemID),
	).Order(message.ByCreatedAt()).First(ctx)
	switch {
	case store.IsNotFound(err):
		return types.RequestIDNil, ErrProblemNotFound
	case err != nil:
		return types.RequestIDNil, fmt.Errorf("get request id by problem id: %v", err)
	}
	return msg.InitialRequestID, nil
}

func (r *Repo) GetActiveManager(ctx context.Context, chatID types.ChatID) (types.UserID, error) {
	p, err := r.db.Problem(ctx).Query().Where(
		problem.ChatID(chatID),
	).First(ctx)
	switch {
	case store.IsNotFound(err):
		return types.UserIDNil, ErrProblemNotFound
	case err != nil:
		return types.UserIDNil, fmt.Errorf("get chat's active manager: %v", err)
	default:
		return p.ManagerID, nil
	}
}
