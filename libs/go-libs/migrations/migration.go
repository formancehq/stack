package migrations

import (
	"github.com/uptrace/bun"
)

type Migration struct {
	Name string
	Up   func(tx bun.Tx) error
}
