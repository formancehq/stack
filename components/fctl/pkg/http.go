package fctl

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/TylerBrock/colorjson"
)

func GetHttpClient(flags *flag.FlagSet, defaultHeaders map[string][]string, out io.Writer) *http.Client {
	return NewHTTPClient(
		GetBool(flags, InsecureTlsFlag),
		GetBool(flags, DebugFlag),
		defaultHeaders,
		out,
	)
}

type RoundTripperFn func(req *http.Request) (*http.Response, error)

func (fn RoundTripperFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func printBody(out io.Writer, data []byte) {
	if len(data) == 0 {
		return
	}
	raw := make(map[string]any)
	if err := json.Unmarshal(data, &raw); err == nil {
		f := colorjson.NewFormatter()
		f.Indent = 2
		colorized, err := f.Marshal(raw)
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(out, string(colorized))
	} else {
		fmt.Fprintln(out, string(data))
	}
}

func debugRoundTripper(rt http.RoundTripper, out io.Writer) RoundTripperFn {
	return func(req *http.Request) (*http.Response, error) {
		data, err := httputil.DumpRequest(req, false)
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(out, string(data))

		if req.Body != nil {
			data, err = io.ReadAll(req.Body)
			if err != nil {
				panic(err)
			}
			req.Body.Close()
			req.Body = io.NopCloser(bytes.NewBuffer(data))
			printBody(out, data)
		}

		rsp, err := rt.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		data, err = httputil.DumpResponse(rsp, false)
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(out, string(data))

		if rsp.Body != nil {
			data, err = io.ReadAll(rsp.Body)
			if err != nil {
				panic(err)
			}
			rsp.Body.Close()
			rsp.Body = io.NopCloser(bytes.NewBuffer(data))
			printBody(out, data)
		}

		return rsp, nil
	}
}

func defaultHeadersRoundTripper(rt http.RoundTripper, headers map[string][]string) RoundTripperFn {
	return func(req *http.Request) (*http.Response, error) {
		for k, v := range headers {
			for _, vv := range v {
				req.Header.Add(k, vv)
			}
		}
		return rt.RoundTrip(req)
	}
}

func NewHTTPClient(insecureTLS bool, debug bool, defaultHeaders map[string][]string, out io.Writer) *http.Client {
	var transport http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecureTLS,
		},
	}
	if debug {
		transport = debugRoundTripper(transport, out)
	}
	if len(defaultHeaders) > 0 {
		transport = defaultHeadersRoundTripper(transport, defaultHeaders)
	}
	return &http.Client{
		Transport: transport,
	}
}
