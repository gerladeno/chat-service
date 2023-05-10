package sendclientmessagejob

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	msgproducer "github.com/gerladeno/chat-service/internal/services/msg-producer"
	"github.com/gerladeno/chat-service/internal/services/outbox"
	"github.com/gerladeno/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=sendclientmessagejobmocks

const Name = "send-client-message"

type messageProducer interface {
	ProduceMessage(ctx context.Context, message msgproducer.Message) error
}

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	messageProducer   messageProducer   `option:"mandatory"`
	messageRepository messageRepository `option:"mandatory"`
}

type Job struct {
	outbox.DefaultJob
	Options
}

func New(opts Options) (*Job, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validating send client message job options: %v", err)
	}
	return &Job{Options: opts}, nil
}

func (j *Job) Name() string {
	return Name
}

func (j *Job) Handle(ctx context.Context, payload string) (err error) {
	defer func() {
		if err != nil {
			zap.L().With(zap.String("payload", payload), zap.Error(err)).Info("failed")
		} else {
			zap.L().With(zap.String("payload", payload)).Info("success")
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
	if err = j.messageProducer.ProduceMessage(ctx, msgproducer.Message{
		ID:         msg.ID,
		ChatID:     msg.ChatID,
		Body:       msg.Body,
		FromClient: !msg.AuthorID.IsZero(),
	}); err != nil {
		return fmt.Errorf("producing message: %v", err)
	}
	return nil
}
