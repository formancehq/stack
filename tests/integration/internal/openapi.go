package internal

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"

	formance "github.com/formancehq/formance-sdk-go"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var sdkClient *formance.Formance

type openapiCheckerRoundTripper struct {
	router     routers.Router
	underlying http.RoundTripper
}

func (c *openapiCheckerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	route, pathParams, err := c.router.FindRoute(req)
	Expect(errors.Wrapf(err, "retrieving operation for route %s %s", req.Method, req.URL.String())).
		WithOffset(6).To(Succeed())

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
		WithOffset(6).To(Succeed())

	_, err = httputil.DumpRequest(req, true)
	Expect(err).ToNot(HaveOccurred())

	rsp, err := c.underlying.RoundTrip(req)
	Expect(err).WithOffset(6).To(Succeed())

	data, err := io.ReadAll(rsp.Body)
	Expect(err).WithOffset(6).To(Succeed())

	rsp.Body = io.NopCloser(bytes.NewBuffer(data))

	err = openapi3filter.ValidateResponse(req.Context(), &openapi3filter.ResponseValidationInput{
		RequestValidationInput: input,
		Status:                 rsp.StatusCode,
		Header:                 rsp.Header,
		Body:                   io.NopCloser(bytes.NewBuffer(data)),
		Options:                options,
	})
	Expect(err).WithOffset(6).To(Succeed())

	return rsp, nil
}

var _ http.RoundTripper = &openapiCheckerRoundTripper{}

func newOpenapiCheckerTransport(ctx context.Context, rt http.RoundTripper) (*openapiCheckerRoundTripper, error) {
	openapiRawSpec, err := os.ReadFile(filepath.Join("..", "..", "..", "openapi", "build", "generate.json"))
	if err != nil {
		return nil, err
	}

	loader := &openapi3.Loader{
		Context:               ctx,
		IsExternalRefsAllowed: true,
	}
	doc, err := loader.LoadFromData(openapiRawSpec)
	if err != nil {
		return nil, err
	}

	// Override default servers
	doc.Servers = []*openapi3.Server{{
		URL: "http://127.0.0.1",
	}}
	doc.Security = openapi3.SecurityRequirements{}
	doc.Components.SecuritySchemes = openapi3.SecuritySchemes{}

	err = doc.Validate(ctx)
	if err != nil {
		return nil, err
	}

	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		return nil, err
	}

	return &openapiCheckerRoundTripper{
		router:     router,
		underlying: rt,
	}, nil
}
