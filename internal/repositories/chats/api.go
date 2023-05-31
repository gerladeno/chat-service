package chats

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gerladeno/chat-service/internal/store"
	"github.com/gerladeno/chat-service/internal/store/chat"
	"github.com/gerladeno/chat-service/internal/store/problem"
	"github.com/gerladeno/chat-service/internal/types"
)

var ErrChatNotFound = errors.New("chat not found")

type Chat struct {
	ID       types.ChatID
	ClientID types.UserID
}

func (r *Repo) CreateIfNotExists(ctx context.Context, userID types.UserID) (types.ChatID, error) {
	chatID, err := r.db.Chat(ctx).Create().
		SetCreatedAt(time.Now()).
		SetClientID(userID).OnConflictColumns("client_id").UpdateClientID().
		ID(ctx)
	if err != nil {
		return types.NewChatID(), fmt.Errorf("upserting chat: %v", err)
	}
	return chatID, nil
}

func (r *Repo) GetUserID(ctx context.Context, chatID types.ChatID) (types.UserID, error) {
	chatEntity, err := r.db.Chat(ctx).Get(ctx, chatID)
	switch {
	case store.IsNotFound(err):
		return types.UserIDNil, ErrChatNotFound
	case err != nil:
		return types.UserIDNil, fmt.Errorf("get chat by id: %v", err)
	}
	return chatEntity.ClientID, nil
}

func (r *Repo) GetChatsForManager(ctx context.Context, managerID types.UserID) ([]Chat, error) {
	chats, err := r.db.Chat(ctx).Query().Where(
		chat.HasProblemsWith(problem.ResolvedAtIsNil(), problem.ManagerID(managerID)),
	).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("get unresolved chats by manager: %v", err)
	}
	result := make([]Chat, 0, len(chats))
	for _, chatEntity := range chats {
		result = append(result, Chat{
			ID:       chatEntity.ID,
			ClientID: chatEntity.ClientID,
		})
	}
	return result, nil
}

func (r *Repo) GetClientID(ctx context.Context, chatID types.ChatID) (types.UserID, error) {
	c, err := r.db.Chat(ctx).Get(ctx, chatID)
	switch {
	case store.IsNotFound(err):
		return types.UserIDNil, ErrChatNotFound
	case err != nil:
		return types.UserIDNil, fmt.Errorf("get chat by id: %v", err)
	}
	return c.ClientID, nil
}
