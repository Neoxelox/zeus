package database

import (
	"regexp"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

var codeExtractor = regexp.MustCompile(`\(SQLSTATE (.*)\)`)

var (
	ErrNoRows             = errors.New("No rows in result set")
	ErrIntegrityViolation = errors.New("Integrity constraint violation")
)

// Error transforms error into a database layer error.
func Error(err error) error {
	if err == nil {
		return nil
	}

	if code := codeExtractor.FindStringSubmatch(err.Error()); len(code) == 2 { // nolint
		if cErr := codeError(code[1]); cErr != nil {
			return cErr
		}
	}

	if cErr := symmetricError(err); cErr != nil {
		return cErr
	}

	return err
}

// nolint
func codeError(code string) error {
	switch code {
	case pgerrcode.IntegrityConstraintViolation, pgerrcode.RestrictViolation, pgerrcode.NotNullViolation,
		pgerrcode.ForeignKeyViolation, pgerrcode.UniqueViolation, pgerrcode.CheckViolation,
		pgerrcode.ExclusionViolation:
		return ErrIntegrityViolation
	default:
		return nil
	}
}

// nolint
func symmetricError(err error) error {
	switch err.Error() {
	case pgx.ErrNoRows.Error():
		return ErrNoRows
	default:
		return nil
	}
}
