package server

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Schemes enumerates the possible schemes.
var Schemes = struct {
	HTTP  string
	HTTPS string
}{"http", "https"}

// Environments enumerates the possible environments.
var Environments = struct {
	PRODUCTION  string
	STAGING     string
	DEVELOPMENT string
	TESTING     string
}{"production", "staging", "development", "testing"}

type (
	_app struct {
		Host            []string
		Port            int
		Scheme          string
		Environment     string
		Name            string
		Version         string
		Release         string
		GracefulTimeout int
	}

	_database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		SSLMode  string
		Dsn      string
	}

	// Configuration describes the application configuration.
	Configuration struct {
		App      _app
		Database _database
	}
)

func (s *Server) addConfiguration() error { // nolint
	s.Configuration = Configuration{
		App: _app{
			Host:            getEnvAsSlice("ZEUS_HOST", []string{"localhost"}),
			Port:            getEnvAsInt("ZEUS_PORT", 1111),
			Scheme:          getEnvAsString("ZEUS_SCHEME", "http"),
			Environment:     getEnvAsString("ZEUS_ENVIRONMENT", "development"),
			Name:            getEnvAsString("ZEUS_NAME", "zeus"),
			Version:         getEnvAsString("ZEUS_VERSION", "fakeVersion"),
			Release:         getEnvAsString("ZEUS_RELEASE", "fakeRelease"),
			GracefulTimeout: getEnvAsInt("ZEUS_GRACEFUL_TIMEOUT", 15),
		},

		Database: _database{
			Host:     getEnvAsString("DATABASE_HOST", "postgres"),
			Port:     getEnvAsInt("DATABASE_PORT", 5432),
			User:     getEnvAsString("DATABASE_USER", "zeus"),
			Password: getEnvAsString("DATABASE_PASSWORD", "zeus"),
			Name:     getEnvAsString("DATABASE_NAME", "zeus"),
			SSLMode:  getEnvAsString("DATABASE_SSLMODE", "disable"),
			Dsn: fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
				getEnvAsString("DATABASE_USER", "zeus"),
				getEnvAsString("DATABASE_PASSWORD", "zeus"),
				getEnvAsString("DATABASE_HOST", "postgres"),
				getEnvAsInt("DATABASE_PORT", 5432),
				getEnvAsString("DATABASE_NAME", "zeus"),
				getEnvAsString("DATABASE_SSLMODE", "disable"),
			),
		},
	}

	return nil
}

func getEnvAsString(key string, def string) string { // nolint
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return def
}

func getEnvAsInt(key string, def int) int { // nolint
	valueStr := getEnvAsString(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return def
}

func getEnvAsBool(key string, def bool) bool { // nolint
	valueStr := getEnvAsString(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}

	return def
}

func getEnvAsSlice(key string, def []string) []string { // nolint
	valueStr := getEnvAsString(key, "")
	if value := strings.Split(valueStr, ","); len(value) >= 1 {
		return value
	}

	return def
}
