package eventstream

import (
	"go.uber.org/multierr"

	"github.com/gerladeno/chat-service/internal/types"
)

const (
	TypeMessageEventSent    = `MessageSentEvent`
	TypeMessageEventBlocked = `MessageBlockedEvent`
	TypeNewMessageEvent     = `NewMessageEvent`
	TypeNewChatEvent        = `NewChatEvent`
)

type Event interface {
	eventMarker()
	Validate() error
}

type event struct{}         //
func (*event) eventMarker() {}

type CoreEventFields struct {
	EventID   types.EventID
	EventType string
	RequestID types.RequestID
}

func (e CoreEventFields) Matches(x any) bool {
	val, ok := x.(CoreEventFields)
	if !ok {
		return false
	}
	return e.EventType == val.EventType && e.RequestID == val.RequestID
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
	return er
}
