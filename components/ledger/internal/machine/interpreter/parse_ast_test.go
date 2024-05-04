package interpreter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSend(t *testing.T) {
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

func TestParseSeq(t *testing.T) {
	got := CompileFull(`send [COIN 42] (
		source = {
			@s1
			@s2
		}
		destination = @dest	
	)`)

	if len(got.Errors) != 0 {
		t.Fatalf(`Unexpected compilation errors = %#v`, got.Errors)
	}

	expected := Program{
		Statements: []SendStatement{
			{
				Amount: 42,
				Source: &SeqSrc{[]Source{
					&AccountSrc{Name: "s1"},
					&AccountSrc{Name: "s2"},
				}},
				Destination: &AccountDest{"dest"},
			},
		},
	}

	require.Equal(t, &expected, got.Program, "Program should be the same.")
}
