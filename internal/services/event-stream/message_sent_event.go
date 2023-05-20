package eventstream

import "github.com/gerladeno/chat-service/internal/types"

// MessageSentEvent indicates that the message was checked by AFC
// and was sent to the manager. Two gray ticks.
type MessageSentEvent struct {
	event
	CoreEventFields
}

func (e MessageSentEvent) Matches(x any) bool {
	val, ok := x.(*MessageSentEvent)
	if !ok {
		return false
	}
	return e.CoreEventFields.Matches(val.CoreEventFields)
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
