package wallet

import (
	"encoding/json"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestAccountUnmarshal(t *testing.T) {
	type accountResponse struct {
		Data AccountWithVolumesAndBalances `json:"data"`
	}
	account := accountResponse{}
	data := `{"data":{"address":"wallets:4cf8f8c0015f4e5e90b8c2039e8384c1:main","metadata":{"wallets/balances":"true","wallets/balances/name":"main","wallets/createdAt":"2023-05-16T11:41:02.413258715Z","wallets/custom_data":{},"wallets/id":"4cf8f8c0-015f-4e5e-90b8-c2039e8384c1","wallets/name":"test","wallets/spec/type":"wallets.primary"},"volumes":{"USD":{"input":36893488147427550068,"output":0,"balance":36893488147427550068}},"balances":{"USD":36893488147427550068}}}`
	require.NoError(t, json.Unmarshal([]byte(data), &account))

	spew.Dump(account)
}
