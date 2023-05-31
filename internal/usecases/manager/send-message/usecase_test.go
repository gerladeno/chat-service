package sendmessage_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	problemsrepo "github.com/gerladeno/chat-service/internal/repositories/problems"
	sendmanagermessagejob "github.com/gerladeno/chat-service/internal/services/outbox/jobs/send-manager-message"
	"github.com/gerladeno/chat-service/internal/testingh"
	"github.com/gerladeno/chat-service/internal/types"
	sendmessage "github.com/gerladeno/chat-service/internal/usecases/manager/send-message"
	sendmessagemocks "github.com/gerladeno/chat-service/internal/usecases/manager/send-message/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl        *gomock.Controller
	msgRepo     *sendmessagemocks.MockmessagesRepository
	problemRepo *sendmessagemocks.MockproblemsRepository
	txtor       *sendmessagemocks.Mocktransactor
	outBoxSvc   *sendmessagemocks.MockoutboxService
	uCase       sendmessage.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.msgRepo = sendmessagemocks.NewMockmessagesRepository(s.ctrl)
	s.outBoxSvc = sendmessagemocks.NewMockoutboxService(s.ctrl)
	s.problemRepo = sendmessagemocks.NewMockproblemsRepository(s.ctrl)
	s.txtor = sendmessagemocks.NewMocktransactor(s.ctrl)

	var err error
	s.uCase, err = sendmessage.New(sendmessage.NewOptions(s.msgRepo, s.outBoxSvc, s.problemRepo, s.txtor))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestRequestValidationError() {
	// Arrange.
	req := sendmessage.Request{}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.Require().ErrorIs(err, sendmessage.ErrInvalidRequest)
}

func (s *UseCaseSuite) TestCreateMessage_ProblemNotFound() {
	// Arrange.
	reqID := types.NewRequestID()
	managerID := types.NewUserID()
	chatID := types.NewChatID()

	s.problemRepo.EXPECT().GetAssignedProblemID(gomock.Any(), managerID, chatID).
		Return(types.ProblemIDNil, problemsrepo.ErrProblemNotFound)

	req := sendmessage.Request{
		ID:          reqID,
		ManagerID:   managerID,
		ChatID:      chatID,
		MessageBody: "wtf",
	}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.Require().ErrorIs(err, sendmessage.ErrProblemNotFound)
}

func (s *UseCaseSuite) TestCreateMessage_UnexpectedMsgCreateError() {
	// Arrange.
	reqID := types.NewRequestID()
	managerID := types.NewUserID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()

	s.problemRepo.EXPECT().GetAssignedProblemID(gomock.Any(), managerID, chatID).Return(problemID, nil)

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.msgRepo.EXPECT().CreateFullVisible(
		gomock.Any(),
		reqID,
		problemID,
		chatID,
		managerID,
		`wtf`,
	).Return(nil, errors.New("unexpected"))

	req := sendmessage.Request{
		ID:          reqID,
		ManagerID:   managerID,
		ChatID:      chatID,
		MessageBody: "wtf",
	}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestCreateMessage_UnexpectedJobCreateError() {
	// Arrange.
	reqID := types.NewRequestID()
	managerID := types.NewUserID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()
	msgID := types.NewMessageID()

	msg := messagesrepo.Message{
		ID:                  msgID,
		ChatID:              chatID,
		AuthorID:            types.UserIDNil,
		Body:                `wtf`,
		CreatedAt:           time.Now(),
		IsVisibleForClient:  true,
		IsVisibleForManager: true,
		IsBlocked:           false,
		IsService:           false,
	}

	s.problemRepo.EXPECT().GetAssignedProblemID(gomock.Any(), managerID, chatID).Return(problemID, nil)

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})

	s.msgRepo.EXPECT().CreateFullVisible(
		gomock.Any(),
		reqID,
		problemID,
		chatID,
		managerID,
		`wtf`,
	).Return(&msg, nil)

	s.outBoxSvc.EXPECT().
		Put(gomock.Any(), sendmanagermessagejob.Name, sendmanagermessagejob.NewPayload(msgID, managerID), gomock.Any()).
		Return(types.JobIDNil, errors.New("surprise motherfarmer"))

	req := sendmessage.Request{
		ID:          reqID,
		ManagerID:   managerID,
		ChatID:      chatID,
		MessageBody: "wtf",
	}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestCreateMessage_TxCommitError() {
	// Arrange.
	reqID := types.NewRequestID()
	managerID := types.NewUserID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()
	msgID := types.NewMessageID()

	msg := messagesrepo.Message{
		ID:                  msgID,
		ChatID:              chatID,
		AuthorID:            types.UserIDNil,
		Body:                `wtf`,
		CreatedAt:           time.Now(),
		IsVisibleForClient:  true,
		IsVisibleForManager: true,
		IsBlocked:           false,
		IsService:           false,
	}

	s.problemRepo.EXPECT().GetAssignedProblemID(gomock.Any(), managerID, chatID).Return(problemID, nil)

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			_ = f(ctx)
			return sql.ErrTxDone
		})

	s.msgRepo.EXPECT().CreateFullVisible(
		gomock.Any(),
		reqID,
		problemID,
		chatID,
		managerID,
		`wtf`,
	).Return(&msg, nil)

	s.outBoxSvc.EXPECT().
		Put(gomock.Any(), sendmanagermessagejob.Name, sendmanagermessagejob.NewPayload(msgID, managerID), gomock.Any()).
		Return(types.NewJobID(), nil)

	req := sendmessage.Request{
		ID:          reqID,
		ManagerID:   managerID,
		ChatID:      chatID,
		MessageBody: "wtf",
	}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestCreateMessage_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	managerID := types.NewUserID()
	chatID := types.NewChatID()
	problemID := types.NewProblemID()
	msgID := types.NewMessageID()

	msg := messagesrepo.Message{
		ID:                  msgID,
		ChatID:              chatID,
		AuthorID:            types.UserIDNil,
		Body:                `wtf`,
		CreatedAt:           time.Now(),
		IsVisibleForClient:  true,
		IsVisibleForManager: true,
		IsBlocked:           false,
		IsService:           false,
	}

	s.problemRepo.EXPECT().GetAssignedProblemID(gomock.Any(), managerID, chatID).Return(problemID, nil)

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})

	s.msgRepo.EXPECT().CreateFullVisible(
		gomock.Any(),
		reqID,
		problemID,
		chatID,
		managerID,
		`wtf`,
	).Return(&msg, nil)

	s.outBoxSvc.EXPECT().
		Put(gomock.Any(), sendmanagermessagejob.Name, sendmanagermessagejob.NewPayload(msgID, managerID), gomock.Any()).
		Return(types.NewJobID(), nil)

	req := sendmessage.Request{
		ID:          reqID,
		ManagerID:   managerID,
		ChatID:      chatID,
		MessageBody: "wtf",
	}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
	s.Require().Equal(sendmessage.Response{
		MessageID: msgID,
		CreatedAt: msg.CreatedAt,
	}, resp)
}
