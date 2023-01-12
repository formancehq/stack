package stripe

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/formancehq/payments/internal/pkg/fifo"

	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/go-libs/logging/logginglogrus"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v72"
)

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Verbose() {
		l := logrus.New()
		l.Level = logrus.DebugLevel
		logging.SetFactory(logging.StaticLoggerFactory(logginglogrus.New(l)))
	}

	os.Exit(m.Run())
}

type ClientMockExpectation struct {
	query   url.Values
	hasMore bool
	items   []*stripe.BalanceTransaction
}

func (e *ClientMockExpectation) QueryParam(key string, value any) *ClientMockExpectation {
	var qpvalue string
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		qpvalue = fmt.Sprintf("%d", value)
	default:
		qpvalue = fmt.Sprintf("%s", value)
	}
	e.query.Set(key, qpvalue)

	return e
}

func (e *ClientMockExpectation) StartingAfter(v string) *ClientMockExpectation {
	e.QueryParam("starting_after", v)

	return e
}

func (e *ClientMockExpectation) CreatedLte(v time.Time) *ClientMockExpectation {
	e.QueryParam("created[lte]", v.Unix())

	return e
}

func (e *ClientMockExpectation) Limit(v int) *ClientMockExpectation {
	e.QueryParam("limit", v)

	return e
}

func (e *ClientMockExpectation) RespondsWith(hasMore bool,
	txs ...*stripe.BalanceTransaction,
) *ClientMockExpectation {
	e.hasMore = hasMore
	e.items = txs

	return e
}

func (e *ClientMockExpectation) handle(options ...ClientOption) ([]*stripe.BalanceTransaction, bool, error) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	for _, option := range options {
		option.apply(req)
	}

	for key := range e.query {
		if req.URL.Query().Get(key) != e.query.Get(key) {
			return nil, false, fmt.Errorf("mismatch query params, expected query param '%s' "+
				"with value '%s', got '%s'", key, e.query.Get(key), req.URL.Query().Get(key))
		}
	}

	return e.items, e.hasMore, nil
}

type ClientMock struct {
	expectations *fifo.FIFO[*ClientMockExpectation]
}

func (m *ClientMock) ForAccount(account string) Client {
	return m
}

func (m *ClientMock) BalanceTransactions(ctx context.Context,
	options ...ClientOption,
) ([]*stripe.BalanceTransaction, bool, error) {
	e, ok := m.expectations.Pop()
	if !ok {
		return nil, false, fmt.Errorf("no more expectation")
	}

	return e.handle(options...)
}

func (m *ClientMock) Expect() *ClientMockExpectation {
	e := &ClientMockExpectation{
		query: url.Values{},
	}
	m.expectations.Push(e)

	return e
}

func NewClientMock(t *testing.T, expectationsShouldBeConsumed bool) *ClientMock {
	t.Helper()

	m := &ClientMock{
		expectations: &fifo.FIFO[*ClientMockExpectation]{},
	}

	if expectationsShouldBeConsumed {
		t.Cleanup(func() {
			if !m.expectations.Empty() && !t.Failed() {
				t.Errorf("all expectations not consumed")
			}
		})
	}

	return m
}

var _ Client = &ClientMock{}
