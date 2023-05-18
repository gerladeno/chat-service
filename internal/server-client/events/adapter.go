package clientevents

import (
	"errors"

	eventstream "github.com/gerladeno/chat-service/internal/services/event-stream"
	"github.com/gerladeno/chat-service/internal/types"
	websocketstream "github.com/gerladeno/chat-service/internal/websocket-stream"
)

var _ websocketstream.EventAdapter = Adapter{}

var ErrUnsupportedEventType = errors.New("unsupported event type")

type Adapter struct{}

func (Adapter) Adapt(ev eventstream.Event) (any, error) {
	switch v := ev.(type) {
	case *eventstream.NewMessageEvent:
		var userID *types.UserID
		if !v.UserID.IsZero() {
			userID = &v.UserID
		}
		return NewMessageEvent{
			AuthorId:  userID,
			Body:      v.MessageBody,
			CreatedAt: v.CreatedAt,
			EventId:   v.EventID,
			EventType: v.EventType,
			IsService: v.IsService,
			MessageId: v.MessageID,
			RequestId: v.RequestID,
		}, nil
	case *eventstream.MessageSentEvent:
		return MessageEventSent{
			EventId:   v.EventID,
			EventType: v.EventType,
			MessageId: v.MessageID,
			RequestId: v.RequestID,
		}, nil
	}
	return nil, ErrUnsupportedEventType
}
