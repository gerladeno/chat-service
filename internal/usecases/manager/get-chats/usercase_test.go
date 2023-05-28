package getchats_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	chatsrepo "github.com/gerladeno/chat-service/internal/repositories/chats"
	"github.com/gerladeno/chat-service/internal/testingh"
	"github.com/gerladeno/chat-service/internal/types"
	getchats "github.com/gerladeno/chat-service/internal/usecases/manager/get-chats"
	getchatsmock "github.com/gerladeno/chat-service/internal/usecases/manager/get-chats/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl      *gomock.Controller
	chatsRepo *getchatsmock.MockchatsRepo
	uCase     getchats.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	var err error
	s.ctrl = gomock.NewController(s.T())
	s.chatsRepo = getchatsmock.NewMockchatsRepo(s.ctrl)
	s.uCase, err = getchats.New(getchats.NewOptions(s.chatsRepo))
	s.Require().NoError(err)
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *UseCaseSuite) TestHandle() {
	s.Run("invalid request", func() {
		_, err := s.uCase.Handle(s.Ctx, getchats.Request{
			ID:        types.RequestID{},
			ManagerID: types.UserID{},
		})
		s.Require().Error(err)
	})

	s.Run("repo error", func() {
		managerID := types.NewUserID()
		s.chatsRepo.EXPECT().GetChatsForManager(s.Ctx, managerID).Return(nil, errors.New("opa"))
		_, err := s.uCase.Handle(s.Ctx, getchats.Request{
			ID:        types.NewRequestID(),
			ManagerID: managerID,
		})
		s.Require().Error(err)
	})

	s.Run("positive", func() {
		managerID := types.NewUserID()
		s.chatsRepo.EXPECT().GetChatsForManager(s.Ctx, managerID).Return([]chatsrepo.Chat{{}}, nil)
		resp, err := s.uCase.Handle(s.Ctx, getchats.Request{
			ID:        types.NewRequestID(),
			ManagerID: managerID,
		})
		s.Require().NoError(err)
		s.Require().Len(resp.Chats, 1)
	})
}
