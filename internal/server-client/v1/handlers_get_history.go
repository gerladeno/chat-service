package clientv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	servererrors "github.com/gerladeno/chat-service/internal/errors"
	"github.com/gerladeno/chat-service/internal/middlewares"
	gethistory "github.com/gerladeno/chat-service/internal/usecases/client/get-history"
)

func (h Handlers) PostGetHistory(eCtx echo.Context, params PostGetHistoryParams) error {
	ctx := eCtx.Request().Context()
	clientID := middlewares.MustUserID(eCtx)
	var req gethistory.Request
	if err := eCtx.Bind(&req); err != nil {
		return fmt.Errorf("binding request for requestId %s: %w", params.XRequestID, err)
	}
	req.ID = params.XRequestID
	req.ClientID = clientID
	resp, err := h.getHistoryUseCase.Handle(ctx, req)
	switch {
	case errors.Is(err, gethistory.ErrInvalidCursor) || errors.Is(err, gethistory.ErrInvalidRequest):
		return servererrors.NewServerError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err)
	case err != nil:
		return err
	}
	if err = eCtx.JSON(http.StatusOK, GetHistoryResponse{Data: getHistory2MessagesPage(resp)}); err != nil {
		return fmt.Errorf("sending response for requestId %s: %w", params.XRequestID, err)
	}
	return nil
}

func getHistory2MessagesPage(resp gethistory.Response) *MessagesPage {
	mp := MessagesPage{Next: resp.NextCursor}
	mp.Messages = make([]Message, 0, len(resp.Messages))
	var tmp Message
	for i := range resp.Messages {
		tmp = Message{
			Body:       resp.Messages[i].Body,
			CreatedAt:  resp.Messages[i].CreatedAt,
			Id:         resp.Messages[i].ID,
			IsBlocked:  resp.Messages[i].IsBlocked,
			IsReceived: resp.Messages[i].IsReceived,
			IsService:  resp.Messages[i].IsService,
		}
		if !resp.Messages[i].AuthorID.IsZero() {
			tmp.AuthorId = &resp.Messages[i].AuthorID
		}
		mp.Messages = append(mp.Messages, tmp)
	}
	return &mp
}
