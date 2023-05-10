package store

import (
	"database/sql"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"fmt"
)

//go:generate options-gen -out-filename=client_psql_options.gen.go -from-struct=PSQLOptions
type PSQLOptions struct {
	address  string `option:"mandatory" validate:"required,hostname_port"`
	username string `option:"mandatory" validate:"required"`
	password string `option:"mandatory" validate:"required"`
	database string `option:"mandatory" validate:"required"`
	debug    bool
}

func NewPSQLClient(opts PSQLOptions) (*Client, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}

	db, err := NewPgxDB(NewPgxOptions(opts.address, opts.username, opts.password, opts.database))
	if err != nil {
		return nil, fmt.Errorf("init db driver: %v", err)
	}

	clientOpts := []Option{Driver(entsql.OpenDB(dialect.Postgres, db))}
	if opts.debug {
		clientOpts = append(clientOpts, Debug())
	}
	return NewClient(clientOpts...), nil
}

//go:generate options-gen -out-filename=client_psql_pgx_options.gen.go -from-struct=PgxOptions
type PgxOptions struct {
	address  string `option:"mandatory" validate:"required,hostname_port"`
	username string `option:"mandatory" validate:"required"`
	password string `option:"mandatory" validate:"required"`
	database string `option:"mandatory" validate:"required"`
}

const pgxDSNTmpl = `postgresql://%s:%s@%s/%s?sslmode=disable`

func NewPgxDB(opts PgxOptions) (*sql.DB, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}
	db, err := sql.Open("pgx", fmt.Sprintf(pgxDSNTmpl, opts.username, opts.password, opts.address, opts.database))
	if err != nil {
		return nil, fmt.Errorf("connect to pg: %v", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping pg: %v", err)
	}
	return db, nil
}