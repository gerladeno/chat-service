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
	jobsrepo "github.com/gerladeno/chat-service/internal/repositories/jobs"
	messagesrepo "github.com/gerladeno/chat-service/internal/repositories/messages"
	problemsrepo "github.com/gerladeno/chat-service/internal/repositories/problems"
	clientv1 "github.com/gerladeno/chat-service/internal/server-client/v1"
	serverdebug "github.com/gerladeno/chat-service/internal/server-debug"
	managerv1 "github.com/gerladeno/chat-service/internal/server-manager/v1"
	managerload "github.com/gerladeno/chat-service/internal/services/manager-load"
	inmemmanagerpool "github.com/gerladeno/chat-service/internal/services/manager-pool/in-mem"
	msgproducer "github.com/gerladeno/chat-service/internal/services/msg-producer"
	"github.com/gerladeno/chat-service/internal/services/outbox"
	sendclientmessagejob "github.com/gerladeno/chat-service/internal/services/outbox/jobs/send-client-message"
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

	// Config and logs
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

	// Swagger
	clientSwagger, err := clientv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("get client swagger: %v", err)
	}
	managerSwagger, err := managerv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("get manager swagger: %v", err)
	}

	// Keycloak
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

	// Storage
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

	// Repos
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
	jobsRepo, err := jobsrepo.New(jobsrepo.NewOptions(db))
	if err != nil {
		return fmt.Errorf("init jobs repo: %v", err)
	}

	// Init services
	msgProducer, err := msgproducer.New(msgproducer.NewOptions(
		msgproducer.NewKafkaWriter(
			cfg.Services.MsgProducer.Brokers,
			cfg.Services.MsgProducer.Topic,
			cfg.Services.MsgProducer.BatchSize,
		), msgproducer.WithEncryptKey(cfg.Services.MsgProducer.EncryptKey),
	))
	if err != nil {
		return fmt.Errorf("init msg producer: %v", err)
	}

	sendClientMessageJob, err := sendclientmessagejob.New(sendclientmessagejob.NewOptions(
		msgProducer,
		msgRepo,
	))
	if err != nil {
		return fmt.Errorf("init send client message job: %v", err)
	}

	outboxService, err := outbox.New(outbox.NewOptions(
		cfg.Services.Outbox.Workers,
		cfg.Services.Outbox.IdleTime,
		cfg.Services.Outbox.ReserveFor,
		jobsRepo,
		db,
	))
	if err != nil {
		return fmt.Errorf("init outbox service: %v", err)
	}

	managerPool := inmemmanagerpool.New()

	managerLoad, err := managerload.New(managerload.NewOptions(
		cfg.Services.ManagerLoad.MaxProblemsAtSameTime,
		problemsRepo,
	))
	if err != nil {
		return fmt.Errorf("init manager load service")
	}

	outboxService.MustRegisterJob(sendClientMessageJob)

	// Init servers
	srvDebug, err := serverdebug.New(serverdebug.NewOptions(
		cfg.Servers.Debug.Addr,
		clientSwagger,
		managerSwagger,
	))
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}

	srvManager, err := initServerManager(
		cfg.Global.IsProd(),
		cfg.Servers.Manager.Addr,
		cfg.Servers.Manager.AllowOrigins,
		managerSwagger,

		kcClient,
		cfg.Servers.Manager.RequiredAccess.Resource,
		cfg.Servers.Manager.RequiredAccess.Role,

		managerLoad,
		managerPool,
	)
	if err != nil {
		return fmt.Errorf("init manager chat server: %v", err)
	}

	srvClient, err := initServerClient(
		cfg.Global.IsProd(),
		cfg.Servers.Client.Addr,
		cfg.Servers.Client.AllowOrigins,
		clientSwagger,

		kcClient,
		cfg.Servers.Client.RequiredAccess.Resource,
		cfg.Servers.Client.RequiredAccess.Role,

		db,
		msgRepo,
		chatRepo,
		problemsRepo,
		outboxService,
	)
	if err != nil {
		return fmt.Errorf("init client chat server: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		<-ctx.Done()
		return psqlClient.Close()
	})
	// Run servers and services.
	eg.Go(func() error { return srvDebug.Run(ctx) })

	eg.Go(func() error { return srvClient.Run(ctx) })

	eg.Go(func() error { return srvManager.Run(ctx) })

	eg.Go(func() error { return outboxService.Run(ctx) })
	// Ждут своего часа.
	// ...

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	return nil
}
