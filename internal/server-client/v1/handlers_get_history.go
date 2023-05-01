package clientv1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/gerladeno/chat-service/internal/types"
)

var stub = MessagesPage{Messages: []Message{
	{
		AuthorId:  types.NewUserID(),
		Body:      "Здравствуйте! Разберёмся.",
		CreatedAt: time.Now(),
		Id:        types.NewMessageID(),
	},
	{
		AuthorId:  types.MustParse[types.UserID]("35f6dc0c-57c6-4066-aa93-75437138bed9"),
		Body:      "Привет! Не могу снять денег с карты,\nпишет 'карта заблокирована'",
		CreatedAt: time.Now().Add(-time.Minute),
		Id:        types.NewMessageID(),
	},
}}

func (h Handlers) PostGetHistory(eCtx echo.Context, params PostGetHistoryParams) error {
	if err := eCtx.JSON(http.StatusOK, GetHistoryResponse{Data: stub}); err != nil {
		return fmt.Errorf("sending stub response for requestId %s: %w", params.XRequestID, err)
	}
	return nil
}
