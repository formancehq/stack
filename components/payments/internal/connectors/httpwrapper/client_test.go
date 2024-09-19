package httpwrapper_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}

type successRes struct {
	ID string `json:"id"`
}

type errorRes struct {
	Code string `json:"code"`
}

var _ = Describe("ClientWrapper", func() {
	var (
		config *httpwrapper.Config
		client httpwrapper.Client
		server *httptest.Server
	)

	BeforeEach(func() {
		config = &httpwrapper.Config{Timeout: 30 * time.Millisecond}
		var err error
		client, err = httpwrapper.NewClient(config)
		Expect(err).To(BeNil())
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params, err := url.ParseQuery(r.URL.RawQuery)
			Expect(err).To(BeNil())

			code := params.Get("code")
			statusCode, err := strconv.Atoi(code)
			Expect(err).To(BeNil())
			if statusCode == http.StatusOK {
				w.Write([]byte(`{"id":"someid"}`))
				return
			}

			w.WriteHeader(statusCode)
			w.Write([]byte(`{"code":"err123"}`))
		}))
	})
	AfterEach(func() {
		server.Close()
	})

	Context("making a request with default client settings", func() {
		It("unmarshals successful responses when acceptable status code seen", func(ctx SpecContext) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"?code=200", http.NoBody)
			Expect(err).To(BeNil())

			res := &successRes{}
			code, doErr := client.Do(req, res, nil)
			Expect(code).To(Equal(http.StatusOK))
			Expect(doErr).To(BeNil())
			Expect(res.ID).To(Equal("someid"))
		})
		It("unmarshals error responses when bad status code seen", func(ctx SpecContext) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"?code=500", http.NoBody)
			Expect(err).To(BeNil())

			res := &errorRes{}
			code, doErr := client.Do(req, &successRes{}, res)
			Expect(code).To(Equal(http.StatusInternalServerError))
			Expect(doErr).To(MatchError(httpwrapper.ErrStatusCodeUnexpected))
			Expect(res.Code).To(Equal("err123"))
		})
		It("responds with error when HTTP request fails", func(ctx SpecContext) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "notaurl", http.NoBody)
			Expect(err).To(BeNil())

			res := &errorRes{}
			code, doErr := client.Do(req, &successRes{}, res)
			Expect(code).To(Equal(0))
			Expect(doErr).To(MatchError(ContainSubstring("failed to make request")))
		})
	})
})
