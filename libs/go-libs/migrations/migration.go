package migrations

import (
	"database/sql"
)

type Migration struct {
	Up func(tx *sql.Tx) error
}
