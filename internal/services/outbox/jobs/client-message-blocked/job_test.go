package clientmessageblockedjob_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	eventstream "github.com/gerladeno/chat-service/internal/services/event-stream"
	"github.com/gerladeno/chat-service/internal/services/outbox"
	clientmessageblockedjob "github.com/gerladeno/chat-service/internal/services/outbox/jobs/client-message-blocked"
	clientmessageblockedjobmocks "github.com/gerladeno/chat-service/internal/services/outbox/jobs/client-message-blocked/mocks"
	"github.com/gerladeno/chat-service/internal/types"
)

func TestJob_Handle(t *testing.T) {
	// Arrange.
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	msgRepo := clientmessageblockedjobmocks.NewMockmessageRepository(ctrl)
	eventStream := clientmessageblockedjobmocks.NewMockeventStream(ctrl)
	job, err := clientmessageblockedjob.New(clientmessageblockedjob.NewOptions(msgRepo, eventStream))
	require.NoError(t, err)

	clientID := types.NewUserID()
	msgID := types.NewMessageID()
	chatID := types.NewChatID()
	const body = "Hello!"

	msg := messagesrepo.Message{
		ID:                  msgID,
		ChatID:              chatID,
		AuthorID:            clientID,
		Body:                body,
		CreatedAt:           time.Now(),
		IsVisibleForClient:  true,
		IsVisibleForManager: false,
		IsBlocked:           true,
		IsService:           false,
	}
	msgRepo.EXPECT().GetMessageByID(gomock.Any(), msgID).Return(&msg, nil)
	eventStream.EXPECT().Publish(ctx, msg.AuthorID, eventstream.NewMessageBlockedEvent(
		types.NewEventID(),
		msg.RequestID,
		msg.ID,
	))

	// Action & assert.
	payload, err := outbox.MarshalPayload(msgID)
	require.NoError(t, err)
	err = job.Handle(ctx, payload)
	require.NoError(t, err)
}
