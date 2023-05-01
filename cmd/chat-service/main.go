package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	keycloakclient "github.com/gerladeno/chat-service/internal/clients/keycloak"
	"github.com/gerladeno/chat-service/internal/config"
	"github.com/gerladeno/chat-service/internal/logger"
	clientv1 "github.com/gerladeno/chat-service/internal/server-client/v1"
	serverdebug "github.com/gerladeno/chat-service/internal/server-debug"
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
		return fmt.Errorf("init logger: %w", err)
	}
	defer logger.Sync()

	swagger, err := clientv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("get swagger: %w", err)
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
		return fmt.Errorf("init keycloak client: %w", err)
	}
	srvClient, err := initServerClient(
		cfg.Servers.Client.Addr,
		cfg.Servers.Client.AllowOrigins,
		swagger,
		kcClient,
		cfg.Servers.Client.RequiredAccess.Resource,
		cfg.Servers.Client.RequiredAccess.Role,
	)
	if err != nil {
		return fmt.Errorf("init chat server: %w", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

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
