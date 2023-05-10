package sendclientmessagejob

import (
	"github.com/gerladeno/chat-service/internal/types"
)

func MarshalPayload(messageID types.MessageID) (string, error) {
	if messageID.IsZero() {
		return "", types.ErrEntityIsNil
	}
	return messageID.String(), nil
}
