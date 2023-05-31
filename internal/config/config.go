package config

import "time"

type Config struct {
	Global   GlobalConfig  `toml:"global"`
	Log      LogConfig     `toml:"log"`
	Servers  ServersConfig `toml:"servers"`
	Sentry   SentryConfig  `toml:"sentry"`
	Clients  ClientConfig  `toml:"clients"`
	DB       DBConfig      `toml:"db"`
	Services ServiceConfig `toml:"services"`
}

type GlobalConfig struct {
	Env string `toml:"env" validate:"required,oneof=dev stage prod"`
}

func (gc GlobalConfig) IsProd() bool {
	return gc.Env == "prod"
}

type LogConfig struct {
	Level string `toml:"level" validate:"required,oneof=debug info warn error"`
}

type ServersConfig struct {
	Debug   DebugServerConfig `toml:"debug"`
	Client  ServerConfig      `toml:"client"`
	Manager ServerConfig      `toml:"manager"`
}

type DebugServerConfig struct {
	Addr string `toml:"addr" validate:"required,hostname_port"`
}

type SentryConfig struct {
	DSN string `toml:"dsn"`
}

type ServerConfig struct {
	Addr           string         `toml:"addr" validate:"required,hostname_port"`
	AllowOrigins   []string       `toml:"allow_origins" validate:"required"`
	SecWSProtocol  string         `toml:"sec_ws_protocol" validate:"required"`
	RequiredAccess RequiredAccess `toml:"required_access" validate:"required"`
}

type RequiredAccess struct {
	Resource string `toml:"resource" validate:"required"`
	Role     string `toml:"role" validate:"required"`
}

type ClientConfig struct {
	Keycloak Keycloak `toml:"keycloak" validate:"required"`
}

type Keycloak struct {
	BasePath     string `toml:"base_path" validate:"required"`
	Realm        string `toml:"realm" validate:"required"`
	ClientID     string `toml:"client_id" validate:"required"`
	ClientSecret string `toml:"client_secret" validate:"required"`
	DebugMode    bool   `toml:"debug_mode"`
}

type DBConfig struct {
	Postgres PGConfig `toml:"postgres"`
}

type PGConfig struct {
	User      string `toml:"user" validate:"required"`
	Password  string `toml:"password" validate:"required"`
	Addr      string `toml:"addr" validate:"required,hostname_port"`
	Database  string `toml:"database" validate:"required"`
	DebugMode bool   `toml:"debug_mode"`
}

type ServiceConfig struct {
	MsgProducer         MsgProducerConfig         `toml:"msg_producer"`
	Outbox              OutboxConfig              `toml:"outbox"`
	ManagerLoad         ManagerLoadConfig         `toml:"manager_load"`
	AFCVerdictProcessor AFCVerdictProcessorConfig `toml:"afc_verdicts_processor"`
	ManagerScheduler    ManagerSchedulerConfig    `toml:"manager_scheduler"`
}

type MsgProducerConfig struct {
	Brokers    []string `toml:"brokers" validate:"required,dive,hostname_port"`
	Topic      string   `toml:"topic" validate:"required"`
	BatchSize  int      `toml:"batch_size"`
	EncryptKey string   `toml:"encrypt_key"`
}

type OutboxConfig struct {
	Workers    int           `toml:"workers" validate:"required"`
	IdleTime   time.Duration `toml:"idle_time" validate:"required"`
	ReserveFor time.Duration `toml:"reserve_for"`
}

type ManagerLoadConfig struct {
	MaxProblemsAtSameTime int `toml:"max_problems_at_same_time" validate:"required,min=1,max=30"`
}

type AFCVerdictProcessorConfig struct {
	BackoffInitialInterval time.Duration `toml:"backoff_initial_interval" validate:"min=50ms,max=1s"`
	BackoffMaxElapsedTime  time.Duration `toml:"backoff_max_elapsed_time" validate:"min=500ms,max=1m"`
	BackoffFactor          float64       `toml:"backoff_factor" validate:"min=1.01,max=10"`
	Brokers                []string      `toml:"brokers" validate:"required,dive,hostname_port"`
	Consumers              int           `toml:"consumers" validate:"required,min=1"`
	ConsumerGroup          string        `toml:"consumer_group" validate:"required"`
	VerdictTopic           string        `toml:"verdict_topic" validate:"required"`
	VerdictSignKey         string        `toml:"verdicts_signing_public_key"`
	VerdictTopicDLQ        string        `toml:"verdict_topic_dlq" validate:"required"`
	ProcessBatchSize       int           `toml:"process_batch_size" validate:"min=1,max=1000"`
	ProcessBatchMaxTimeout time.Duration `toml:"process_batch_max_timeout" validate:"min=50ms,max=10s"`
	Retries                int           `toml:"retries" validate:"min=1,max=10"`
}

type ManagerSchedulerConfig struct {
	Period time.Duration `toml:"period" validate:"required"`
}
