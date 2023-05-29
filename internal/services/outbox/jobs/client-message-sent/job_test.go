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

func TestJob_Handle_NoManager(t *testing.T) {
	// Arrange.
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	msgRepo := clientmessagesentjobmocks.NewMockmessageRepository(ctrl)
	problemsRepo := clientmessagesentjobmocks.NewMockproblemsRepository(ctrl)
	eventStream := clientmessagesentjobmocks.NewMockeventStream(ctrl)
	job, err := clientmessagesentjob.New(clientmessagesentjob.NewOptions(msgRepo, problemsRepo, eventStream))
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
	problemsRepo.EXPECT().GetActiveManager(gomock.Any(), chatID).Return(types.UserIDNil, nil)
	eventStream.EXPECT().Publish(gomock.Any(), msg.AuthorID, eventstream.NewMessageSentEvent(
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

func TestJob_Handle_ManagerAssigned(t *testing.T) {
	// Arrange.
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	msgRepo := clientmessagesentjobmocks.NewMockmessageRepository(ctrl)
	problemsRepo := clientmessagesentjobmocks.NewMockproblemsRepository(ctrl)
	eventStream := clientmessagesentjobmocks.NewMockeventStream(ctrl)
	job, err := clientmessagesentjob.New(clientmessagesentjob.NewOptions(msgRepo, problemsRepo, eventStream))
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
	managerID := types.NewUserID()
	msgRepo.EXPECT().GetMessageByID(gomock.Any(), msgID).Return(&msg, nil)
	problemsRepo.EXPECT().GetActiveManager(gomock.Any(), chatID).Return(managerID, nil)
	eventStream.EXPECT().Publish(gomock.Any(), msg.AuthorID, eventstream.NewMessageSentEvent(
		types.NewEventID(),
		msg.RequestID,
		msg.ID,
	))
	eventStream.EXPECT().Publish(gomock.Any(), managerID, eventstream.NewNewMessageEvent(
		types.NewEventID(),
		msg.RequestID,
		msg.ChatID,
		msg.ID,
		msg.AuthorID,
		msg.CreatedAt,
		msg.Body,
		false))

	// Action & assert.
	payload, err := outbox.MarshalPayload(msgID)
	require.NoError(t, err)
	err = job.Handle(ctx, payload)
	require.NoError(t, err)
}
