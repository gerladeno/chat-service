package managerv1_test

import (
	"errors"
	"fmt"
	"net/http"

	internalerrors "github.com/gerladeno/chat-service/internal/errors"
	managerv1 "github.com/gerladeno/chat-service/internal/server-manager/v1"
	"github.com/gerladeno/chat-service/internal/types"
	closechat "github.com/gerladeno/chat-service/internal/usecases/manager/close-chat"
)

func (s *HandlersSuite) TestCloseChat_Usecase_SomeError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()

	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat",
		fmt.Sprintf(`{"chatId": %q}`, chatID))

	s.closeChatUseCase.EXPECT().Handle(eCtx.Request().Context(), closechat.Request{
		ID:        reqID,
		ManagerID: s.managerID,
		ChatID:    chatID,
	}).Return(errors.New("gigabooo"))

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestCloseChat_Usecase_NoProblemError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()

	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat",
		fmt.Sprintf(`{"chatId": %q}`, chatID))

	s.closeChatUseCase.EXPECT().Handle(eCtx.Request().Context(), closechat.Request{
		ID:        reqID,
		ManagerID: s.managerID,
		ChatID:    chatID,
	}).Return(closechat.ErrNoActiveProblem)

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Empty(resp.Body)
	s.Require().Equal(managerv1.ErrorCodeNoActiveProblem, internalerrors.GetServerErrorCode(err))
}

func (s *HandlersSuite) TestCloseChat_Usecase_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()

	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat",
		fmt.Sprintf(`{"chatId": %q}`, chatID))

	s.closeChatUseCase.EXPECT().Handle(eCtx.Request().Context(), closechat.Request{
		ID:        reqID,
		ManagerID: s.managerID,
		ChatID:    chatID,
	}).Return(nil)

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.Code)
	s.JSONEq(`
{
    "data": null
}`, resp.Body.String())
}
