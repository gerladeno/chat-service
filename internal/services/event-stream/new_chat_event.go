package eventstream

import (
	"go.uber.org/multierr"

	"github.com/gerladeno/chat-service/internal/types"
)

type NewChatEvent struct {
	event
	CoreEventFields
	ChatID              types.ChatID
	UserID              types.UserID
	CanTakeMoreProblems bool
}

func (e NewChatEvent) Validate() error {
	er := e.CoreEventFields.Validate()
	if err := e.ChatID.Validate(); err != nil {
		er = multierr.Append(er, err)
	}
	if err := e.UserID.Validate(); err != nil {
		er = multierr.Append(er, err)
	}
	return er
}

func (e NewChatEvent) Matches(x any) bool {
	val, ok := x.(*NewChatEvent)
	if !ok {
		return false
	}
	return e.CoreEventFields.Matches(val.CoreEventFields) &&
		e.ChatID == val.ChatID &&
		e.UserID == val.UserID &&
		e.CanTakeMoreProblems == val.CanTakeMoreProblems
}

func NewNewChatEvent(
	eventID types.EventID,
	requestID types.RequestID,
	chatID types.ChatID,
	userID types.UserID,
	canTakeMoreProblems bool,
) Event {
	return &NewChatEvent{
		event: event{},
		CoreEventFields: CoreEventFields{
			EventID:   eventID,
			EventType: TypeNewChatEvent,
			RequestID: requestID,
		},
		ChatID:              chatID,
		UserID:              userID,
		CanTakeMoreProblems: canTakeMoreProblems,
	}
}
