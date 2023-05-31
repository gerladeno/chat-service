package managerv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	servererrors "github.com/gerladeno/chat-service/internal/errors"
	"github.com/gerladeno/chat-service/internal/middlewares"
	gethistory "github.com/gerladeno/chat-service/internal/usecases/client/get-history"
	getchathistory "github.com/gerladeno/chat-service/internal/usecases/manager/get-chat-history"
)

func (h Handlers) PostGetChatHistory(eCtx echo.Context, params PostGetChatHistoryParams) error {
	ctx := eCtx.Request().Context()
	managerID := middlewares.MustUserID(eCtx)
	var req getchathistory.Request
	if err := eCtx.Bind(&req); err != nil {
		return fmt.Errorf("binding request for requestId %s: %w", params.XRequestID, err)
	}
	req.ManagerID = managerID
	resp, err := h.getChatHistoryUseCase.Handle(ctx, req)
	switch {
	case errors.Is(err, gethistory.ErrInvalidCursor) || errors.Is(err, gethistory.ErrInvalidRequest):
		return servererrors.NewServerError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err)
	case err != nil:
		return err
	}
	if err = eCtx.JSON(http.StatusOK, GetChatHistoryResponse{Data: getChatHistory2MessagesPage(resp)}); err != nil {
		return fmt.Errorf("sending response for requestId %s: %w", params.XRequestID, err)
	}
	return nil
}

func getChatHistory2MessagesPage(resp getchathistory.Response) *MessagesPage {
	mp := MessagesPage{Next: resp.NextCursor}
	mp.Messages = make([]Message, 0, len(resp.Messages))
	for i := range resp.Messages {
		mp.Messages = append(mp.Messages, Message{
			Body:      resp.Messages[i].Body,
			AuthorId:  resp.Messages[i].AuthorID,
			CreatedAt: resp.Messages[i].CreatedAt,
			Id:        resp.Messages[i].ID,
		})
	}
	return &mp
}
