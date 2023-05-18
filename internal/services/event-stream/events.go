package eventstream

import (
	"time"

	"go.uber.org/multierr"

	"github.com/gerladeno/chat-service/internal/types"
)

const (
	typeMessageEventSent = `MessageSentEvent`
	typeNewMessageEvent  = `NewMessageEvent`
)

type Event interface {
	eventMarker()
	Validate() error
}

type event struct{}         //
func (*event) eventMarker() {}

// MessageSentEvent indicates that the message was checked by AFC
// and was sent to the manager. Two gray ticks.
type MessageSentEvent struct {
	event
	CoreEventFields
}

type CoreEventFields struct {
	EventID   types.EventID
	EventType string
	RequestID types.RequestID
	MessageID types.MessageID
}

func (e CoreEventFields) Matches(x any) bool {
	val, ok := x.(CoreEventFields)
	if !ok {
		return false
	}
	return e.EventType == val.EventType && e.MessageID == val.MessageID && e.RequestID == val.RequestID
}

func (e CoreEventFields) String() string {
	return e.EventType
}

func (e CoreEventFields) Validate() error {
	var er error
	if err := e.EventID.Validate(); err != nil {
		er = multierr.Append(er, err)
	}
	if err := e.RequestID.Validate(); err != nil {
		er = multierr.Append(er, err)
	}
	if err := e.MessageID.Validate(); err != nil {
		er = multierr.Append(er, err)
	}
	return er
}

// NewMessageEvent is a signal about the appearance of a new message in the chat.
type NewMessageEvent struct {
	event
	CoreEventFields
	ChatID      types.ChatID
	UserID      types.UserID
	CreatedAt   time.Time
	MessageBody string
	IsService   bool
}

func (e NewMessageEvent) Validate() error {
	er := e.CoreEventFields.Validate()
	if err := e.ChatID.Validate(); err != nil {
		er = multierr.Append(er, err)
	}
	if err := e.UserID.Validate(); err != nil {
		er = multierr.Append(er, err)
	}
	return er
}

func (e NewMessageEvent) Matches(x any) bool {
	val, ok := x.(*NewMessageEvent)
	if !ok {
		return false
	}
	return e.CoreEventFields.Matches(val.CoreEventFields) &&
		e.ChatID == val.ChatID &&
		e.UserID == val.UserID &&
		e.CreatedAt == val.CreatedAt &&
		e.MessageBody == val.MessageBody &&
		e.IsService == val.IsService
}

func NewNewMessageEvent(
	eventID types.EventID,
	requestID types.RequestID,
	chatID types.ChatID,
	messageID types.MessageID,
	userID types.UserID,
	createdAt time.Time,
	body string,
	isService bool,
) Event {
	return &NewMessageEvent{
		event: event{},
		CoreEventFields: CoreEventFields{
			EventID:   eventID,
			EventType: typeNewMessageEvent,
			RequestID: requestID,
			MessageID: messageID,
		},
		ChatID:      chatID,
		UserID:      userID,
		CreatedAt:   createdAt,
		MessageBody: body,
		IsService:   isService,
	}
}

func NewMessageSentEvent(
	eventID types.EventID,
	requestID types.RequestID,
	messageID types.MessageID,
) Event {
	return &MessageSentEvent{
		event: event{},
		CoreEventFields: CoreEventFields{
			EventID:   eventID,
			EventType: typeMessageEventSent,
			RequestID: requestID,
			MessageID: messageID,
		},
	}
}
