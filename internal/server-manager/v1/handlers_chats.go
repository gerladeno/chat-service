package managerv1

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/gerladeno/chat-service/internal/middlewares"
	getchats "github.com/gerladeno/chat-service/internal/usecases/manager/get-chats"
)

func (h Handlers) PostGetChats(eCtx echo.Context, params PostGetChatsParams) error {
	ctx := eCtx.Request().Context()
	managerID := middlewares.MustUserID(eCtx)
	resp, err := h.getChatsUseCase.Handle(ctx, getchats.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
	})
	if err != nil {
		return fmt.Errorf("handling getChats: %v", err)
	}
	if err = eCtx.JSON(http.StatusOK, GetChatsResponse{
		Data: chats2chatsList(resp.Chats),
	}); err != nil {
		return fmt.Errorf("sending getChats response: %v", err)
	}
	return nil
}

func chats2chatsList(chats []getchats.Chat) *ChatList {
	result := ChatList{Chats: make([]Chat, 0, len(chats))}
	for _, chat := range chats {
		result.Chats = append(result.Chats, Chat{
			ChatId:   chat.ID,
			ClientId: chat.ClientID,
		})
	}
	return &result
}
