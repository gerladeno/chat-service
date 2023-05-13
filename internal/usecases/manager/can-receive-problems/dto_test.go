package canreceiveproblems_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gerladeno/chat-service/internal/types"
	canreceiveproblems "github.com/gerladeno/chat-service/internal/usecases/manager/can-receive-problems"
)

func TestRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		id        types.RequestID
		managerID types.UserID
		err       bool
	}{
		{
			name:      "positive",
			id:        types.NewRequestID(),
			managerID: types.NewUserID(),
			err:       false,
		},
		{
			name:      "zero request id",
			id:        types.RequestIDNil,
			managerID: types.NewUserID(),
			err:       true,
		},
		{
			name:      "zero manager id",
			id:        types.NewRequestID(),
			managerID: types.UserIDNil,
			err:       true,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			r := canreceiveproblems.Request{
				ID:        test.id,
				ManagerID: test.managerID,
			}
			if test.err {
				require.Error(t, r.Validate())
			} else {
				require.NoError(t, r.Validate())
			}
		})
	}
}
