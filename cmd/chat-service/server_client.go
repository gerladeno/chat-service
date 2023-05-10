package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	keycloakclient "github.com/gerladeno/chat-service/internal/clients/keycloak"
	chatsrepo "github.com/gerladeno/chat-service/internal/repositories/chats"
	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	problemsrepo "github.com/gerladeno/chat-service/internal/repositories/problems"
	serverclient "github.com/gerladeno/chat-service/internal/server-client"
	"github.com/gerladeno/chat-service/internal/server-client/errhandler"
	clientv1 "github.com/gerladeno/chat-service/internal/server-client/v1"
	"github.com/gerladeno/chat-service/internal/store"
	gethistory "github.com/gerladeno/chat-service/internal/usecases/client/get-history"
	sendmessage "github.com/gerladeno/chat-service/internal/usecases/client/send-message"
)

const nameServerClient = "server-client"

func initServerClient(
	addr string,
	allowOrigins []string,
	v1Swagger *openapi3.T,
	client *keycloakclient.Client,
	resource string,
	role string,
	isProd bool,
	msgRepo *messagesrepo.Repo,
	chatRepo *chatsrepo.Repo,
	problemRepo *problemsrepo.Repo,
	db *store.Database,
) (*serverclient.Server, error) {
	lg := zap.L().Named(nameServerClient)

	getHistoryUseCase, err := gethistory.New(gethistory.NewOptions(msgRepo))
	if err != nil {
		return nil, fmt.Errorf("create getHistoryUseCase: %v", err)
	}
	sendMessageUseCae, err := sendmessage.New(sendmessage.NewOptions(chatRepo, msgRepo, problemRepo, db))
	if err != nil {
		return nil, fmt.Errorf("create sendMessageUseCae: %v", err)
	}

	v1Handlers, err := clientv1.NewHandlers(clientv1.NewOptions(lg, getHistoryUseCase, sendMessageUseCae))
	if err != nil {
		return nil, fmt.Errorf("create v1 handlers: %v", err)
	}

	errHandler, err := errhandler.New(errhandler.NewOptions(lg, isProd, errhandler.ResponseBuilder))
	if err != nil {
		return nil, fmt.Errorf("create error handler: %v", err)
	}

	srv, err := serverclient.New(serverclient.NewOptions(
		lg,
		addr,
		allowOrigins,
		v1Swagger,
		v1Handlers,
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
