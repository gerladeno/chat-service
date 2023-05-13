package canreceiveproblems_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/gerladeno/chat-service/internal/testingh"
	"github.com/gerladeno/chat-service/internal/types"
	canreceiveproblems "github.com/gerladeno/chat-service/internal/usecases/manager/can-receive-problems"
	canreceiveproblemsmocks "github.com/gerladeno/chat-service/internal/usecases/manager/can-receive-problems/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl      *gomock.Controller
	mLoadMock *canreceiveproblemsmocks.MockmanagerLoadService
	mPoolMock *canreceiveproblemsmocks.MockmanagerPool
	uCase     *canreceiveproblems.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	var err error
	s.ctrl = gomock.NewController(s.T())
	s.mPoolMock = canreceiveproblemsmocks.NewMockmanagerPool(s.ctrl)
	s.mLoadMock = canreceiveproblemsmocks.NewMockmanagerLoadService(s.ctrl)
	s.uCase, err = canreceiveproblems.New(canreceiveproblems.NewOptions(s.mLoadMock, s.mPoolMock))
	s.Require().NoError(err)
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *UseCaseSuite) TestHandler() {
	req := canreceiveproblems.Request{
		ID:        types.RequestID{},
		ManagerID: types.UserID{},
	}

	s.Run("invalid request", func() {
		_, err := s.uCase.Handle(s.Ctx, req)
		s.Require().Error(err)
	})

	req.ID = types.NewRequestID()
	req.ManagerID = types.NewUserID()
	s.Run("managerPool.Contains err", func() {
		s.mPoolMock.EXPECT().Contains(s.Ctx, req.ManagerID).Return(false, errors.New("wuuut"))
		_, err := s.uCase.Handle(s.Ctx, req)
		s.Require().Error(err)
	})

	s.Run("manager already in a pool", func() {
		s.mPoolMock.EXPECT().Contains(s.Ctx, req.ManagerID).Return(true, nil)
		resp, err := s.uCase.Handle(s.Ctx, req)
		s.Require().NoError(err)
		s.Require().False(resp.Result)
	})

	s.Run("managerLoadService.CanManagerTakeProblem err", func() {
		s.mPoolMock.EXPECT().Contains(s.Ctx, req.ManagerID).Return(false, nil)
		s.mLoadMock.EXPECT().CanManagerTakeProblem(s.Ctx, req.ManagerID).Return(false, errors.New("bang"))
		_, err := s.uCase.Handle(s.Ctx, req)
		s.Require().Error(err)
	})

	s.Run("success true", func() {
		s.mPoolMock.EXPECT().Contains(s.Ctx, req.ManagerID).Return(false, nil)
		s.mLoadMock.EXPECT().CanManagerTakeProblem(s.Ctx, req.ManagerID).Return(true, nil)
		resp, err := s.uCase.Handle(s.Ctx, req)
		s.Require().NoError(err)
		s.Require().True(resp.Result)
	})

	s.Run("success false", func() {
		s.mPoolMock.EXPECT().Contains(s.Ctx, req.ManagerID).Return(false, nil)
		s.mLoadMock.EXPECT().CanManagerTakeProblem(s.Ctx, req.ManagerID).Return(false, nil)
		resp, err := s.uCase.Handle(s.Ctx, req)
		s.Require().NoError(err)
		s.Require().False(resp.Result)
	})
}
