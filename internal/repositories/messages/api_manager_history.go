package messagesrepo

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"

	"github.com/gerladeno/chat-service/internal/store/message"
	"github.com/gerladeno/chat-service/internal/store/problem"
	"github.com/gerladeno/chat-service/internal/types"
)

func (r *Repo) GetManagerChatMessages(
	ctx context.Context,
	chatID types.ChatID,
	managerID types.UserID,
	pageSize int,
	cursor *Cursor,
) ([]Message, *Cursor, error) {
	query := r.db.Message(ctx).Query().Where(
		message.ChatID(chatID),
		message.IsVisibleForManager(true),
		message.IsService(false),
		message.IsBlocked(false),
		message.HasProblemWith(
			problem.ManagerID(managerID),
			problem.ResolvedAtIsNil(),
		),
	)
	switch {
	case cursor != nil:
		if err := cursor.validate(); err != nil {
			return nil, nil, err
		}
		pageSize = cursor.PageSize
		query = query.Where(message.CreatedAtLT(cursor.LastCreatedAt))
	case pageSize != 0:
		if err := validatePageSize(pageSize); err != nil {
			return nil, nil, err
		}
	default:
		return nil, nil, ErrInvalidParams
	}

	messages, err := query.Limit(pageSize + 1).Order(message.ByCreatedAt(sql.OrderDesc())).All(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("querying messages: %v", err)
	}

	l := len(messages)
	if l > pageSize {
		l = pageSize
	}
	result := make([]Message, l)
	for i := range result {
		result[i] = adaptStoreMessage(messages[i])
	}

	cursor = nil
	if len(messages) > len(result) {
		cursor = &Cursor{
			LastCreatedAt: result[l-1].CreatedAt,
			PageSize:      pageSize,
		}
	}
	return result, cursor, nil
}
