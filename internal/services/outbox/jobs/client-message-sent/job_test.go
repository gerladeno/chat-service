package clientmessagesentjob_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	eventstream "github.com/gerladeno/chat-service/internal/services/event-stream"
	"github.com/gerladeno/chat-service/internal/services/outbox"
	clientmessagesentjob "github.com/gerladeno/chat-service/internal/services/outbox/jobs/client-message-sent"
	clientmessagesentjobmocks "github.com/gerladeno/chat-service/internal/services/outbox/jobs/client-message-sent/mocks"
	"github.com/gerladeno/chat-service/internal/types"
)

func TestJob_Handle(t *testing.T) {
	// Arrange.
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	msgRepo := clientmessagesentjobmocks.NewMockmessageRepository(ctrl)
	eventStream := clientmessagesentjobmocks.NewMockeventStream(ctrl)
	job, err := clientmessagesentjob.New(clientmessagesentjob.NewOptions(msgRepo, eventStream))
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
		IsBlocked:           false,
		IsService:           false,
	}
	msgRepo.EXPECT().GetMessageByID(gomock.Any(), msgID).Return(&msg, nil)
	eventStream.EXPECT().Publish(ctx, msg.AuthorID, eventstream.NewMessageSentEvent(
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
