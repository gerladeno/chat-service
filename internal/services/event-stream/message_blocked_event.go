package eventstream

import "github.com/gerladeno/chat-service/internal/types"

// MessageBlockedEvent indicates that the message was checked by AFC
// and was sent to the manager. Two gray ticks.
type MessageBlockedEvent struct {
	event
	CoreEventFields
}

func (e MessageBlockedEvent) Matches(x any) bool {
	val, ok := x.(*MessageBlockedEvent)
	if !ok {
		return false
	}
	return e.CoreEventFields.Matches(val.CoreEventFields)
}

func NewMessageBlockedEvent(
	eventID types.EventID,
	requestID types.RequestID,
	messageID types.MessageID,
) Event {
	return &MessageBlockedEvent{
		event: event{},
		CoreEventFields: CoreEventFields{
			EventID:   eventID,
			EventType: typeMessageEventBlocked,
			RequestID: requestID,
			MessageID: messageID,
		},
	}
}
