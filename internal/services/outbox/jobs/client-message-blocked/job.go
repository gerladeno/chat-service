package clientmessageblockedjob

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	eventstream "github.com/gerladeno/chat-service/internal/services/event-stream"
	"github.com/gerladeno/chat-service/internal/services/outbox"
	"github.com/gerladeno/chat-service/internal/types"
)

const Name = "client-message-blocked"

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package clientmessageblockedjobmocks

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	messageRepository messageRepository `option:"mandatory"`
	eventStream       eventStream       `option:"mandatory"`
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
	if err = j.eventStream.Publish(ctx, msg.AuthorID, eventstream.NewMessageBlockedEvent(
		types.NewEventID(),
		msg.RequestID,
		msg.ID,
	)); err != nil {
		return fmt.Errorf("publishing message: %v", err)
	}
	return nil
}
