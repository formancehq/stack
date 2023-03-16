package errors

import (
	"github.com/formancehq/stack/libs/go-libs/sqlstorage/sqlerrors"
	"github.com/lib/pq"
)

// postgresError is an helper to wrap postgres errors into storage errors
func PostgresError(err error) error {
	if err != nil {
		switch pge := err.(type) {
		case *pq.Error:
			switch pge.Code {
			case "23505":
				return sqlerrors.NewError(sqlerrors.ConstraintFailed, err)
			case "53300":
				return sqlerrors.NewError(sqlerrors.TooManyClient, err)
			}
		}
	}

	return err
}
