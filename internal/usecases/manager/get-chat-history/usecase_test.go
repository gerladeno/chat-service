package getchathistory_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/gerladeno/chat-service/internal/cursor"
	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	"github.com/gerladeno/chat-service/internal/testingh"
	"github.com/gerladeno/chat-service/internal/types"
	getchathistory "github.com/gerladeno/chat-service/internal/usecases/manager/get-chat-history"
	getchathistorymocks "github.com/gerladeno/chat-service/internal/usecases/manager/get-chat-history/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl    *gomock.Controller
	msgRepo *getchathistorymocks.MockmessagesRepository
	uCase   getchathistory.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.msgRepo = getchathistorymocks.NewMockmessagesRepository(s.ctrl)

	var err error
	s.uCase, err = getchathistory.New(getchathistory.NewOptions(s.msgRepo))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestRequestValidationError() {
	// Arrange.
	req := getchathistory.Request{}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, getchathistory.ErrInvalidRequest)
	s.Empty(resp.Messages)
	s.Empty(resp.NextCursor)
}

func (s *UseCaseSuite) TestCursorDecodingError() {
	// Arrange.
	req := getchathistory.Request{
		ChatID:    types.NewChatID(),
		ManagerID: types.NewUserID(),
		Cursor:    "eyJwYWdlX3NpemUiOjEwMA==", // {"page_size":100
	}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, getchathistory.ErrInvalidCursor)
	s.Empty(resp.Messages)
	s.Empty(resp.NextCursor)
}

func (s *UseCaseSuite) TestGetClientChatMessages_InvalidCursor() {
	// Arrange.
	chatID := types.NewChatID()
	managerID := types.NewUserID()

	c := messagesrepo.Cursor{PageSize: -1, LastCreatedAt: time.Now()}
	cursorWithNegativePageSize, err := cursor.Encode(c)
	s.Require().NoError(err)

	s.msgRepo.EXPECT().GetManagerChatMessages(s.Ctx, chatID, managerID, 0, messagesrepo.NewCursorMatcher(c)).
		Return(nil, nil, messagesrepo.ErrInvalidCursor)

	req := getchathistory.Request{
		ChatID:    chatID,
		ManagerID: managerID,
		PageSize:  0,
		Cursor:    cursorWithNegativePageSize,
	}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, getchathistory.ErrInvalidCursor)
	s.Empty(resp.Messages)
	s.Empty(resp.NextCursor)
}

func (s *UseCaseSuite) TestGetClientChatMessages_SomeError() {
	// Arrange.
	chatID := types.NewChatID()
	managerID := types.NewUserID()
	errExpected := errors.New("any error")

	s.msgRepo.EXPECT().GetManagerChatMessages(s.Ctx, chatID, managerID, 20, (*messagesrepo.Cursor)(nil)).
		Return(nil, nil, errExpected)

	req := getchathistory.Request{
		ChatID:    chatID,
		ManagerID: managerID,
		PageSize:  20,
	}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.Empty(resp.Messages)
	s.Empty(resp.NextCursor)
}

func (s *UseCaseSuite) TestGetClientChatMessages_Success_SinglePage() {
	// Arrange.
	const messagesCount = 10
	const pageSize = messagesCount + 1

	chatID := types.NewChatID()
	clientID := types.NewUserID()
	managerID := types.NewUserID()
	expectedMsgs := s.createMessages(messagesCount, clientID, chatID)

	s.msgRepo.EXPECT().GetManagerChatMessages(s.Ctx, chatID, managerID, pageSize, (*messagesrepo.Cursor)(nil)).
		Return(expectedMsgs, nil, nil)

	req := getchathistory.Request{
		ChatID:    chatID,
		ManagerID: managerID,
		PageSize:  pageSize,
	}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)
	s.Require().NoError(err)

	// Assert.
	s.Empty(resp.NextCursor)

	s.Require().Len(resp.Messages, messagesCount)
	for i := 0; i < messagesCount; i++ {
		s.Equal(expectedMsgs[i].ID, resp.Messages[i].ID)
		s.Equal(expectedMsgs[i].AuthorID, resp.Messages[i].AuthorID)
		s.Equal(expectedMsgs[i].Body, resp.Messages[i].Body)
		s.Equal(expectedMsgs[i].CreatedAt.Unix(), resp.Messages[i].CreatedAt.Unix())
	}
}

func (s *UseCaseSuite) TestGetClientChatMessages_Success_LastPage() {
	// Arrange.
	const messagesCount = 10
	const pageSize = messagesCount + 1

	chatID := types.NewChatID()
	clientID := types.NewUserID()
	managerID := types.NewUserID()
	expectedMsgs := s.createMessages(messagesCount, clientID, chatID)

	c := messagesrepo.Cursor{PageSize: pageSize, LastCreatedAt: time.Now()}
	s.msgRepo.EXPECT().GetManagerChatMessages(s.Ctx, chatID, managerID, 0, messagesrepo.NewCursorMatcher(c)).
		Return(expectedMsgs, nil, nil)

	cursorStr, err := cursor.Encode(c)
	s.Require().NoError(err)

	req := getchathistory.Request{
		ChatID:    chatID,
		ManagerID: managerID,
		Cursor:    cursorStr,
	}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)
	s.Require().NoError(err)

	// Assert.
	s.Empty(resp.NextCursor)
	s.Require().Len(resp.Messages, messagesCount)
}

func (s *UseCaseSuite) createMessages(count int, authorID types.UserID, chatID types.ChatID) []messagesrepo.Message {
	s.T().Helper()

	result := make([]messagesrepo.Message, 0, count)
	for i := 0; i < count; i++ {
		result = append(result, messagesrepo.Message{
			ID:        types.NewMessageID(),
			ChatID:    chatID,
			AuthorID:  authorID,
			Body:      uuid.New().String(),
			CreatedAt: time.Now(),
		})
	}
	return result
}
