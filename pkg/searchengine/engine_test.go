package searchengine

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/aquasecurity/esquery"
	"github.com/numary/ledger/pkg/core"
	"github.com/stretchr/testify/assert"
)

func testEngine(t *testing.T) {
	ledger := "quickstart"
	insertTransaction(t, ledger, "transaction0", time.Now(), core.Transaction{
		Metadata: core.Metadata{
			"foo": json.RawMessage(`{"foo": "bar"}`),
		},
		Postings: core.Postings{
			{
				Source:      "world",
				Destination: "central:bank",
				Asset:       "USD",
				Amount:      100,
			},
		},
	})

	q := NewSingleDocTypeSearch("TRANSACTION")
	q.WithLedgers(ledger)
	q.WithTerms("central:bank")
	response, err := q.Do(context.Background(), engine)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Len(t, response.Items, 1) {
		return
	}
}

func testMatchingAllFields(t *testing.T) {

	now := time.Now().Round(time.Second).UTC()
	insertTransaction(t, "quickstart", "transaction0", now.Add(-time.Minute), core.Transaction{
		Postings: core.Postings{
			{
				Source:      "world",
				Destination: "central_bank",
				Amount:      100,
				Asset:       "USD",
			},
		},
	})
	insertTransaction(t, "quickstart", "transaction1", now, core.Transaction{})
	insertTransaction(t, "quickstart", "transaction2", now.Add(time.Minute), core.Transaction{})

	q := NewMultiDocTypeSearch()
	q.WithLedgers("quickstart")
	q.WithTerms("USD")

	response, err := q.Do(context.Background(), engine)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Len(t, response["TRANSACTION"], 1) {
		return
	}

	q = NewMultiDocTypeSearch()
	q.WithLedgers("quickstart")
	q.WithTerms("US")
	response, err = q.Do(context.Background(), engine)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Len(t, response["TRANSACTION"], 1) {
		return
	}
}

func testSort(t *testing.T) {
	now := time.Now().Round(time.Second).UTC()
	const count = 20
	for i := 0; i < count; i++ {
		insertTransaction(t, "quickstart", fmt.Sprintf("transaction%d", i), now.Add(time.Duration(i)*time.Minute), core.Transaction{})
	}

	q := NewSingleDocTypeSearch("TRANSACTION")
	q.WithLedgers("quickstart")
	q.WithSize(20)
	q.WithSort("txid", esquery.OrderAsc)

	_, err := openSearchClient.Indices.GetMapping()
	if !assert.NoError(t, err) {
		return
	}

	response, err := q.Do(context.Background(), engine)
	if !assert.NoError(t, err) {
		return
	}

	if !assert.Len(t, response.Items, count) {
		return
	}
}

func testPagination(t *testing.T) {
	now := time.Now().Round(time.Second).UTC()

	for i := 0; i < 20; i++ {
		at := now.Add(time.Duration(i) * time.Minute)
		insertTransaction(t, "quickstart", fmt.Sprintf("transaction%d", i), at, core.Transaction{
			Timestamp: at.Format(time.RFC3339),
		})
	}

	var (
		searchAfter []interface{}
	)
	for i := 0; ; i++ {
		q := NewSingleDocTypeSearch("TRANSACTION")
		q.WithLedgers("quickstart")
		q.WithSize(5)
		q.WithSort("timestamp", esquery.OrderDesc)
		q.WithSearchAfter(searchAfter)

		_, err := openSearchClient.Indices.GetMapping()
		if !assert.NoError(t, err) {
			return
		}

		response, err := q.Do(context.Background(), engine)
		if !assert.NoError(t, err) {
			return
		}

		tx := core.Transaction{}
		if !assert.NoError(t, json.Unmarshal(response.Items[0], &tx)) {
			return
		}

		if i < 3 {
			if !assert.Len(t, response.Items, 5) {
				return
			}
			if !assert.Equal(t, tx.Timestamp, now.Add(19*time.Minute).Add(-time.Duration(i)*5*time.Minute).UTC().Format(time.RFC3339)) {
				return
			}
		} else {
			if !assert.Len(t, response.Items, 5) {
				return
			}
			if !assert.Equal(t, tx.Timestamp, now.Add(19*time.Minute).Add(-time.Duration(i)*5*time.Minute).UTC().Format(time.RFC3339)) {
				return
			}
			break
		}
		lastTx := core.Transaction{}
		if !assert.NoError(t, json.Unmarshal(response.Items[4], &lastTx)) {
			return
		}

		searchAfter = []interface{}{lastTx.Timestamp}
	}

}

