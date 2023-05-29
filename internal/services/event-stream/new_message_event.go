package eventstream

import (
	"errors"
	"time"

	"go.uber.org/multierr"

	"github.com/gerladeno/chat-service/internal/types"
)

// NewMessageEvent is a signal about the appearance of a new message in the chat.
type NewMessageEvent struct {
	event
	CoreEventFields
	MessageID   types.MessageID
	ChatID      types.ChatID
	AuthorID    types.UserID
	CreatedAt   time.Time
	MessageBody string
	IsService   bool
}

func (e NewMessageEvent) Validate() error {
	er := e.CoreEventFields.Validate()
	if err := e.ChatID.Validate(); err != nil {
		er = multierr.Append(er, err)
	}
	if err := e.AuthorID.Validate(); err != nil && !errors.Is(err, types.ErrZeroID) {
		er = multierr.Append(er, err)
	}
	if err := e.MessageID.Validate(); err != nil {
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
		e.AuthorID == val.AuthorID &&
		e.CreatedAt == val.CreatedAt &&
		e.MessageBody == val.MessageBody &&
		e.IsService == val.IsService &&
		e.MessageID == val.MessageID
}

func NewNewMessageEvent(
	eventID types.EventID,
	requestID types.RequestID,
	chatID types.ChatID,
	messageID types.MessageID,
	authorID types.UserID,
	createdAt time.Time,
	body string,
	isService bool,
) Event {
	return &NewMessageEvent{
		event: event{},
		CoreEventFields: CoreEventFields{
			EventID:   eventID,
			EventType: TypeNewMessageEvent,
			RequestID: requestID,
		},
		ChatID:      chatID,
		AuthorID:    authorID,
		CreatedAt:   createdAt,
		MessageBody: body,
		IsService:   isService,
		MessageID:   messageID,
	}
}
