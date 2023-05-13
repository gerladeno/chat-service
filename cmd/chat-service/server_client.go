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
	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	problemsrepo "github.com/gerladeno/chat-service/internal/repositories/problems"
	"github.com/gerladeno/chat-service/internal/server"
	clientv1 "github.com/gerladeno/chat-service/internal/server-client/v1"
	"github.com/gerladeno/chat-service/internal/server/errhandler"
	"github.com/gerladeno/chat-service/internal/services/outbox"
	"github.com/gerladeno/chat-service/internal/store"
	gethistory "github.com/gerladeno/chat-service/internal/usecases/client/get-history"
	sendmessage "github.com/gerladeno/chat-service/internal/usecases/client/send-message"
	websocketstream "github.com/gerladeno/chat-service/internal/websocket-stream"
)

const nameServerClient = "server-client"

func initServerClient(
	isProd bool,
	addr string,
	allowOrigins []string,
	v1Swagger *openapi3.T,

	client *keycloakclient.Client,
	resource string,
	role string,
	wsSecProtocol string,

	db *store.Database,
	msgRepo *messagesrepo.Repo,
	chatRepo *chatsrepo.Repo,
	problemRepo *problemsrepo.Repo,

	outboxService *outbox.Service,
	wsHandler *websocketstream.HTTPHandler,
) (*server.Server, error) {
	lg := zap.L().Named(nameServerClient)

	getHistoryUseCase, err := gethistory.New(gethistory.NewOptions(msgRepo))
	if err != nil {
		return nil, fmt.Errorf("create getHistoryUseCase: %v", err)
	}
	sendMessageUseCase, err := sendmessage.New(sendmessage.NewOptions(chatRepo, msgRepo, outboxService, problemRepo, db))
	if err != nil {
		return nil, fmt.Errorf("create sendMessageUseCase: %v", err)
	}

	v1Handlers, err := clientv1.NewHandlers(clientv1.NewOptions(lg, getHistoryUseCase, sendMessageUseCase))
	if err != nil {
		return nil, fmt.Errorf("create v1 handlers: %v", err)
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
			clientv1.RegisterHandlers(v1, v1Handlers)
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
