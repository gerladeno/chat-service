//go:build integration

package chats_test

import (
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"

	chatsrepo "github.com/gerladeno/chat-service/internal/repositories/chats"
	"github.com/gerladeno/chat-service/internal/testingh"
	"github.com/gerladeno/chat-service/internal/types"
)

type ChatsRepoSuite struct {
	testingh.DBSuite
	repo *chatsrepo.Repo
}

func TestChatsRepoSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ChatsRepoSuite{DBSuite: testingh.NewDBSuite("TestChatsRepoSuite")})
}

func (s *ChatsRepoSuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	var err error

	s.repo, err = chatsrepo.New(chatsrepo.NewOptions(s.Database))
	s.Require().NoError(err)
}

func (s *ChatsRepoSuite) Test_CreateIfNotExists() {
	s.Run("chat does not exist, should be created", func() {
		clientID := types.NewUserID()

		chatID, err := s.repo.CreateIfNotExists(s.Ctx, clientID)
		s.Require().NoError(err)
		s.NotEmpty(chatID)
	})

	s.Run("chat already exists", func() {
		clientID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		chatID, err := s.repo.CreateIfNotExists(s.Ctx, clientID)
		s.Require().NoError(err)
		s.Require().NotEmpty(chatID)
		s.Equal(chat.ID, chatID)
	})
}

func (s *ChatsRepoSuite) Test_GetUserID() {
	s.Run("chat doesn't exist", func() {
		_, err := s.repo.GetUserID(s.Ctx, types.NewChatID())
		s.Require().ErrorIs(err, chatsrepo.ErrChatNotFound)
	})

	s.Run("create chat and find", func() {
		clientID := types.NewUserID()

		chatID, err := s.repo.CreateIfNotExists(s.Ctx, clientID)
		s.Require().NoError(err)
		s.NotEmpty(chatID)

		userID, err := s.repo.GetUserID(s.Ctx, chatID)
		s.Require().NoError(err)
		s.Require().Equal(clientID, userID)
	})
}

func (s *ChatsRepoSuite) Test_GetChatsForManager() {
	s.Run("no chats", func() {
		chats, err := s.repo.GetChatsForManager(s.Ctx, types.NewUserID())
		s.Require().NoError(err)
		s.Require().Empty(chats)
	})

	s.Run("no assigned problems", func() {
		_, _ = s.createChatAndProblem(types.NewUserID())
		chats, err := s.repo.GetChatsForManager(s.Ctx, types.NewUserID())
		s.Require().NoError(err)
		s.Require().Empty(chats)
	})

	s.Run("problems assigned to other manager", func() {
		_, problemID := s.createChatAndProblem(types.NewUserID())
		_, err := s.Database.Problem(s.Ctx).UpdateOneID(problemID).SetManagerID(types.NewUserID()).Save(s.Ctx)
		s.Require().NoError(err)
		chats, err := s.repo.GetChatsForManager(s.Ctx, types.NewUserID())
		s.Require().NoError(err)
		s.Require().Empty(chats)
	})

	s.Run("have assigned problem", func() {
		managerID := types.NewUserID()
		_, problemID := s.createChatAndProblem(types.NewUserID())
		_, err := s.Database.Problem(s.Ctx).UpdateOneID(problemID).SetManagerID(managerID).Save(s.Ctx)
		s.Require().NoError(err)
		chats, err := s.repo.GetChatsForManager(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Require().Len(chats, 1)
	})

	s.Run("problem assigned but resolved", func() {
		managerID := types.NewUserID()
		_, problemID := s.createChatAndProblem(types.NewUserID())
		_, err := s.Database.Problem(s.Ctx).UpdateOneID(problemID).
			SetManagerID(managerID).
			SetResolvedAt(time.Now()).
			Save(s.Ctx)
		s.Require().NoError(err)
		chats, err := s.repo.GetChatsForManager(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Require().Empty(chats)
	})
}

func (s *ChatsRepoSuite) createChatAndProblem(clientID types.UserID) (types.ChatID, types.ProblemID) {
	s.T().Helper()
	chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
	s.Require().NoError(err)
	problem, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).Save(s.Ctx)
	s.Require().NoError(err)
	return chat.ID, problem.ID
}
