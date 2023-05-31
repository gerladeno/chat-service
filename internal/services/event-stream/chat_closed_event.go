package eventstream

import (
	"go.uber.org/multierr"

	"github.com/gerladeno/chat-service/internal/types"
)

type ChatClosedEvent struct {
	event
	CoreEventFields
	ChatID              types.ChatID
	CanTakeMoreProblems bool
}

func (e ChatClosedEvent) Validate() error {
	er := e.CoreEventFields.Validate()
	if err := e.ChatID.Validate(); err != nil {
		er = multierr.Append(er, err)
	}
	return er
}

func (e ChatClosedEvent) Matches(x any) bool {
	val, ok := x.(*ChatClosedEvent)
	if !ok {
		return false
	}
	return e.CoreEventFields.Matches(val.CoreEventFields) &&
		e.ChatID == val.ChatID &&
		e.CanTakeMoreProblems == val.CanTakeMoreProblems
}

func NewChatClosedEvent(
	eventID types.EventID,
	requestID types.RequestID,
	chatID types.ChatID,
	canTakeMoreProblems bool,
) Event {
	return &ChatClosedEvent{
		event: event{},
		CoreEventFields: CoreEventFields{
			EventID:   eventID,
			EventType: TypeChatClosedEvent,
			RequestID: requestID,
		},
		ChatID:              chatID,
		CanTakeMoreProblems: canTakeMoreProblems,
	}
}
