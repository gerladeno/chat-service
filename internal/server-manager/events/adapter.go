package managerevents

import (
	"errors"

	eventstream "github.com/gerladeno/chat-service/internal/services/event-stream"
	"github.com/gerladeno/chat-service/internal/types"
	websocketstream "github.com/gerladeno/chat-service/internal/websocket-stream"
)

var _ websocketstream.EventAdapter = Adapter{}

type Adapter struct{}

var ErrUnsupportedEventType = errors.New("unsupported event type")

func (Adapter) Adapt(ev eventstream.Event) (any, error) {
	var event Event
	var err error
	switch v := ev.(type) {
	case *eventstream.NewMessageEvent:
		event.EventId = v.EventID
		event.RequestId = v.RequestID
		var userID *types.UserID
		if !v.AuthorID.IsZero() {
			userID = &v.AuthorID
		}
		err = event.FromNewMessageEvent(NewMessageEvent{
			AuthorId:  userID,
			Body:      v.MessageBody,
			CreatedAt: v.CreatedAt,
			IsService: v.IsService,
			MessageId: v.MessageID,
		})
	case *eventstream.MessageSentEvent:
		event.EventId = v.EventID
		event.RequestId = v.RequestID
		err = event.FromMessageSentEvent(MessageSentEvent{
			MessageId: v.MessageID,
		})
	case *eventstream.NewChatEvent:
		event.EventId = v.EventID
		event.RequestId = v.RequestID
		err = event.FromNewChatEvent(NewChatEvent{
			CanTakeMoreProblems: v.CanTakeMoreProblems,
			ChatId:              v.ChatID,
			ClientId:            v.UserID,
		})
	default:
		return nil, ErrUnsupportedEventType
	}
	if err != nil {
		return nil, err
	}
	return event, nil
}
