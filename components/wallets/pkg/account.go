package wallet

import (
	"github.com/formancehq/stack/libs/go-libs/metadata"
)

type Account interface {
	metadata.Owner
	GetAddress() string
}
