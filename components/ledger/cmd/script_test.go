package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
	"github.com/numary/ledger/pkg/api/controllers"
	"github.com/numary/ledger/pkg/core"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func Test_ScriptCommands(t *testing.T) {
	name := uuid.NewString()
	viper.Set("name", name)
	_ = NewStorageInit().Execute()

	d1 := []byte(`
		send [EUR 1] (
			source = @world
			destination = @alice
		)`)
	require.NoError(t, os.WriteFile("/tmp/script", d1, 0644))

	cmd := NewScriptCheck()
	cmd.SetArgs([]string{"/tmp/script"})
	_ = cmd.Execute()

	httpServer := httptest.NewServer(http.HandlerFunc(scriptSuccessHandler))
	defer func() {
		httpServer.CloseClientConnections()
		httpServer.Close()
	}()

	viper.Set(serverHttpBindAddressFlag, httpServer.URL[7:])
	cmd = NewScriptExec()
	cmd.SetArgs([]string{name, "/tmp/script"})
	_ = cmd.Execute()

	cmd = NewScriptExec()
	viper.Set(previewFlag, true)
	cmd.SetArgs([]string{name, "/tmp/script"})
	_ = cmd.Execute()
}

func scriptSuccessHandler(w http.ResponseWriter, _ *http.Request) {
	resp := controllers.ScriptResponse{
		ErrorResponse: api.ErrorResponse{},
		Transaction: &core.ExpandedTransaction{
			Transaction: core.Transaction{
				TransactionData: core.TransactionData{
					Postings: core.Postings{
						{
							Source:      "world",
							Destination: "alice",
							Amount:      core.NewMonetaryInt(1),
							Asset:       "EUR",
						},
					},
					Timestamp: time.Now(),
				},
			},
			PreCommitVolumes:  nil,
			PostCommitVolumes: nil,
		},
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Printf("ERR:%s\n", err)
	}
}
