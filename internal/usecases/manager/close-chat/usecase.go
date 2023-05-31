package closechat

import (
	"context"
	"errors"
	"fmt"
	"time"

	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	problemsrepo "github.com/gerladeno/chat-service/internal/repositories/problems"
	closechatjob "github.com/gerladeno/chat-service/internal/services/outbox/jobs/close-chat"
	"github.com/gerladeno/chat-service/internal/types"
)

const serviceMessageBody = `Your question has been marked as resolved.
Thank you for being with us!`

var ErrNoActiveProblem = errors.New("no active problem")

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mocks.gen.go -package=closechatmocks

type msgRepo interface {
	CreateService(
		ctx context.Context,
		reqID types.RequestID,
		problemID types.ProblemID,
		chatID types.ChatID,
		msgBody string,
	) (*messagesrepo.Message, error)
}

type problemsRepo interface {
	ResolveProblem(ctx context.Context, chatID types.ChatID, managerID types.UserID, requestID types.RequestID,
	) (types.ProblemID, error)
}

type outboxService interface {
	Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	msgRepo       msgRepo       `option:"mandatory" validate:"required"`
	problemsRepo  problemsRepo  `option:"mandatory" validate:"required"`
	outboxService outboxService `option:"mandatory" validate:"required"`
	transactor    transactor    `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	if err := opts.Validate(); err != nil {
		return UseCase{}, fmt.Errorf("validate close chat usecase options: %v", err)
	}
	return UseCase{Options: opts}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("request validation: %v", err)
	}
	err := u.transactor.RunInTx(ctx, func(ctx context.Context) error {
		problemID, err := u.problemsRepo.ResolveProblem(ctx, req.ChatID, req.ManagerID, req.ID)
		if err != nil {
			return fmt.Errorf("resolve problem: %w", err)
		}
		msg, err := u.msgRepo.CreateService(ctx, req.ID, problemID, req.ChatID, serviceMessageBody)
		if err != nil {
			return fmt.Errorf("create service message: %v", err)
		}
		if _, err = u.outboxService.Put(
			ctx,
			closechatjob.Name,
			closechatjob.NewPayload(req.ID, req.ChatID, req.ManagerID, msg.ID),
			time.Now(),
		); err != nil {
			return fmt.Errorf("create outbox job: %v", err)
		}
		return nil
	})
	switch {
	case errors.Is(err, problemsrepo.ErrProblemNotFound):
		return ErrNoActiveProblem
	case err != nil:
		return err
	}
	return nil
}
