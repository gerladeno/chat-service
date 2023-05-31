package managerassignedtoproblemjob_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	eventstream "github.com/gerladeno/chat-service/internal/services/event-stream"
	managerassignedtoproblemjob "github.com/gerladeno/chat-service/internal/services/outbox/jobs/manager-assigned-to-problem"
	managerassignedtoproblemjobbmocks "github.com/gerladeno/chat-service/internal/services/outbox/jobs/manager-assigned-to-problem/mocks"
	"github.com/gerladeno/chat-service/internal/types"
)

func TestJob_Handle(t *testing.T) {
	// Arrange.
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	msgRepo := managerassignedtoproblemjobbmocks.NewMockmessageRepository(ctrl)
	chatRepo := managerassignedtoproblemjobbmocks.NewMockchatRepository(ctrl)
	managerLoad := managerassignedtoproblemjobbmocks.NewMockmanagerLoad(ctrl)
	eventStream := managerassignedtoproblemjobbmocks.NewMockeventStream(ctrl)
	job, err := managerassignedtoproblemjob.New(managerassignedtoproblemjob.NewOptions(
		msgRepo,
		chatRepo,
		managerLoad,
		eventStream,
	))
	require.NoError(t, err)

	clientID := types.NewUserID()
	msgID := types.NewMessageID()
	chatID := types.NewChatID()
	managerID := types.MustParse[types.UserID]("e354e5b1-29e2-4edb-8485-8b88f7130d20")
	const body = `Manager e354e5b1-29e2-4edb-8485-8b88f7130d20 will answer you`
	msg := messagesrepo.Message{
		ID:                  msgID,
		ChatID:              chatID,
		AuthorID:            types.UserIDNil,
		Body:                body,
		CreatedAt:           time.Now(),
		IsVisibleForClient:  true,
		IsVisibleForManager: false,
		IsBlocked:           false,
		IsService:           true,
	}

	msgRepo.EXPECT().GetMessageByID(ctx, msgID).Return(&msg, nil)
	chatRepo.EXPECT().GetUserID(ctx, chatID).Return(clientID, nil)
	eventStream.EXPECT().Publish(ctx, clientID, eventstream.NewNewMessageEvent(
		types.NewEventID(),
		msg.RequestID,
		msg.ChatID,
		msg.ID,
		msg.AuthorID,
		msg.CreatedAt,
		msg.Body,
		msg.IsService,
	)).Return(nil)
	managerLoad.EXPECT().CanManagerTakeProblem(ctx, managerID).Return(true, nil)
	eventStream.EXPECT().Publish(ctx, managerID, eventstream.NewNewChatEvent(
		types.NewEventID(),
		msg.RequestID,
		msg.ChatID,
		clientID,
		true,
	)).Return(nil)

	// Action & assert.
	payload := managerassignedtoproblemjob.NewPayload(msgID, managerID)
	err = job.Handle(ctx, payload)
	require.NoError(t, err)
}
