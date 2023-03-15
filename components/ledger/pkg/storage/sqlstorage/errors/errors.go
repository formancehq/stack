package errors

import (
	"github.com/formancehq/stack/libs/go-libs/api/apierrors"
	"github.com/lib/pq"
)

// postgresError is an helper to wrap postgres errors into storage errors
func PostgresError(err error) error {
	if err != nil {
		switch pge := err.(type) {
		case *pq.Error:
			switch pge.Code {
			case "23505":
				return apierrors.NewStorageError(apierrors.ConstraintFailed, err)
			case "53300":
				return apierrors.NewStorageError(apierrors.TooManyClient, err)
			}
		}
	}

	return err
}
