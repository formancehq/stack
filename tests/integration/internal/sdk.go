package internal

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/stack/libs/go-libs/httpclient"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var sdkClient *formance.APIClient

type openapiCheckerRoundTripper struct {
	router     routers.Router
	underlying http.RoundTripper
}

func (c *openapiCheckerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	route, pathParams, err := c.router.FindRoute(req)
	Expect(errors.Wrapf(err, "retrieving operation for route %s %s", req.Method, req.URL.String())).
		WithOffset(8).To(Succeed())

	options := &openapi3filter.Options{
		IncludeResponseStatus: true,
		MultiError:            true,
		AuthenticationFunc:    openapi3filter.NoopAuthenticationFunc,
	}
	input := &openapi3filter.RequestValidationInput{
		Request:     req,
		PathParams:  pathParams,
		QueryParams: req.URL.Query(),
		Route:       route,
		Options:     options,
	}

	Expect(errors.Wrap(openapi3filter.ValidateRequest(req.Context(), input), "validating request")).
		WithOffset(8).To(Succeed())

	_, err = httputil.DumpRequest(req, true)
	Expect(err).NotTo(HaveOccurred())

	rsp, err := c.underlying.RoundTrip(req)
	Expect(err).WithOffset(8).To(Succeed())

	data, err := io.ReadAll(rsp.Body)
	Expect(err).WithOffset(8).To(Succeed())
	rsp.Body = io.NopCloser(bytes.NewBuffer(data))

	err = openapi3filter.ValidateResponse(req.Context(), &openapi3filter.ResponseValidationInput{
		RequestValidationInput: input,
		Status:                 rsp.StatusCode,
		Header:                 rsp.Header,
		Body:                   io.NopCloser(bytes.NewBuffer(data)),
		Options:                options,
	})
	Expect(errors.Wrap(err, "validating response")).WithOffset(8).To(Succeed())

	return rsp, nil
}

var _ http.RoundTripper = &openapiCheckerRoundTripper{}

func newOpenapiCheckerTransport(rt http.RoundTripper) *openapiCheckerRoundTripper {
	openapiRawSpec, err := os.ReadFile(filepath.Join("..", "..", "..", "openapi", "build", "generate.json"))
	Expect(err).NotTo(HaveOccurred())

	loader := &openapi3.Loader{
		Context:               TestContext(),
		IsExternalRefsAllowed: true,
	}
	doc, err := loader.LoadFromData(openapiRawSpec)
	Expect(err).NotTo(HaveOccurred())

	// Override default servers
	doc.Servers = []*openapi3.Server{{
		URL: "http://127.0.0.1",
	}}
	doc.Security = openapi3.SecurityRequirements{}
	doc.Components.SecuritySchemes = openapi3.SecuritySchemes{}

	err = doc.Validate(ctx)
	Expect(err).NotTo(HaveOccurred())

	router, err := gorillamux.NewRouter(doc)
	Expect(err).NotTo(HaveOccurred())

	return &openapiCheckerRoundTripper{
		router:     router,
		underlying: rt,
	}
}

func configureSDK() {
	gatewayUrl, err := url.Parse(gatewayServer.URL)
	if err != nil {
		panic(err)
	}

	configuration := formance.NewConfiguration()
	configuration.Host = gatewayUrl.Host
	configuration.Servers = []formance.ServerConfiguration{{
		URL: gatewayUrl.String(),
	}}
	configuration.HTTPClient = &http.Client{
		Transport: newOpenapiCheckerTransport(
			httpclient.NewDebugHTTPTransport(http.DefaultTransport),
		),
	}
	sdkClient = formance.NewAPIClient(configuration)
}

func Client() *formance.APIClient {
	return sdkClient
}
