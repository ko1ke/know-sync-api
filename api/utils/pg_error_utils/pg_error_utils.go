package pg_error_utils

import (
	"errors"

	"github.com/jackc/pgconn"
)

// Parse errors only database can detect.
func ParseError(err error) error {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return err
	}

	switch pgErr.Code {
	case "23505":
		switch pgErr.Message {
		case "duplicate key value violates unique constraint \"users_email_key\"":
			return errors.New("Eメールは既に登録されています")
		default:
			return pgErr
		}
	}
	return err
}