func testMatchingSpecificField(t *testing.T) {

	now := time.Now().Round(time.Second).UTC()
	insertTransaction(t, "quickstart", "transaction0", now.Add(-time.Minute), core.Transaction{
		Timestamp: now.Add(-time.Minute).Format(time.RFC3339),
		Postings: core.Postings{
			{
				Source:      "world",
				Destination: "central_bank",
				Amount:      100,
				Asset:       "USD",
			},
		},
	})
	insertTransaction(t, "quickstart", "transaction1", now.Add(time.Minute), core.Transaction{
		Timestamp: now.Add(time.Minute).Format(time.RFC3339),
		Postings: core.Postings{
			{
				Source:      "central_bank",
				Destination: "user:001",
				Amount:      1000,
				Asset:       "USD",
			},
			{
				Source:      "world",
				Destination: "central_bank",
				Amount:      10000,
				Asset:       "USD",
			},
		},
	})

	type testCase struct {
		name          string
		term          string
		expectedCount int
	}

	testCases := []testCase{
		{
			name:          "equality-using-equal",
			term:          "amount=100",
			expectedCount: 1,
		},
		{
			name:          "greater-than-on-long",
			term:          "amount>500",
			expectedCount: 1,
		},
		{
			name:          "greater-than-on-date-millis",
			term:          fmt.Sprintf("timestamp>%d", now.UnixMilli()),
			expectedCount: 1,
		},
		{
			name:          "greater-than-on-date-rfc3339",
			term:          fmt.Sprintf("timestamp>%s", now.Format(time.RFC3339)),
			expectedCount: 1,
		},
		{
			name:          "lower-than",
			term:          "amount<5000",
			expectedCount: 2,
		},
		{
			name:          "greater-than-or-equal",
			term:          "amount>=1000",
			expectedCount: 1,
		},
		{
			name:          "lower-than-or-equal",
			term:          "amount<=100",
			expectedCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			q := NewMultiDocTypeSearch()
			q.WithLedgers("quickstart")
			q.WithTerms(tc.term)

			response, err := q.Do(context.Background(), engine)
			if !assert.NoError(t, err) {
				return
			}
			if !assert.Len(t, response["TRANSACTION"], tc.expectedCount) {
				return
			}
		})
	}
}

func testUsingOrPolicy(t *testing.T) {

	now := time.Now().Round(time.Second).UTC()
	insertTransaction(t, "quickstart", "transaction0", now.Add(-time.Minute), core.Transaction{
		Postings: core.Postings{
			{
				Source:      "world",
				Destination: "central_bank1",
				Amount:      100,
				Asset:       "USD",
			},
		},
	})
	insertTransaction(t, "quickstart", "transaction1", now.Add(time.Minute), core.Transaction{
		Postings: core.Postings{
			{
				Source:      "world",
				Destination: "central_bank2",
				Amount:      1000,
				Asset:       "USD",
			},
		},
	})
	insertTransaction(t, "quickstart", "transaction2", now.Add(time.Minute), core.Transaction{
		Postings: core.Postings{
			{
				Source:      "world",
				Destination: "central_bank3",
				Amount:      1000,
				Asset:       "USD",
			},
		},
	})

	q := NewSingleDocTypeSearch("TRANSACTION")
	q.WithLedgers("quickstart")
	q.WithTerms("destination=central_bank1", "destination=central_bank2")
	q.WithPolicy(TermPolicyOR)

	response, err := q.Do(context.Background(), engine)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Len(t, response.Items, 2) {
		return
	}
}

func testAssetDecimals(t *testing.T) {

	now := time.Now().Round(time.Second).UTC()
	insertTransaction(t, "quickstart", "transaction0", now.Add(-time.Minute), core.Transaction{
		Postings: core.Postings{
			{
				Source:      "world",
				Destination: "central_bank",
				Amount:      10045,
				Asset:       "USD/2",
			},
		},
	})
	insertTransaction(t, "quickstart", "transaction1", now.Add(-time.Minute), core.Transaction{
		Postings: core.Postings{
			{
				Source:      "world",
				Destination: "central_bank",
				Amount:      1000,
				Asset:       "USD",
			},
		},
	})

	type testCase struct {
		name          string
		term          string
		expectedCount int
	}

	testCases := []testCase{
		{
			name:          "colon",
			term:          "amount=100.45",
			expectedCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			q := NewMultiDocTypeSearch()
			q.WithTerms(tc.term)
			q.WithLedgers("quickstart")

			response, err := q.Do(context.Background(), engine)
			if !assert.NoError(t, err) {
				return
			}
			if !assert.Len(t, response["TRANSACTION"], tc.expectedCount) {
				return
			}
		})
	}

}

func testSearchInTransactionMetadata(t *testing.T) {
	now := time.Now().Round(time.Second).UTC()
	metadata := core.Metadata{
		"Hello": json.RawMessage("\"guys!\""),
		"John":  json.RawMessage("\"Snow!\""),
	}
	insertTransaction(t, "quickstart", "transaction0", now, core.Transaction{
		Metadata: metadata,
	})
	insertTransaction(t, "quickstart", "transaction1", now, core.Transaction{})

	q := NewMultiDocTypeSearch()
	q.WithTerms("John")
	response, err := q.Do(context.Background(), engine)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Len(t, response["TRANSACTION"], 1) {
		return
	}

	tx := core.Transaction{}
	if !assert.NoError(t, json.Unmarshal(response["TRANSACTION"][0], &tx)) {
		return
	}
	if !assert.Equal(t, metadata, tx.Metadata) {
		return
	}
}

func testKeepOnlyLastDocument(t *testing.T) {

	now := time.Now().Round(time.Hour)
	for i := 0; i < 10; i++ {
		insertAccount(t, "quickstart", fmt.Sprintf("account%d", i), now, core.Account{
			Address: fmt.Sprintf("user:00%d", i),
		})
	}
	for i := 0; i < 20; i++ {
		insertTransaction(t, "quickstart", fmt.Sprintf("transaction%d", i), now.Add(2*time.Minute), core.Transaction{
			ID:        int64(i),
			Timestamp: now.Add(time.Hour).Format(time.RFC3339),
		})
	}

	q := NewMultiDocTypeSearch()
	q.WithSize(5)

	response, err := q.Do(context.Background(), engine)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Len(t, response["TRANSACTION"], 5) {
		return
	}
	if !assert.Len(t, response["ACCOUNT"], 5) {
		return
	}
}
