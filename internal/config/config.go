package config

type Config struct {
	Global  GlobalConfig  `toml:"global"`
	Log     LogConfig     `toml:"log"`
	Servers ServersConfig `toml:"servers"`
	Sentry  SentryConfig  `toml:"sentry"`
	Clients ClientConfig  `toml:"clients"`
	DB      DBConfig      `toml:"db"`
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
	Debug  DebugServerConfig  `toml:"debug"`
	Client ClientServerConfig `toml:"client"`
}

type DebugServerConfig struct {
	Addr string `toml:"addr" validate:"required,hostname_port"`
}

type SentryConfig struct {
	DSN string `toml:"dsn"`
}

type ClientServerConfig struct {
	Addr           string         `toml:"addr" validate:"required,hostname_port"`
	AllowOrigins   []string       `toml:"allow_origins" validate:"required"`
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
