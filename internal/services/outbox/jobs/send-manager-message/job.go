package sendmanagermessagejob

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"

	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	eventstream "github.com/gerladeno/chat-service/internal/services/event-stream"
	msgproducer "github.com/gerladeno/chat-service/internal/services/msg-producer"
	"github.com/gerladeno/chat-service/internal/services/outbox"
	"github.com/gerladeno/chat-service/internal/types"
)

const Name = "send-manager-message"

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=sendmanagermessagejobmocks
type chatsRepository interface {
	GetUserID(ctx context.Context, chatID types.ChatID) (types.UserID, error)
}

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

type messageProducer interface {
	ProduceMessage(ctx context.Context, message msgproducer.Message) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	chatsRepository   chatsRepository   `option:"mandatory" validate:"required"`
	messageRepository messageRepository `option:"mandatory" validate:"required"`
	eventStream       eventStream       `option:"mandatory" validate:"required"`
	messageProducer   messageProducer   `option:"mandatory" validate:"required"`
}

type Job struct {
	outbox.DefaultJob
	Options
}

func New(opts Options) (*Job, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validating send manager message job options: %v", err)
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
	messageID, managerID, err := idsFromPayload(payload)
	if err != nil {
		return err
	}
	msg, err := j.messageRepository.GetMessageByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("getting message by id: %v", err)
	}
	clientID, err := j.chatsRepository.GetUserID(ctx, msg.ChatID)
	if err != nil {
		return fmt.Errorf("getting client id: %v", err)
	}
	if err = j.messageProducer.ProduceMessage(ctx, msgproducer.Message{
		ID:         messageID,
		ChatID:     msg.ChatID,
		Body:       msg.Body,
		FromClient: false,
	}); err != nil {
		return fmt.Errorf("produce message: %v", err)
	}
	if err = j.eventStream.Publish(ctx, clientID, eventstream.NewNewMessageEvent(
		types.NewEventID(),
		msg.RequestID,
		msg.ChatID,
		messageID,
		managerID,
		msg.CreatedAt,
		msg.Body,
		false,
	)); err != nil {
		return fmt.Errorf("publishing new new message event for client: %v", err)
	}
	if err = j.eventStream.Publish(ctx, managerID, eventstream.NewNewMessageEvent(
		types.NewEventID(),
		msg.RequestID,
		msg.ChatID,
		messageID,
		managerID,
		msg.CreatedAt,
		msg.Body,
		false,
	)); err != nil {
		return fmt.Errorf("publishing new new message event for client: %v", err)
	}
	if err = j.eventStream.Publish(ctx, managerID, eventstream.NewMessageSentEvent(
		types.NewEventID(),
		msg.RequestID,
		messageID,
	)); err != nil {
		return fmt.Errorf("publishing new message event sent: %v", err)
	}
	return nil
}

func idsFromPayload(payload string) (types.MessageID, types.UserID, error) {
	parts := strings.Split(payload, ".")
	if len(parts) != 2 {
		return types.MessageIDNil, types.UserIDNil, fmt.Errorf("err parsing payload: %s", payload)
	}
	msgID, err := types.Parse[types.MessageID](parts[0])
	if err != nil {
		return types.MessageIDNil, types.UserIDNil, fmt.Errorf("parsing messageID: %v", err)
	}
	managerID, err := types.Parse[types.UserID](parts[1])
	if err != nil {
		return types.MessageIDNil, types.UserIDNil, fmt.Errorf("parsing managerID: %v", err)
	}
	return msgID, managerID, nil
}

func NewPayload(messageID types.MessageID, managerID types.UserID) string {
	return fmt.Sprintf("%s.%s", messageID, managerID)
}
