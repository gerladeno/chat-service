package managerscheduler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	problemsrepo "github.com/gerladeno/chat-service/internal/repositories/problems"
	managerpool "github.com/gerladeno/chat-service/internal/services/manager-pool"
	managerassignedtoproblemjob "github.com/gerladeno/chat-service/internal/services/outbox/jobs/manager-assigned-to-problem"
	"github.com/gerladeno/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/service_mock.gen.go -package=managerschedulermocks

const serviceName = "manager-scheduler"

type problemsRepo interface {
	GetProblemsWithoutManager(ctx context.Context) ([]problemsrepo.Problem, error)
	AssignManager(ctx context.Context, problemID types.ProblemID, managerID types.UserID) error
	GetRequestID(ctx context.Context, problemID types.ProblemID) (types.RequestID, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

type msgRepo interface {
	CreateService(
		ctx context.Context,
		reqID types.RequestID,
		problemID types.ProblemID,
		chatID types.ChatID,
		msgBody string,
	) (*messagesrepo.Message, error)
}

type outbox interface {
	Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
}

//go:generate options-gen -out-filename=service_options.gen.go -from-struct=Options
type Options struct {
	period time.Duration `option:"mandatory" validate:"min=100ms,max=1m"`

	mngrPool     managerpool.Pool `option:"mandatory" validate:"required"`
	msgRepo      msgRepo          `option:"mandatory" validate:"required"`
	outbox       outbox           `option:"mandatory" validate:"required"`
	problemsRepo problemsRepo     `option:"mandatory" validate:"required"`
	tx           transactor       `option:"mandatory" validate:"required"`
}

type Service struct {
	Options
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate new managerscheduler options: %v", err)
	}
	return &Service{Options: opts}, nil
}

func (s *Service) Run(ctx context.Context) error {
	zap.L().Named(serviceName).Info("started manager scheduler")
	defer zap.L().Named(serviceName).Info("stopped")
	t := time.NewTicker(s.period)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			if err := s.performProblemsDistribution(ctx); err != nil && !errors.Is(err, context.Canceled) {
				return fmt.Errorf("problems distribution error: %v", err)
			}
		}
	}
}

func (s *Service) performProblemsDistribution(ctx context.Context) error {
	if err := s.tx.RunInTx(ctx, func(ctx context.Context) error {
		problems, err := s.problemsRepo.GetProblemsWithoutManager(ctx)
		if err != nil {
			return fmt.Errorf("get problems for distribution: %w", err)
		}
		for _, problem := range problems {
			managerID, err := s.mngrPool.Get(ctx)
			switch {
			case errors.Is(err, managerpool.ErrNoAvailableManagers):
				return nil
			case err != nil:
				return fmt.Errorf("get available manager: %w", err)
			}
			if err = s.problemsRepo.AssignManager(ctx, problem.ID, managerID); err != nil {
				return fmt.Errorf("assign manager to a problem: %v", err)
			}
			requestID, err := s.problemsRepo.GetRequestID(ctx, problem.ID)
			if err != nil {
				return fmt.Errorf("get problems request id: %w", err)
			}
			msg, err := s.msgRepo.CreateService(ctx, requestID, problem.ID, problem.ChatID, managerAssignedString(managerID))
			if err != nil {
				return fmt.Errorf("create service message: %w", err)
			}
			if _, err = s.outbox.Put(
				ctx,
				managerassignedtoproblemjob.Name,
				managerassignedtoproblemjob.NewPayload(msg.ID, managerID),
				time.Now(),
			); err != nil {
				return fmt.Errorf("put a job to outbox: %w", err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("assign manager in tx: %w", err)
	}
	return nil
}

func managerAssignedString(managerID types.UserID) string {
	return fmt.Sprintf(`Manager %s will answer you`, managerID)
}
