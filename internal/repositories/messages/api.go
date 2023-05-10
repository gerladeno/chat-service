package messagesrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/gerladeno/chat-service/internal/store"
	"github.com/gerladeno/chat-service/internal/store/message"
	"github.com/gerladeno/chat-service/internal/types"
)

var ErrMsgNotFound = errors.New("message not found")

func (r *Repo) GetMessageByRequestID(ctx context.Context, reqID types.RequestID) (*Message, error) {
	msg, err := r.db.Message(ctx).Query().Where(message.InitialRequestID(reqID)).Only(ctx)
	switch {
	case store.IsNotFound(err):
		return nil, ErrMsgNotFound
	case err != nil:
		return nil, fmt.Errorf("getting msg by reqID: %v", err)
	}
	result := adaptStoreMessage(msg)
	return &result, nil
}

// CreateClientVisible creates a message that is visible only to the client.
func (r *Repo) CreateClientVisible(
	ctx context.Context,
	reqID types.RequestID,
	problemID types.ProblemID,
	chatID types.ChatID,
	authorID types.UserID,
	msgBody string,
) (*Message, error) {
	msg, err := r.db.Message(ctx).Create().
		SetID(types.NewMessageID()).
		SetInitialRequestID(reqID).
		SetProblemID(problemID).
		SetChatID(chatID).
		SetAuthorID(authorID).
		SetBody(msgBody).
		SetIsVisibleForClient(true).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create client-visible message: %v", err)
	}
	result := adaptStoreMessage(msg)
	return &result, nil
}
