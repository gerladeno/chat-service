package messagesrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/gerladeno/chat-service/internal/store/chat"
	"github.com/gerladeno/chat-service/internal/store/message"
	"github.com/gerladeno/chat-service/internal/types"
)

var (
	ErrInvalidPageSize = errors.New("invalid page size")
	ErrInvalidCursor   = errors.New("invalid cursor")
	ErrInvalidParams   = errors.New("either page size or valid cursor must be present")
)

type Cursor struct {
	LastCreatedAt time.Time
	PageSize      int
}

func (c Cursor) validate() error {
	if err := validatePgeSize(c.PageSize); err != nil || c.LastCreatedAt.IsZero() {
		return ErrInvalidCursor
	}
	return nil
}

// GetClientChatMessages returns Nth page of messages in the chat for client side.
func (r *Repo) GetClientChatMessages(
	ctx context.Context,
	clientID types.UserID,
	pageSize int,
	cursor *Cursor,
) ([]Message, *Cursor, error) {
	query := r.db.Chat(ctx).Query().Where(chat.ClientIDEQ(clientID)).QueryMessages()
	switch {
	case cursor != nil:
		if err := cursor.validate(); err != nil {
			return nil, nil, err
		}
		pageSize = cursor.PageSize
		query = query.Where(message.IsVisibleForClient(true), message.CreatedAtLT(cursor.LastCreatedAt))
	case pageSize != 0:
		if err := validatePgeSize(pageSize); err != nil {
			return nil, nil, err
		}
		query = query.Where(message.IsVisibleForClient(true))
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

func validatePgeSize(pageSize int) error {
	if pageSize < 10 || pageSize > 100 {
		return ErrInvalidPageSize
	}
	return nil
}
