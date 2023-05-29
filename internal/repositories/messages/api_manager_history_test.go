//go:build integration

package messagesrepo_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	"github.com/gerladeno/chat-service/internal/testingh"
	"github.com/gerladeno/chat-service/internal/types"
)

type MsgRepoManagerHistoryAPISuite struct {
	testingh.DBSuite
	repo      *messagesrepo.Repo
	chatID    types.ChatID
	problemID types.ProblemID
	managerID types.UserID
	clientID  types.UserID
}

func TestMsgRepoManagerHistoryAPISuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &MsgRepoManagerHistoryAPISuite{DBSuite: testingh.NewDBSuite("TestMsgRepoManagerHistoryAPISuite")})
}

func (s *MsgRepoManagerHistoryAPISuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	var err error
	s.repo, err = messagesrepo.New(messagesrepo.NewOptions(s.Database))
	s.Require().NoError(err)
	s.Ctx = context.Background()
	s.chatID = s.createChat()
	s.managerID = types.NewUserID()
	s.clientID = types.NewUserID()
	anotherManagerID := types.NewUserID()
	p1 := s.createProblem(s.chatID, anotherManagerID, true)
	p2 := s.createProblem(s.chatID, s.managerID, true)
	s.problemID = s.createProblem(s.chatID, s.managerID, false)
	for i := 0; i < 15; i++ {
		s.createMessage(s.chatID, p1, anotherManagerID,
			fmt.Sprintf(`another_manager_%d`, i), true)
	}
	for i := 0; i < 15; i++ {
		s.createMessage(s.chatID, p1, s.clientID,
			fmt.Sprintf(`client_with_another_manager_%d`, i), true)
	}
	for i := 0; i < 15; i++ {
		s.createMessage(s.chatID, p2, s.clientID,
			fmt.Sprintf(`client_same_manager_resolved_problem_%d`, i), true)
	}
	for i := 0; i < 15; i++ {
		s.createMessage(s.chatID, p2, s.managerID,
			fmt.Sprintf(`same_manager_resolved_problem_%d`, i), true)
	}
	for i := 0; i < 15; i++ {
		s.createMessage(s.chatID, s.problemID, s.managerID,
			fmt.Sprintf(`same_manager_actual_problem_%d`, i), true)
	}
	for i := 0; i < 15; i++ {
		s.createMessage(s.chatID, s.problemID, s.clientID,
			fmt.Sprintf(`invisible_client_same_manager_actual_problem_%d`, i), false)
	}
	for i := 0; i < 15; i++ {
		s.createMessage(s.chatID, s.problemID, s.clientID,
			fmt.Sprintf(`client_same_manager_actual_problem_%d`, i), true)
	}
}

