package searchhttp

import (
	"testing"

	"github.com/aquasecurity/esquery"
	"github.com/formancehq/search/pkg/searchengine"
	"github.com/stretchr/testify/require"
)

func TestNextToken(t *testing.T) {
	nti := cursorTokenInfo{
		Target: "ACCOUNT",
		Sort: []searchengine.Sort{
			{
				Key:   "slug",
				Order: esquery.OrderDesc,
			},
		},
		SearchAfter: []interface{}{
			"ACCOUNT-2",
		},
		Ledgers: []string{"quickstart"},
	}
	tok := EncodeCursorToken(nti)
	decoded := cursorTokenInfo{}
	require.NoError(t, DecodeCursorToken(tok, &decoded))
	require.EqualValues(t, nti, decoded)
}
