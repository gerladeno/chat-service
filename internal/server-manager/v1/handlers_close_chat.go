package managerv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	servererrors "github.com/gerladeno/chat-service/internal/errors"
	"github.com/gerladeno/chat-service/internal/middlewares"
	closechat "github.com/gerladeno/chat-service/internal/usecases/manager/close-chat"
)

const (
	ErrorCodeNoActiveProblem = 5001
	NoActiveProblemError     = `no active problem in the chat`
)

func (h Handlers) PostCloseChat(eCtx echo.Context, params PostCloseChatParams) error {
	ctx := eCtx.Request().Context()
	reqID := params.XRequestID
	managerID := middlewares.MustUserID(eCtx)
	var req closechat.Request
	if err := eCtx.Bind(&req); err != nil {
		return fmt.Errorf("bind request: %v", err)
	}
	req.ID = reqID
	req.ManagerID = managerID
	err := h.closeChatUseCase.Handle(ctx, req)
	switch {
	case errors.Is(err, closechat.ErrNoActiveProblem):
		return servererrors.NewServerError(ErrorCodeNoActiveProblem, NoActiveProblemError, err)
	case err != nil:
		return fmt.Errorf("handle close chat request: %v", err)
	}
	if err = eCtx.JSON(http.StatusOK, CloseChatResponse{}); err != nil {
		return fmt.Errorf("send response: %v", err)
	}
	return nil
}
