package chats

import (
	"context"
	"fmt"
	"time"

	"github.com/gerladeno/chat-service/internal/types"
)

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
