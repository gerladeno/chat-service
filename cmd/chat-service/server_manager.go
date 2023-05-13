package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	keycloakclient "github.com/gerladeno/chat-service/internal/clients/keycloak"
	"github.com/gerladeno/chat-service/internal/server"
	managerv1 "github.com/gerladeno/chat-service/internal/server-manager/v1"
	"github.com/gerladeno/chat-service/internal/server/errhandler"
	managerload "github.com/gerladeno/chat-service/internal/services/manager-load"
	managerpool "github.com/gerladeno/chat-service/internal/services/manager-pool"
	canreceiveproblems "github.com/gerladeno/chat-service/internal/usecases/manager/can-receive-problems"
	freehands "github.com/gerladeno/chat-service/internal/usecases/manager/free-hands"
)

const nameServerManager = `server-manager`

func initServerManager(
	isProd bool,
	addr string,
	allowOrigins []string,
	v1Swagger *openapi3.T,

	client *keycloakclient.Client,
	resource string,
	role string,

	managerLoad *managerload.Service,
	managerPool managerpool.Pool,
) (*server.Server, error) {
	lg := zap.L().Named(nameServerManager)

	canReceiveProblemsUseCase, err := canreceiveproblems.New(canreceiveproblems.NewOptions(
		managerLoad,
		managerPool,
	))
	if err != nil {
		return nil, fmt.Errorf("initing canReceiveProblemsUseCase: %v", err)
	}
	freeHandsUseCase, err := freehands.New(freehands.NewOptions(
		managerLoad,
		managerPool,
	))
	if err != nil {
		return nil, fmt.Errorf("initing freeHandsUseCase: %v", err)
	}

	v1Handlers, err := managerv1.NewHandlers(managerv1.NewOptions(lg, canReceiveProblemsUseCase, freeHandsUseCase))
	if err != nil {
		return nil, fmt.Errorf("initing v1Handlers: %v", err)
	}

	errHandler, err := errhandler.New(errhandler.NewOptions(lg, isProd, errhandler.ResponseBuilder))
	if err != nil {
		return nil, fmt.Errorf("create error handler: %v", err)
	}

	srv, err := server.New(server.NewOptions(
		lg,
		addr,
		allowOrigins,
		v1Swagger,

		func(g *echo.Group) {
			managerv1.RegisterHandlers(g, v1Handlers)
		},
		client,
		resource,
		role,
		errHandler.Handle,
	))
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return srv, nil
}
