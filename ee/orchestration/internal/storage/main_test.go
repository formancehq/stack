package storage

import (
	"os"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/pgtesting"
)

func TestMain(t *testing.M) {
	if err := pgtesting.CreatePostgresServer(); err != nil {
		panic(err)
	}

	code := t.Run()

	if err := pgtesting.DestroyPostgresServer(); err != nil {
		panic(err)
	}

	os.Exit(code)
}
