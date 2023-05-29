package managerv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	servererrors "github.com/gerladeno/chat-service/internal/errors"
	"github.com/gerladeno/chat-service/internal/middlewares"
	canreceiveproblems "github.com/gerladeno/chat-service/internal/usecases/manager/can-receive-problems"
	freehands "github.com/gerladeno/chat-service/internal/usecases/manager/free-hands"
)

const (
	ErrorCodeManagerOverloaded = 5000
	ManagerOverloadedError     = `manager overloaded`
)

func (h Handlers) PostGetFreeHandsBtnAvailability(eCtx echo.Context, params PostGetFreeHandsBtnAvailabilityParams) error {
	ctx := eCtx.Request().Context()
	managerID := middlewares.MustUserID(eCtx)
	req := canreceiveproblems.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
	}
	can, err := h.canReceiveProblemsUseCase.Handle(ctx, req)
	if err != nil {
		return fmt.Errorf("handling can receive problem button availability: %v", err)
	}
	if err = eCtx.JSON(http.StatusOK, GetFreeHandsBtnAvailabilityResponse{
		Data: &FreeHandsBtnAvailability{Available: can.Result},
	}); err != nil {
		return fmt.Errorf("sending response on can receive problem button availability: %v", err)
	}
	return nil
}

func (h Handlers) PostFreeHands(eCtx echo.Context, params PostFreeHandsParams) error {
	ctx := eCtx.Request().Context()
	managerID := middlewares.MustUserID(eCtx)
	req := freehands.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
	}
	err := h.freeHandsUseCase.Handle(ctx, req)
	switch {
	case errors.Is(err, freehands.ErrManagerOverloaded):
		return servererrors.NewServerError(ErrorCodeManagerOverloaded, ManagerOverloadedError, err)
	case err != nil:
		return fmt.Errorf("freeHandsUseCase: %v", err)
	}
	if err = eCtx.JSON(http.StatusOK, FreeHandsResponse{Data: nil}); err != nil {
		return fmt.Errorf("send response on freeHands call: %v", err)
	}
	return nil
}
