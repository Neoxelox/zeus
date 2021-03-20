package database

import (
	"github.com/cockroachdb/errors"
)

var (
	ErrNoRows             = errors.New("No rows in result set")
	ErrIntegrityViolation = errors.New("Integrity constraint violation")
)

// TODO(alex): psql error interface.
