package sendmanagermessagejob_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	eventstream "github.com/gerladeno/chat-service/internal/services/event-stream"
	msgproducer "github.com/gerladeno/chat-service/internal/services/msg-producer"
	sendmanagermessagejob "github.com/gerladeno/chat-service/internal/services/outbox/jobs/send-manager-message"
	sendmanagermessagejobmocks "github.com/gerladeno/chat-service/internal/services/outbox/jobs/send-manager-message/mocks"
	"github.com/gerladeno/chat-service/internal/types"
)

func TestJob_Handle(t *testing.T) {
	// Arrange.
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chatRepo := sendmanagermessagejobmocks.NewMockchatsRepository(ctrl)
	msgRepo := sendmanagermessagejobmocks.NewMockmessageRepository(ctrl)
	eventStream := sendmanagermessagejobmocks.NewMockeventStream(ctrl)
	msgProducer := sendmanagermessagejobmocks.NewMockmessageProducer(ctrl)
	job, err := sendmanagermessagejob.New(sendmanagermessagejob.NewOptions(
		chatRepo,
		msgRepo,
		eventStream,
		msgProducer,
	))
	require.NoError(t, err)

	clientID := types.NewUserID()
	msgID := types.NewMessageID()
	chatID := types.NewChatID()
	managerID := types.NewUserID()
	msg := messagesrepo.Message{
		ID:                  msgID,
		ChatID:              chatID,
		AuthorID:            types.UserIDNil,
		Body:                `Yo client!`,
		CreatedAt:           time.Now(),
		IsVisibleForClient:  true,
		IsVisibleForManager: true,
		IsBlocked:           false,
		IsService:           false,
	}

	msgRepo.EXPECT().GetMessageByID(ctx, msgID).Return(&msg, nil)
	chatRepo.EXPECT().GetUserID(ctx, chatID).Return(clientID, nil)
	msgProducer.EXPECT().ProduceMessage(ctx, msgproducer.Message{
		ID:         msgID,
		ChatID:     msg.ChatID,
		Body:       msg.Body,
		FromClient: false,
	}).Return(nil)
	eventStream.EXPECT().Publish(ctx, clientID, eventstream.NewNewMessageEvent(
		types.NewEventID(),
		msg.RequestID,
		msg.ChatID,
		msgID,
		managerID,
		msg.CreatedAt,
		msg.Body,
		false,
	)).Return(nil)
	eventStream.EXPECT().Publish(ctx, managerID, eventstream.NewNewMessageEvent(
		types.NewEventID(),
		msg.RequestID,
		msg.ChatID,
		msgID,
		managerID,
		msg.CreatedAt,
		msg.Body,
		false,
	)).Return(nil)
	eventStream.EXPECT().Publish(ctx, managerID, eventstream.NewMessageSentEvent(
		types.NewEventID(),
		msg.RequestID,
		msgID,
	)).Return(nil)

	// Action & assert.
	payload := sendmanagermessagejob.NewPayload(msgID, managerID)
	err = job.Handle(ctx, payload)
	require.NoError(t, err)
}
