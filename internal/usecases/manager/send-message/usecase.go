package sendmessage

import (
	"context"
	"errors"
	"fmt"
	"time"

	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	problemsrepo "github.com/gerladeno/chat-service/internal/repositories/problems"
	sendmanagermessagejob "github.com/gerladeno/chat-service/internal/services/outbox/jobs/send-manager-message"
	"github.com/gerladeno/chat-service/internal/types"
)

var (
	ErrProblemNotFound = errors.New("problem not found")
	ErrInvalidRequest  = errors.New("invalid request")
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=sendmessagemocks

type messagesRepository interface {
	CreateFullVisible(
		ctx context.Context,
		reqID types.RequestID,
		problemID types.ProblemID,
		chatID types.ChatID,
		authorID types.UserID,
		msgBody string,
	) (*messagesrepo.Message, error)
}

type outboxService interface {
	Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
}

type problemsRepository interface {
	GetAssignedProblemID(ctx context.Context, managerID types.UserID, chatID types.ChatID) (types.ProblemID, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	messagesRepository messagesRepository `option:"mandatory" validate:"required"`
	outboxService      outboxService      `option:"mandatory" validate:"required"`
	problemsRepository problemsRepository `option:"mandatory" validate:"required"`
	transactor         transactor         `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	if err := opts.Validate(); err != nil {
		return UseCase{}, fmt.Errorf("validate send manager message usecase option validation: %v", err)
	}
	return UseCase{Options: opts}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, ErrInvalidRequest
	}
	problemID, err := u.problemsRepository.GetAssignedProblemID(ctx, req.ManagerID, req.ChatID)
	switch {
	case errors.Is(err, problemsrepo.ErrProblemNotFound):
		return Response{}, ErrProblemNotFound
	case err != nil:
		return Response{}, fmt.Errorf("get assigned problem: %v", err)
	}
	var msg *messagesrepo.Message
	if err = u.transactor.RunInTx(ctx, func(ctx context.Context) error {
		msg, err = u.messagesRepository.CreateFullVisible(ctx, req.ID, problemID, req.ChatID, req.ManagerID, req.MessageBody)
		if err != nil {
			return fmt.Errorf("create manager message: %v", err)
		}
		if _, err = u.outboxService.Put(
			ctx,
			sendmanagermessagejob.Name,
			sendmanagermessagejob.NewPayload(msg.ID, req.ManagerID),
			time.Now(),
		); err != nil {
			return fmt.Errorf("put outbox job to publish manager message: %v", err)
		}
		return nil
	}); err != nil {
		return Response{}, err
	}
	return Response{MessageID: msg.ID, CreatedAt: msg.CreatedAt}, nil
}