func (s *MsgRepoManagerHistoryAPISuite) TestSomething() {
	s.Run("too small page size", func() {
		msgs, next, err := s.repo.GetManagerChatMessages(s.Ctx, s.chatID, s.managerID, 9, nil)
		s.Require().ErrorIs(err, messagesrepo.ErrInvalidPageSize)
		s.Nil(next)
		s.Empty(msgs)
	})

	s.Run("too big page size", func() {
		msgs, next, err := s.repo.GetManagerChatMessages(s.Ctx, s.chatID, s.managerID, 101, nil)
		s.Require().ErrorIs(err, messagesrepo.ErrInvalidPageSize)
		s.Nil(next)
		s.Empty(msgs)
	})

	s.Run("no last created at in cursor", func() {
		msgs, next, err := s.repo.GetManagerChatMessages(s.Ctx, s.chatID, s.managerID, 0, &messagesrepo.Cursor{
			LastCreatedAt: time.Time{},
			PageSize:      50,
		})
		s.Require().ErrorIs(err, messagesrepo.ErrInvalidCursor)
		s.Nil(next)
		s.Empty(msgs)
	})

	s.Run("too small page size in cursor", func() {
		msgs, next, err := s.repo.GetManagerChatMessages(s.Ctx, s.chatID, s.managerID, 0, &messagesrepo.Cursor{
			LastCreatedAt: time.Now(),
			PageSize:      9,
		})
		s.Require().ErrorIs(err, messagesrepo.ErrInvalidCursor)
		s.Nil(next)
		s.Empty(msgs)
	})

	s.Run("too big page size in cursor", func() {
		msgs, next, err := s.repo.GetManagerChatMessages(s.Ctx, s.chatID, s.managerID, 0, &messagesrepo.Cursor{
			LastCreatedAt: time.Now(),
			PageSize:      101,
		})
		s.Require().ErrorIs(err, messagesrepo.ErrInvalidCursor)
		s.Nil(next)
		s.Empty(msgs)
	})

	s.Run("chat has not got any messages", func() {
		msgs, next, err := s.repo.GetManagerChatMessages(s.Ctx, types.NewChatID(), s.managerID, 50, nil)
		s.Require().NoError(err)
		s.Nil(next)
		s.Empty(msgs)
	})

	s.Run("neither page_size nor cursor", func() {
		msgs, next, err := s.repo.GetManagerChatMessages(s.Ctx, s.chatID, s.managerID, 0, nil)
		s.Require().ErrorIs(err, messagesrepo.ErrInvalidParams)
		s.Nil(next)
		s.Empty(msgs)
	})

	s.Run("check visibility", func() {
		msgs, next, err := s.repo.GetManagerChatMessages(s.Ctx, s.chatID, s.managerID, 100, nil)
		s.Require().NoError(err)
		s.Nil(next)
		s.Require().Len(msgs, 30)
		for _, m := range msgs {
			s.Require().True(strings.HasPrefix(m.Body, "same_manager_actual_problem_") ||
				strings.HasPrefix(m.Body, "client_same_manager_actual_problem_"))
		}
	})

	s.Run("check cursor logic", func() {
		msgs, next, err := s.repo.GetManagerChatMessages(s.Ctx, s.chatID, s.managerID, 10, nil)
		s.Require().NoError(err)
		s.Require().Len(msgs, 10)
		for _, m := range msgs {
			s.Require().True(strings.HasPrefix(m.Body, "client_same_manager_actual_problem_"))
		}
		msgs, next, err = s.repo.GetManagerChatMessages(s.Ctx, s.chatID, s.managerID, 0, next)
		s.Require().NoError(err)
		// second 10 should be a mix of the two
		s.Require().Len(msgs, 10)
		msgs, next, err = s.repo.GetManagerChatMessages(s.Ctx, s.chatID, s.managerID, 0, next)
		s.Require().NoError(err)
		s.Nil(next)
		s.Require().Len(msgs, 10)
		for _, m := range msgs {
			s.Require().True(strings.HasPrefix(m.Body, "same_manager_actual_problem_"))
		}
	})
}

func (s *MsgRepoManagerHistoryAPISuite) createChat() types.ChatID {
	s.T().Helper()
	c, err := s.Database.Chat(s.Ctx).Create().SetClientID(types.NewUserID()).Save(s.Ctx)
	s.Require().NoError(err)
	return c.ID
}

func (s *MsgRepoManagerHistoryAPISuite) createProblem(chatID types.ChatID, managerID types.UserID, resolved bool,
) types.ProblemID {
	s.T().Helper()
	builder := s.Database.Problem(s.Ctx).Create().SetChatID(chatID).SetManagerID(managerID)
	if resolved {
		builder = builder.SetResolvedAt(time.Now())
	}
	p, err := builder.Save(s.Ctx)
	s.Require().NoError(err)
	return p.ID
}

func (s *MsgRepoManagerHistoryAPISuite) createMessage(
	chatID types.ChatID,
	problemID types.ProblemID,
	authorID types.UserID,
	body string,
	visibleForManager bool,
) {
	_, err := s.Database.Message(s.Ctx).Create().
		SetBody(body).
		SetProblemID(problemID).
		SetChatID(chatID).
		SetAuthorID(authorID).
		SetIsVisibleForManager(visibleForManager).
		SetInitialRequestID(types.NewRequestID()).
		Save(s.Ctx)
	s.Require().NoError(err)
}
