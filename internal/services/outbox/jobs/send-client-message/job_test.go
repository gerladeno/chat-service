package sendclientmessagejob_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	eventstream "github.com/gerladeno/chat-service/internal/services/event-stream"
	msgproducer "github.com/gerladeno/chat-service/internal/services/msg-producer"
	sendclientmessagejob "github.com/gerladeno/chat-service/internal/services/outbox/jobs/send-client-message"
	sendclientmessagejobmocks "github.com/gerladeno/chat-service/internal/services/outbox/jobs/send-client-message/mocks"
	"github.com/gerladeno/chat-service/internal/types"
)

func TestJob_Handle(t *testing.T) {
	// Arrange.
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	msgProducer := sendclientmessagejobmocks.NewMockmessageProducer(ctrl)
	msgRepo := sendclientmessagejobmocks.NewMockmessageRepository(ctrl)
	eventStream := sendclientmessagejobmocks.NewMockeventStream(ctrl)
	job, err := sendclientmessagejob.New(sendclientmessagejob.NewOptions(msgProducer, msgRepo, eventStream))
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

	msgProducer.EXPECT().ProduceMessage(gomock.Any(), msgproducer.Message{
		ID:         msgID,
		ChatID:     chatID,
		Body:       body,
		FromClient: true,
	}).Return(nil)

	eventStream.EXPECT().Publish(ctx, msg.AuthorID, eventstream.NewNewMessageEvent(
		types.NewEventID(),
		msg.RequestID,
		msg.ChatID,
		msg.ID,
		msg.AuthorID,
		msg.CreatedAt,
		msg.Body,
		msg.IsService,
	))

	// Action & assert.
	payload, err := sendclientmessagejob.MarshalPayload(msgID)
	require.NoError(t, err)

	err = job.Handle(ctx, payload)
	require.NoError(t, err)
}
