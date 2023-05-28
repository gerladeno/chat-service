//go:build integration

package problems_test

import (
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"

	problemsrepo "github.com/gerladeno/chat-service/internal/repositories/problems"
	storeproblem "github.com/gerladeno/chat-service/internal/store/problem"
	"github.com/gerladeno/chat-service/internal/testingh"
	"github.com/gerladeno/chat-service/internal/types"
)

type ProblemsRepoSuite struct {
	testingh.DBSuite
	repo *problemsrepo.Repo
}

func TestProblemsRepoSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ProblemsRepoSuite{DBSuite: testingh.NewDBSuite("TestProblemsRepoSuite")})
}

func (s *ProblemsRepoSuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	var err error

	s.repo, err = problemsrepo.New(problemsrepo.NewOptions(s.Database))
	s.Require().NoError(err)
}

func (s *ProblemsRepoSuite) Test_CreateIfNotExists() {
	s.Run("problem does not exist, should be created", func() {
		clientID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		problemID, err := s.repo.CreateIfNotExists(s.Ctx, chat.ID)
		s.Require().NoError(err)
		s.NotEmpty(problemID)

		problem, err := s.Database.Problem(s.Ctx).Get(s.Ctx, problemID)
		s.Require().NoError(err)
		s.Equal(problemID, problem.ID)
		s.Equal(chat.ID, problem.ChatID)
	})

	s.Run("resolved problem already exists, should be created", func() {
		clientID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		// Create problem.
		problem, err := s.Database.Problem(s.Ctx).Create().
			SetChatID(chat.ID).
			SetManagerID(types.NewUserID()).
			SetResolvedAt(time.Now()).Save(s.Ctx)
		s.Require().NoError(err)

		problemID, err := s.repo.CreateIfNotExists(s.Ctx, chat.ID)
		s.Require().NoError(err)
		s.NotEmpty(problemID)
		s.NotEqual(problem.ID, problemID)
	})

	s.Run("problem already exists", func() {
		clientID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		// Create problem.
		problem, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).Save(s.Ctx)
		s.Require().NoError(err)

		problemID, err := s.repo.CreateIfNotExists(s.Ctx, chat.ID)
		s.Require().NoError(err)
		s.NotEmpty(problemID)
		s.Equal(problem.ID, problemID)
	})
}

func (s *ProblemsRepoSuite) Test_GetManagerOpenProblemsCount() {
	s.Run("manager has no open problems", func() {
		managerID := types.NewUserID()

		count, err := s.repo.GetManagerOpenProblemsCount(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Empty(count)
	})

	s.Run("manager has open problems", func() {
		const (
			problemsCount         = 20
			resolvedProblemsCount = 3
		)

		managerID := types.NewUserID()
		problems := make([]types.ProblemID, 0, problemsCount)

		for i := 0; i < problemsCount; i++ {
			_, pID := s.createChatWithProblemAssignedTo(managerID)
			problems = append(problems, pID)
		}

		// Create problems for other managers.
		for i := 0; i < problemsCount; i++ {
			s.createChatWithProblemAssignedTo(types.NewUserID())
		}

		count, err := s.repo.GetManagerOpenProblemsCount(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Equal(problemsCount, count)

		// Resolve some problems.
		for i := 0; i < resolvedProblemsCount; i++ {
			pID := problems[i*resolvedProblemsCount]
			_, err := s.Database.Problem(s.Ctx).
				Update().
				Where(storeproblem.ID(pID)).
				SetResolvedAt(time.Now()).
				Save(s.Ctx)
			s.Require().NoError(err)
		}

		count, err = s.repo.GetManagerOpenProblemsCount(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Equal(problemsCount-resolvedProblemsCount, count)
	})
}

func (s *ProblemsRepoSuite) createChatWithProblemAssignedTo(managerID types.UserID) (types.ChatID, types.ProblemID) {
	s.T().Helper()

	// 1 chat can have only 1 open problem.

	chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(types.NewUserID()).Save(s.Ctx)
	s.Require().NoError(err)

	p, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).SetManagerID(managerID).Save(s.Ctx)
	s.Require().NoError(err)

	return chat.ID, p.ID
}

func (s *ProblemsRepoSuite) Test_GetProblemsWithoutManager() {
	s.Run("No problems", func() {
		problems, err := s.repo.GetProblemsWithoutManager(s.Ctx)
		s.Require().NoError(err)
		s.Require().Empty(problems)
	})

	s.Run("Unassigned problem exists, but no visible messages", func() {
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(types.NewUserID()).Save(s.Ctx)
		s.Require().NoError(err)

		_, err = s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).Save(s.Ctx)
		s.Require().NoError(err)

		problems, err := s.repo.GetProblemsWithoutManager(s.Ctx)
		s.Require().NoError(err)
		s.Require().Empty(problems)
	})

	s.Run("Unassigned problem exists", func() {
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(types.NewUserID()).Save(s.Ctx)
		s.Require().NoError(err)

		p, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).Save(s.Ctx)
		s.Require().NoError(err)

		_, err = s.Database.Message(s.Ctx).Create().
			SetProblemID(p.ID).
			SetID(types.NewMessageID()).
			SetChatID(chat.ID).
			SetIsVisibleForManager(true).
			SetInitialRequestID(types.NewRequestID()).
			SetBody("biba").
			Save(s.Ctx)
		s.Require().NoError(err)

		problems, err := s.repo.GetProblemsWithoutManager(s.Ctx)
		s.Require().NoError(err)
		s.Require().Len(problems, 1)
	})
}

func (s *ProblemsRepoSuite) Test_AssignManager() {
	s.Run("Problem not found", func() {
		err := s.repo.AssignManager(s.Ctx, types.NewProblemID(), types.NewUserID())
		s.Require().ErrorIs(err, problemsrepo.ErrProblemNotFound)
	})

	s.Run("assign manager", func() {
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(types.NewUserID()).Save(s.Ctx)
		s.Require().NoError(err)

		p, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).Save(s.Ctx)
		s.Require().NoError(err)

		managerID := types.NewUserID()
		err = s.repo.AssignManager(s.Ctx, p.ID, managerID)
		s.Require().NoError(err)

		p = s.Database.Problem(s.Ctx).GetX(s.Ctx, p.ID)
		s.Require().Equal(managerID, p.ManagerID)
	})
}

func (s *ProblemsRepoSuite) Test_GetRequestID() {
	s.Run("problem not found", func() {
		_, err := s.repo.GetRequestID(s.Ctx, types.NewProblemID())
		s.Require().ErrorIs(err, problemsrepo.ErrProblemNotFound)
	})

	s.Run("get request id", func() {
		reqID := types.NewRequestID()
		chatID, pID := s.createChatWithProblemAssignedTo(types.NewUserID())
		_, err := s.Database.Message(s.Ctx).Create().
			SetProblemID(pID).
			SetID(types.NewMessageID()).
			SetChatID(chatID).
			SetIsVisibleForManager(true).
			SetInitialRequestID(reqID).
			SetBody("biba").
			Save(s.Ctx)
		s.Require().NoError(err)

		_, err = s.Database.Message(s.Ctx).Create().
			SetProblemID(pID).
			SetID(types.NewMessageID()).
			SetChatID(chatID).
			SetIsVisibleForManager(true).
			SetInitialRequestID(types.NewRequestID()).
			SetBody("boba").
			Save(s.Ctx)
		s.Require().NoError(err)

		req, err := s.repo.GetRequestID(s.Ctx, pID)
		s.Require().NoError(err)
		s.Require().Equal(reqID, req)
	})
}
