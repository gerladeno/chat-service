package main

import (
	"fmt"

	oapimdlwr "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	keycloakclient "github.com/gerladeno/chat-service/internal/clients/keycloak"
	chatsrepo "github.com/gerladeno/chat-service/internal/repositories/chats"
	"github.com/gerladeno/chat-service/internal/server"
	managerv1 "github.com/gerladeno/chat-service/internal/server-manager/v1"
	"github.com/gerladeno/chat-service/internal/server/errhandler"
	managerload "github.com/gerladeno/chat-service/internal/services/manager-load"
	managerpool "github.com/gerladeno/chat-service/internal/services/manager-pool"
	canreceiveproblems "github.com/gerladeno/chat-service/internal/usecases/manager/can-receive-problems"
	freehands "github.com/gerladeno/chat-service/internal/usecases/manager/free-hands"
	getchats "github.com/gerladeno/chat-service/internal/usecases/manager/get-chats"
	websocketstream "github.com/gerladeno/chat-service/internal/websocket-stream"
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
	wsSecProtocol string,

	chatRepo *chatsrepo.Repo,

	managerLoad *managerload.Service,
	managerPool managerpool.Pool,
	wsHandler *websocketstream.HTTPHandler,
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
	getChatsUseCase, err := getchats.New(getchats.NewOptions(chatRepo))
	if err != nil {
		return nil, fmt.Errorf("init getChatsUseCase: %v", err)
	}

	v1Handlers, err := managerv1.NewHandlers(managerv1.NewOptions(
		lg,
		canReceiveProblemsUseCase,
		freeHandsUseCase,
		getChatsUseCase,
	))
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

		func(e *echo.Echo) {
			e.GET("/ws", wsHandler.Serve)
			v1 := e.Group("v1", oapimdlwr.OapiRequestValidatorWithOptions(v1Swagger, &oapimdlwr.Options{
				Options: openapi3filter.Options{
					ExcludeRequestBody:  false,
					ExcludeResponseBody: true,
					AuthenticationFunc:  openapi3filter.NoopAuthenticationFunc,
				},
			}))
			managerv1.RegisterHandlers(v1, v1Handlers)
		},
		client,
		resource,
		role,
		wsSecProtocol,
		errHandler.Handle,
	))
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return srv, nil
}
