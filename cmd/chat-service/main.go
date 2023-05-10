package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/sync/errgroup"

	keycloakclient "github.com/gerladeno/chat-service/internal/clients/keycloak"
	"github.com/gerladeno/chat-service/internal/config"
	"github.com/gerladeno/chat-service/internal/logger"
	chatsrepo "github.com/gerladeno/chat-service/internal/repositories/chats"
	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	problemsrepo "github.com/gerladeno/chat-service/internal/repositories/problems"
	clientv1 "github.com/gerladeno/chat-service/internal/server-client/v1"
	serverdebug "github.com/gerladeno/chat-service/internal/server-debug"
	"github.com/gerladeno/chat-service/internal/store"
)

var configPath = flag.String("config", "configs/config.toml", "Path to config file")

func main() {
	if err := run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}

func run() (errReturned error) {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.ParseAndValidate(*configPath)
	if err != nil {
		return fmt.Errorf("parse and validate config %q: %v", *configPath, err)
	}

	if err = logger.Init(logger.NewOptions(
		cfg.Log.Level,
		logger.WithProductionMode(cfg.Global.IsProd()),
		logger.WithSentryDSN(cfg.Sentry.DSN),
		logger.WithEnv(cfg.Global.Env),
	)); err != nil {
		return fmt.Errorf("init logger: %v", err)
	}
	defer logger.Sync()

	swagger, err := clientv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("get swagger: %v", err)
	}
	srvDebug, err := serverdebug.New(serverdebug.NewOptions(
		cfg.Servers.Debug.Addr,
		swagger,
	))
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}

	kcClient, err := keycloakclient.New(keycloakclient.NewOptions(
		cfg.Clients.Keycloak.BasePath,
		cfg.Clients.Keycloak.Realm,
		cfg.Clients.Keycloak.ClientID,
		cfg.Clients.Keycloak.ClientSecret,
		cfg.Global.IsProd(),
		keycloakclient.WithDebugMode(cfg.Clients.Keycloak.DebugMode),
	))
	if err != nil {
		return fmt.Errorf("init keycloak client: %v", err)
	}

	psqlClient, err := store.NewPSQLClient(store.NewPSQLOptions(
		cfg.DB.Postgres.Addr,
		cfg.DB.Postgres.User,
		cfg.DB.Postgres.Password,
		cfg.DB.Postgres.Database,
		store.WithDebug(cfg.DB.Postgres.DebugMode),
	))
	if err != nil {
		return fmt.Errorf("init psql client: %v", err)
	}
	if err = psqlClient.Schema.Create(ctx); err != nil {
		return fmt.Errorf("migrate schema: %v", err)
	}

	db := store.NewDatabase(psqlClient)
	msgRepo, err := messagesrepo.New(messagesrepo.NewOptions(db))
	if err != nil {
		return fmt.Errorf("init messages repo: %v", err)
	}
	chatRepo, err := chatsrepo.New(chatsrepo.NewOptions(db))
	if err != nil {
		return fmt.Errorf("init chats repo: %v", err)
	}
	problemsRepo, err := problemsrepo.New(problemsrepo.NewOptions(db))
	if err != nil {
		return fmt.Errorf("init problems repo: %v", err)
	}

	srvClient, err := initServerClient(
		cfg.Servers.Client.Addr,
		cfg.Servers.Client.AllowOrigins,
		swagger,
		kcClient,
		cfg.Servers.Client.RequiredAccess.Resource,
		cfg.Servers.Client.RequiredAccess.Role,
		cfg.Global.IsProd(),
		msgRepo,
		chatRepo,
		problemsRepo,
		db,
	)
	if err != nil {
		return fmt.Errorf("init chat server: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		<-ctx.Done()
		return psqlClient.Close()
	})
	// Run servers.
	eg.Go(func() error { return srvDebug.Run(ctx) })

	eg.Go(func() error { return srvClient.Run(ctx) })
	// Run services.
	// Ждут своего часа.
	// ...

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	return nil
}
