package database

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // nolint
	_ "github.com/golang-migrate/migrate/v4/source/file"       // nolint
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Configuration describes the Database configuration.
type Configuration struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
	MinConns int
	MaxConns int
	AppName  string
	Logger   pgx.Logger
	LogLevel pgx.LogLevel
}

// Database describes the database.
type Database struct {
	Connection    *pgxpool.Pool
	configuration Configuration
}

// New creates a new Database instance.
func New(ctx context.Context, retries int, configuration Configuration) (*Database, error) {
	delay := time.NewTicker(1 * time.Second)
	timeout := (time.Duration(retries) * time.Second)
	defer delay.Stop()

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		configuration.User,
		configuration.Password,
		configuration.Host,
		configuration.Port,
		configuration.Name,
		configuration.SSLMode,
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot parse database connection string")
	}

	config.MinConns = int32(configuration.MinConns)
	config.MaxConns = int32(configuration.MaxConns)
	config.ConnConfig.RuntimeParams["standard_conforming_strings"] = "on"
	config.ConnConfig.RuntimeParams["application_name"] = configuration.AppName

	config.ConnConfig.Logger = configuration.Logger
	config.ConnConfig.LogLevel = configuration.LogLevel

	timeoutExceeded := time.After(timeout)
	for {
		select {
		case <-timeoutExceeded:
			return nil, errors.New("Cannot connect to the database")
		case <-delay.C:
			configuration.Logger.Log(ctx, pgx.LogLevelInfo, "Trying to connect to the database", nil)
			connection, err := pgxpool.ConnectConfig(ctx, config)
			if err == nil {
				configuration.Logger.Log(ctx, pgx.LogLevelInfo, "Connected to the database", nil)

				return &Database{
					Connection:    connection,
					configuration: configuration,
				}, nil
			}
		}
	}
}

// Close shutdowns any connection to the database.
func (d *Database) Close(ctx context.Context) error {
	d.configuration.Logger.Log(ctx, pgx.LogLevelInfo, "Closing database", nil)
	d.Connection.Close()

	return nil
}

// Health checks if the database is reachable and running.
func (d *Database) Health(ctx context.Context) error {
	if _, err := d.Connection.Exec(ctx, ";"); err != nil {
		return errors.Wrap(err, "Database unreachable")
	}

	if d.Connection.Stat().TotalConns() < int32(d.configuration.MinConns) {
		return errors.New("Database pool size below minimum")
	}

	return nil
}

// Migrate runs database migrations up to the latest version.
func (d *Database) Migrate(ctx context.Context) error {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s&x-multi-statement=true",
		d.configuration.User,
		d.configuration.Password,
		d.configuration.Host,
		d.configuration.Port,
		d.configuration.Name,
		d.configuration.SSLMode,
	)

	migrator, err := migrate.New("file://./migrations", dsn)
	if err != nil {
		return errors.Wrap(err, "Cannot begin migrator")
	}

	err = migrator.Up()
	switch err { // nolint
	case nil:
		d.configuration.Logger.Log(ctx, pgx.LogLevelInfo, "Applying migrations", nil)
	case migrate.ErrNoChange:
		d.configuration.Logger.Log(ctx, pgx.LogLevelInfo, "No migrations to apply", nil)
	default:
		return errors.Wrap(err, "Error within a migration")
	}

	return nil
}

// Connection is a querier and execerer to interact with the database.
type Connection interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

// BeginTransaction starts a database transaction.
func BeginTransaction(ctx context.Context, db *pgxpool.Pool) (pgx.Tx, error) {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Cannot begin transaction")
	}

	return tx, nil
}

// WatchTransaction couples a watchdog to the given transaction that rollbacks it if panic occurs.
func WatchTransaction(ctx context.Context, tx pgx.Tx) func() {
	return func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx) // nolint
			panic(p)
		}
	}
}

// FinishTransaction commits a database transaction.
func FinishTransaction(ctx context.Context, err error, tx pgx.Tx) error {
	if err != nil {
		if rerr := tx.Rollback(ctx); rerr != nil {
			return errors.Wrap(rerr, "Cannot rollback transaction")
		}

		return errors.Wrap(err, "Error within a transaction")
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "Cannot commit transaction")
	}

	return nil
}
