package interpreter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
	got := CompileFull(`send [COIN 42] (
		source = @src
		destination = @dest	
	)`)

	if len(got.Errors) != 0 {
		t.Fatalf(`Unexpected compilation errors = %#v`, got.Errors)
	}

	expected := Program{
		Statements: []SendStatement{
			{
				Amount:      42,
				Source:      &AccountSrc{"src"},
				Destination: &AccountDest{"dest"},
			},
		},
	}

	require.Equal(t, &expected, got.Program, "Program should be the same.")

}
