package pg

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	pgErrUniqueViolation     = "23505"
	pgErrForeignKeyViolation = "23503"
)

func pgErr(err error) (*pgconn.PgError, bool) {
	var pge *pgconn.PgError
	ok := errors.As(err, &pge)

	return pge, ok
}

func isUniqueViolation(err error) bool {
	if pge, ok := pgErr(err); ok {
		return pge.Code == pgErrUniqueViolation
	}

	return false
}
