package benchmarks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/formancehq/ledger/pkg/api/controllers"
	"github.com/formancehq/ledger/pkg/api/routes"
	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/ledger/pkg/ledger"
	"github.com/formancehq/ledger/pkg/opentelemetry/metrics"
	"github.com/formancehq/ledger/pkg/storage/sqlstorage/sqlstoragetesting"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func BenchmarkParallelWrites(b *testing.B) {

	driver := sqlstoragetesting.StorageDriver(b)
	resolver := ledger.NewResolver(driver)
	b.Cleanup(func() {
		require.NoError(b, resolver.CloseLedgers(context.Background()))
	})

	ledgerName := uuid.NewString()

	backend := controllers.NewDefaultBackend(driver, "latest", resolver)
	router := routes.NewRouter(backend, nil, nil, metrics.NewNoOpMetricsRegistry())
	srv := httptest.NewServer(router)
	defer srv.Close()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.SetParallelism(1000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		buf := bytes.NewBufferString("")
		for pb.Next() {
			buf.Reset()

			err := json.NewEncoder(buf).Encode(controllers.PostTransactionRequest{
				Script: core.Script{
					Plain: fmt.Sprintf(`send [USD/2 100] (
						source = @world
						destination = @accounts:%d
					)`, r.Int()%100),
				},
			})
			require.NoError(b, err)

			req := httptest.NewRequest("POST", "/"+ledgerName+"/transactions", buf)
			req.URL.RawQuery = url.Values{
				"async": []string{os.Getenv("ASYNC")},
			}.Encode()
			rsp := httptest.NewRecorder()

			router.ServeHTTP(rsp, req)

			require.Equal(b, http.StatusOK, rsp.Code)
		}
	})
	b.StopTimer()
}
