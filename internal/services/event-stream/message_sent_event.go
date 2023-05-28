package eventstream

import (
	"go.uber.org/multierr"

	"github.com/gerladeno/chat-service/internal/types"
)

// MessageSentEvent indicates that the message was checked by AFC
// and was sent to the manager. Two gray ticks.
type MessageSentEvent struct {
	event
	CoreEventFields
	MessageID types.MessageID
}

func (e MessageSentEvent) Validate() error {
	er := e.CoreEventFields.Validate()
	if err := e.MessageID.Validate(); err != nil {
		er = multierr.Append(er, err)
	}
	return er
}

func (e MessageSentEvent) Matches(x any) bool {
	val, ok := x.(*MessageSentEvent)
	if !ok {
		return false
	}
	return e.CoreEventFields.Matches(val.CoreEventFields) && e.MessageID == val.MessageID
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
			EventType: TypeMessageEventSent,
			RequestID: requestID,
		},
		MessageID: messageID,
	}
}
