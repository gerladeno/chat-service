package closechatjob_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	eventstream "github.com/gerladeno/chat-service/internal/services/event-stream"
	closechatjob "github.com/gerladeno/chat-service/internal/services/outbox/jobs/close-chat"
	closechatjobmocks "github.com/gerladeno/chat-service/internal/services/outbox/jobs/close-chat/mocks"
	"github.com/gerladeno/chat-service/internal/types"
)

func TestJob_Handle(t *testing.T) {
	// Arrange.
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatsRepo := closechatjobmocks.NewMockchatsRepo(ctrl)
	msgRepo := closechatjobmocks.NewMockmessageRepository(ctrl)
	eventStream := closechatjobmocks.NewMockeventStream(ctrl)
	managerLoad := closechatjobmocks.NewMockmanagerLoad(ctrl)
	job, err := closechatjob.New(closechatjob.NewOptions(chatsRepo, msgRepo, eventStream, managerLoad))
	require.NoError(t, err)

	requestID := types.NewRequestID()
	clientID := types.NewUserID()
	messageID := types.NewMessageID()
	managerID := types.NewUserID()
	chatID := types.NewChatID()
	ts := time.Now()
	body := `Your question has been marked as resolved.
Thank you for being with us!`
	msg := messagesrepo.Message{
		ID:                  messageID,
		RequestID:           requestID,
		ChatID:              chatID,
		AuthorID:            types.UserIDNil,
		Body:                body,
		CreatedAt:           ts,
		IsVisibleForClient:  true,
		IsVisibleForManager: false,
		IsBlocked:           false,
		IsService:           true,
	}
	msgRepo.EXPECT().GetMessageByID(ctx, messageID).Return(&msg, nil)
	chatsRepo.EXPECT().GetClientID(ctx, chatID).Return(clientID, nil)
	eventStream.EXPECT().Publish(ctx, clientID, eventstream.NewNewMessageEvent(
		types.NewEventID(),
		requestID,
		chatID,
		messageID,
		types.UserIDNil,
		ts,
		body,
		true,
	)).Return(nil)
	managerLoad.EXPECT().CanManagerTakeProblem(ctx, managerID).Return(true, nil)
	eventStream.EXPECT().Publish(ctx, managerID, eventstream.NewChatClosedEvent(
		types.NewEventID(),
		requestID,
		chatID,
		true,
	))

	// act
	err = job.Handle(ctx, closechatjob.NewPayload(requestID, chatID, managerID, messageID))

	// assert
	require.NoError(t, err)
}
