package messagesrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/gerladeno/chat-service/internal/store"
	"github.com/gerladeno/chat-service/internal/types"
)

func (r *Repo) MarkAsVisibleForManager(ctx context.Context, msgID types.MessageID) error {
	_, err := r.db.Message(ctx).UpdateOneID(msgID).SetCheckedAt(time.Now()).SetIsVisibleForManager(true).Save(ctx)
	switch {
	case store.IsNotFound(err):
		return ErrMsgNotFound
	case err != nil:
		return fmt.Errorf("set msg visible for manager: %v", err)
	}
	return nil
}

func (r *Repo) BlockMessage(ctx context.Context, msgID types.MessageID) error {
	_, err := r.db.Message(ctx).UpdateOneID(msgID).SetCheckedAt(time.Now()).SetIsBlocked(true).Save(ctx)
	switch {
	case store.IsNotFound(err):
		return ErrMsgNotFound
	case err != nil:
		return fmt.Errorf("set msg is blocked: %v", err)
	}
	return nil
}
