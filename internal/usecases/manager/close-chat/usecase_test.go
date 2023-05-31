package closechat_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	problemsrepo "github.com/gerladeno/chat-service/internal/repositories/problems"
	closechatjob "github.com/gerladeno/chat-service/internal/services/outbox/jobs/close-chat"
	"github.com/gerladeno/chat-service/internal/testingh"
	"github.com/gerladeno/chat-service/internal/types"
	closechat "github.com/gerladeno/chat-service/internal/usecases/manager/close-chat"
	closechatmocks "github.com/gerladeno/chat-service/internal/usecases/manager/close-chat/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl          *gomock.Controller
	problemsRepo  *closechatmocks.MockproblemsRepo
	msgRepo       *closechatmocks.MockmsgRepo
	outboxService *closechatmocks.MockoutboxService
	tx            *closechatmocks.Mocktransactor
	uCase         closechat.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	var err error
	s.ctrl = gomock.NewController(s.T())
	s.problemsRepo = closechatmocks.NewMockproblemsRepo(s.ctrl)
	s.msgRepo = closechatmocks.NewMockmsgRepo(s.ctrl)
	s.tx = closechatmocks.NewMocktransactor(s.ctrl)
	s.outboxService = closechatmocks.NewMockoutboxService(s.ctrl)
	s.uCase, err = closechat.New(closechat.NewOptions(s.msgRepo, s.problemsRepo, s.outboxService, s.tx))
	s.Require().NoError(err)
	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TestHandle() {
	s.Run("invalid request", func() {
		err := s.uCase.Handle(s.Ctx, closechat.Request{
			ID:        types.RequestID{},
			ManagerID: types.UserID{},
		})
		s.Require().Error(err)
	})

	s.Run("repo error", func() {
		managerID := types.NewUserID()
		chatID := types.NewChatID()
		requestID := types.NewRequestID()
		s.tx.EXPECT().RunInTx(s.Ctx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, f func(ctx context.Context) error) error {
				return f(ctx)
			})
		s.problemsRepo.EXPECT().ResolveProblem(s.Ctx, chatID, managerID, requestID).
			Return(types.ProblemIDNil, errors.New("opa"))
		err := s.uCase.Handle(s.Ctx, closechat.Request{
			ID:        requestID,
			ManagerID: managerID,
			ChatID:    chatID,
		})
		s.Require().Error(err)
	})

	s.Run("repo ErrProblemNotFound", func() {
		// arrange
		managerID := types.NewUserID()
		chatID := types.NewChatID()
		reqID := types.NewRequestID()
		s.tx.EXPECT().RunInTx(s.Ctx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, f func(ctx context.Context) error) error {
				return f(ctx)
			})
		s.problemsRepo.EXPECT().ResolveProblem(s.Ctx, chatID, managerID, reqID).
			Return(types.ProblemIDNil, problemsrepo.ErrProblemNotFound)

		// act
		err := s.uCase.Handle(s.Ctx, closechat.Request{
			ID:        reqID,
			ManagerID: managerID,
			ChatID:    chatID,
		})

		// assert
		s.Require().Error(err)
		s.Require().ErrorIs(err, closechat.ErrNoActiveProblem)
	})

	s.Run("positive", func() {
		// arrange
		managerID := types.NewUserID()
		chatID := types.NewChatID()
		reqID := types.NewRequestID()
		problemID := types.NewProblemID()
		messageID := types.NewMessageID()
		body := `Your question has been marked as resolved.
Thank you for being with us!`
		s.tx.EXPECT().RunInTx(s.Ctx, gomock.Any()).DoAndReturn(
			func(ctx context.Context, f func(ctx context.Context) error) error {
				return f(ctx)
			})
		s.problemsRepo.EXPECT().ResolveProblem(s.Ctx, chatID, managerID, reqID).Return(problemID, nil)
		s.msgRepo.EXPECT().CreateService(s.Ctx, reqID, problemID, chatID, body).
			Return(&messagesrepo.Message{Body: body, ID: messageID}, nil)
		s.outboxService.EXPECT().Put(s.Ctx, closechatjob.Name,
			closechatjob.NewPayload(reqID, chatID, managerID, messageID), gomock.Any()).Return(types.NewJobID(), nil)

		// act
		err := s.uCase.Handle(s.Ctx, closechat.Request{
			ID:        reqID,
			ManagerID: managerID,
			ChatID:    chatID,
		})

		// assert
		s.Require().NoError(err)
	})
}
