package managerv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	servererrors "github.com/gerladeno/chat-service/internal/errors"
	"github.com/gerladeno/chat-service/internal/middlewares"
	sendmessage "github.com/gerladeno/chat-service/internal/usecases/manager/send-message"
)

func (h Handlers) PostSendMessage(eCtx echo.Context, params PostSendMessageParams) error {
	ctx := eCtx.Request().Context()
	managerID := middlewares.MustUserID(eCtx)

	var req sendmessage.Request
	if err := eCtx.Bind(&req); err != nil {
		return fmt.Errorf("binding request for requestId %s: %w", params.XRequestID, err)
	}
	req.ID = params.XRequestID
	req.ManagerID = managerID
	resp, err := h.sendMessageUseCase.Handle(ctx, req)
	switch {
	case errors.Is(err, sendmessage.ErrInvalidRequest):
		return servererrors.NewServerError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err)
	case err != nil:
		return err
	}
	if err = eCtx.JSON(http.StatusOK, SendMessageResponse{
		Data: &MessageWithoutBody{
			AuthorId:  managerID,
			CreatedAt: resp.CreatedAt,
			Id:        resp.MessageID,
		},
	}); err != nil {
		return fmt.Errorf("sending response for requestId %s: %w", params.XRequestID, err)
	}
	return nil
}
