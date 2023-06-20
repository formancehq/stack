package searchengine

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/formancehq/stack/libs/go-libs/httpclient"
	"github.com/numary/ledger/pkg/core"
	goOpensearch "github.com/opensearch-project/opensearch-go"
	"github.com/ory/dockertest/v3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var (
	engine           *DefaultEngine
	openSearchClient *goOpensearch.Client
)

type testCase struct {
	name string
	fn   func(t *testing.T)
}

var tests = []testCase{
	{
		name: "nominal",
		fn:   testEngine,
	},
	{
		name: "all-fields",
		fn:   testMatchingAllFields,
	},
	{
		name: "pagination",
		fn:   testPagination,
	},
	{
		name: "specific-field",
		fn:   testMatchingSpecificField,
	},
	{
		name: "assets",
		fn:   testAssetDecimals,
	},
	{
		name: "search-in-transaction-metadata",
		fn:   testSearchInTransactionMetadata,
	},
	{
		name: "keep-only-last-document",
		fn:   testKeepOnlyLastDocument,
	},
	{
		name: "using-or-policy",
		fn:   testUsingOrPolicy,
	},
	{
		name: "sort",
		fn:   testSort,
	},
}

func indexName(t *testing.T) string {
	return strings.Split(t.Name(), "/")[1]
}

func insertESDocument(t *testing.T, id string, pipeline string, doc map[string]interface{}) {
	data, err := json.Marshal(doc)
	require.NoError(t, err)

	index := indexName(t)
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Refresh:    "true",
		Body:       bytes.NewReader(data),
		Pipeline:   pipeline,
	}
	res, err := req.Do(context.Background(), openSearchClient)
	require.NoError(t, err)
	defer res.Body.Close()

	if res.IsError() {
		require.FailNowf(t, "error inserting es", "Error inserting es index: %s [%d]", res.Status(), res.String())
	}
}

func insertTransaction(t *testing.T, ledgerName, id string, when time.Time, transaction core.Transaction) {
	insertESDocument(t, id, "TRANSACTION", map[string]interface{}{
		"kind":   "TRANSACTION",
		"ledger": ledgerName,
		"when":   when,
		"data":   transaction,
		"stack":  "",
	})
}

func insertAccount(t *testing.T, ledgerName, id string, when time.Time, payload core.Account) {
	insertESDocument(t, id, "ACCOUNT", map[string]interface{}{
		"kind":   "ACCOUNT",
		"ledger": ledgerName,
		"when":   when,
		"data":   payload,
		"stack":  "",
	})
}

func TestSearchEngine(t *testing.T) {

	if testing.Verbose() {
		logrus.StandardLogger().Level = logrus.DebugLevel
	}
	logrus.Debugln("starting opensearch container")

	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

	resource, err := pool.Run("opensearchproject/opensearch", "1.2.3", []string{
		"discovery.type=single-node",
		"DISABLE_SECURITY_PLUGIN=true",
		"DISABLE_INSTALL_DEMO_CONFIG=true",
	})
	require.NoError(t, err)

	defer func() {
		err := pool.Purge(resource)
		require.NoError(t, err)
	}()

	esAddress := "http://localhost:" + resource.GetPort("9200/tcp")
	openSearchClient, err = goOpensearch.NewClient(goOpensearch.Config{
		Addresses: []string{esAddress},
		Transport: httpclient.NewDebugHTTPTransport(&http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}),
	})
	require.NoError(t, err)

	err = pool.Retry(func() error {
		_, err = openSearchClient.Ping()
		return err
	})
	require.NoError(t, err)

	pipelineDir := "../../tests/pipelines"
	dir, err := os.ReadDir(pipelineDir)
	require.NoError(t, err)

	for _, pipelineFile := range dir {
		filename := pipelineFile.Name()
		objectType := strings.TrimSuffix(filename, ".json")
		data, err := os.ReadFile(path.Join(pipelineDir, filename))
		require.NoError(t, err)

		rsp, err := openSearchClient.Ingest.PutPipeline(objectType, bytes.NewBuffer(data))
		require.NoError(t, err)
		require.False(t, rsp.IsError())
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.NoError(t, CreateIndex(context.TODO(), openSearchClient, test.name))
			engine = NewDefaultEngine(openSearchClient, "", WithESIndex(test.name))
			test.fn(t)
		})
	}
}
