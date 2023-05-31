package managerassignedtoproblemjob

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

const Name = "manager-assigned-to-problem"

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=managerassignedtoproblemjobmocks

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

type chatRepository interface {
	GetUserID(ctx context.Context, chatID types.ChatID) (types.UserID, error)
}

type managerLoad interface {
	CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	messageRepository messageRepository `option:"mandatory" validate:"required"`
	chatRepository    chatRepository    `option:"mandatory" validate:"required"`
	managerLoad       managerLoad       `option:"mandatory" validate:"required"`
	eventStream       eventStream       `option:"mandatory" validate:"required"`
}

type Job struct {
	outbox.DefaultJob
	Options
}

func New(opts Options) (*Job, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validating manager assigned job options: %v", err)
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
	userID, err := j.chatRepository.GetUserID(ctx, msg.ChatID)
	if err != nil {
		return fmt.Errorf("getting client id: %v", err)
	}
	if err = j.eventStream.Publish(ctx, userID, eventstream.NewNewMessageEvent(
		types.NewEventID(),
		msg.RequestID,
		msg.ChatID,
		msg.ID,
		msg.AuthorID,
		msg.CreatedAt,
		msg.Body,
		msg.IsService,
	)); err != nil {
		return fmt.Errorf("publishing message: %v", err)
	}
	can, err := j.managerLoad.CanManagerTakeProblem(ctx, managerID)
	if err != nil {
		return fmt.Errorf("can manager take problem for newChatEvent: %v", err)
	}
	if err = j.eventStream.Publish(ctx, managerID, eventstream.NewNewChatEvent(
		types.NewEventID(),
		msg.RequestID,
		msg.ChatID,
		userID,
		can,
	)); err != nil {
		return fmt.Errorf("publishing message: %v", err)
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
		return types.MessageIDNil, types.UserIDNil, fmt.Errorf("parsing messageID: %v", err)
	}
	return msgID, managerID, nil
}

func NewPayload(messageID types.MessageID, managerID types.UserID) string {
	return fmt.Sprintf("%s.%s", messageID, managerID)
}
