package database

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Configuration describes the Database configuration.
type Configuration struct {
	MinConns int
	MaxConns int
	AppName  string
	Logger   pgx.Logger
	LogLevel pgx.LogLevel
}

// Database describes the database.
type Database struct {
	connection    *pgxpool.Pool
	configuration Configuration
}

// New creates a new Database instance.
func New(ctx context.Context, dsn string, retries int, configuration Configuration) (*Database, error) {
	delay := time.NewTicker(1 * time.Second)
	timeout := (time.Duration(retries) * time.Second)
	defer delay.Stop()

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
					connection:    connection,
					configuration: configuration,
				}, nil
			}
		}
	}
}

// Close shutdowns any connection to the database.
func (d *Database) Close(ctx context.Context) error {
	d.configuration.Logger.Log(ctx, pgx.LogLevelInfo, "Closing database", nil)
	d.connection.Close()

	return nil
}

// Health checks if the database is reachable and running.
func (d *Database) Health(ctx context.Context) error {
	if _, err := d.connection.Exec(ctx, ";"); err != nil {
		return errors.Wrap(err, "Database unreachable")
	}

	if d.connection.Stat().TotalConns() < int32(d.configuration.MinConns) {
		return errors.New("Database pool size below minimum")
	}

	return nil
}
