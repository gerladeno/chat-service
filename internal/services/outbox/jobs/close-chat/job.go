package closechatjob

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"

	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	eventstream "github.com/gerladeno/chat-service/internal/services/event-stream"
	"github.com/gerladeno/chat-service/internal/services/outbox"
	"github.com/gerladeno/chat-service/internal/types"
)

const (
	Name = `close-chat`
)

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=closechatjobmocks

type chatsRepo interface {
	GetClientID(ctx context.Context, chatID types.ChatID) (types.UserID, error)
}

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

type managerLoad interface {
	CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error)
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	chatsRepo         chatsRepo         `option:"mandatory" validate:"required"`
	messageRepository messageRepository `option:"mandatory" validate:"required"`
	eventStream       eventStream       `option:"mandatory" validate:"required"`
	managerLoad       managerLoad       `option:"mandatory" validate:"required"`
}

type Job struct {
	outbox.DefaultJob
	Options
}

func New(opts Options) (*Job, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validating send close chat job options: %v", err)
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
	requestID, chatID, managerID, messageID, err := idsFromPayload(payload)
	if err != nil {
		return err
	}
	msg, err := j.messageRepository.GetMessageByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("get message by id: %v", err)
	}
	clientID, err := j.chatsRepo.GetClientID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get client id: %v", err)
	}
	if err = j.eventStream.Publish(ctx, clientID, eventstream.NewNewMessageEvent(
		types.NewEventID(),
		msg.RequestID,
		chatID,
		msg.ID,
		types.UserIDNil,
		msg.CreatedAt,
		msg.Body,
		msg.IsService,
	)); err != nil {
		return fmt.Errorf("publish client service message: %v", err)
	}
	can, err := j.managerLoad.CanManagerTakeProblem(ctx, managerID)
	if err != nil {
		return fmt.Errorf("can manager take problems: %v", err)
	}
	if err = j.eventStream.Publish(ctx, managerID, eventstream.NewChatClosedEvent(
		types.NewEventID(),
		requestID,
		chatID,
		can,
	)); err != nil {
		return fmt.Errorf("publish manager close chat event: %v", err)
	}
	return nil
}

func idsFromPayload(payload string) (types.RequestID, types.ChatID, types.UserID, types.MessageID, error) {
	parts := strings.Split(payload, ".")
	if len(parts) != 4 {
		return types.RequestIDNil, types.ChatIDNil, types.UserIDNil, types.MessageIDNil,
			fmt.Errorf("err parsing payload: %s", payload)
	}
	requestID, err := types.Parse[types.RequestID](parts[0])
	if err != nil {
		return types.RequestIDNil, types.ChatIDNil, types.UserIDNil, types.MessageIDNil,
			fmt.Errorf("parsing requestID: %v", err)
	}
	chatID, err := types.Parse[types.ChatID](parts[1])
	if err != nil {
		return types.RequestIDNil, types.ChatIDNil, types.UserIDNil, types.MessageIDNil,
			fmt.Errorf("parsing chatID: %v", err)
	}
	managerID, err := types.Parse[types.UserID](parts[2])
	if err != nil {
		return types.RequestIDNil, types.ChatIDNil, types.UserIDNil, types.MessageIDNil,
			fmt.Errorf("parsing managerID: %v", err)
	}
	messageID, err := types.Parse[types.MessageID](parts[3])
	if err != nil {
		return types.RequestIDNil, types.ChatIDNil, types.UserIDNil, types.MessageIDNil,
			fmt.Errorf("parsing requestID: %v", err)
	}
	return requestID, chatID, managerID, messageID, nil
}

func NewPayload(
	requestID types.RequestID,
	chatID types.ChatID,
	managerID types.UserID,
	messageID types.MessageID,
) string {
	return fmt.Sprintf("%s.%s.%s.%s", requestID, chatID, managerID, messageID)
}
