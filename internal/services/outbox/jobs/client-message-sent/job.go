package clientmessagesentjob

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	eventstream "github.com/gerladeno/chat-service/internal/services/event-stream"
	"github.com/gerladeno/chat-service/internal/services/outbox"
	"github.com/gerladeno/chat-service/internal/types"
)

const Name = "client-message-sent"

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=clientmessagesentjobmocks

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

type problemsRepository interface {
	GetActiveManager(ctx context.Context, chatID types.ChatID) (types.UserID, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	messageRepository  messageRepository  `option:"mandatory" validate:"required"`
	problemsRepository problemsRepository `option:"mandatory" validate:"required"`
	eventStream        eventStream        `option:"mandatory" validate:"required"`
}

type Job struct {
	outbox.DefaultJob
	Options
}

func New(opts Options) (*Job, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validating client message sent job options: %v", err)
	}
	return &Job{Options: opts}, nil
}

func (j *Job) Name() string {
	return Name
}

func (j *Job) Handle(ctx context.Context, payload string) (err error) {
	defer func() {
		if err != nil {
			zap.L().With(zap.String("payload", payload), zap.Error(err)).Debug("failed")
		} else {
			zap.L().With(zap.String("payload", payload)).Debug("success")
		}
	}()
	msgID, err := types.Parse[types.MessageID](payload)
	if err != nil {
		return fmt.Errorf("parsing messageID: %v", err)
	}
	msg, err := j.messageRepository.GetMessageByID(ctx, msgID)
	if err != nil {
		return fmt.Errorf("getting message by id: %v", err)
	}
	managerID, err := j.problemsRepository.GetActiveManager(ctx, msg.ChatID)
	if err != nil {
		return fmt.Errorf("get active manager by chat id: %v", err)
	}
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := j.eventStream.Publish(ctx, msg.AuthorID, eventstream.NewMessageSentEvent(
			types.NewEventID(),
			msg.RequestID,
			msg.ID,
		)); err != nil {
			return fmt.Errorf("client message sent publishing: %v", err)
		}
		return nil
	})
	if !managerID.IsZero() {
		eg.Go(func() error {
			if err := j.eventStream.Publish(ctx, managerID, eventstream.NewNewMessageEvent(
				types.NewEventID(),
				msg.RequestID,
				msg.ChatID,
				msg.ID,
				msg.AuthorID,
				msg.CreatedAt,
				msg.Body,
				false,
			)); err != nil {
				return fmt.Errorf("manager new message publishing: %v", err)
			}
			return nil
		})
	}
	return eg.Wait()
}
