package closechat_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gerladeno/chat-service/internal/types"
	closechat "github.com/gerladeno/chat-service/internal/usecases/manager/close-chat"
)

func TestRequest_Validate(t *testing.T) {
	cases := []struct {
		name    string
		request closechat.Request
		wantErr bool
	}{
		// Positive.
		{
			name: "valid request",
			request: closechat.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.NewUserID(),
				ChatID:    types.NewChatID(),
			},
			wantErr: false,
		},

		// Negative.
		{
			name: "require request id",
			request: closechat.Request{
				ID:        types.RequestIDNil,
				ManagerID: types.NewUserID(),
				ChatID:    types.NewChatID(),
			},
			wantErr: true,
		},
		{
			name: "require manager id",
			request: closechat.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.UserIDNil,
				ChatID:    types.NewChatID(),
			},
			wantErr: true,
		},
		{
			name: "require chat id",
			request: closechat.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.NewUserID(),
				ChatID:    types.ChatIDNil,
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
