package shared

import (
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func IsVersionConflict(err error) bool {
	var pgErr *pgconn.PgError
	ok := errors.As(err, &pgErr)
	return ok && pgErr.Code == pgerrcode.UniqueViolation
}
