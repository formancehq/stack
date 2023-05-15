package fctl

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/TylerBrock/colorjson"
	"github.com/spf13/cobra"
)

func GetHttpClient(cmd *cobra.Command) *http.Client {
	return NewHTTPClient(
		GetBool(cmd, InsecureTlsFlag),
		GetBool(cmd, DebugFlag),
	)
}

type RoundTripperFn func(req *http.Request) (*http.Response, error)

func (fn RoundTripperFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func printBody(data []byte) {
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
		fmt.Println(string(colorized))
	} else {
		fmt.Println(string(data))
	}
}

func debugRoundTripper(rt http.RoundTripper) RoundTripperFn {
	return func(req *http.Request) (*http.Response, error) {
		data, err := httputil.DumpRequest(req, false)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))

		if req.Body != nil {
			data, err = io.ReadAll(req.Body)
			if err != nil {
				panic(err)
			}
			req.Body.Close()
			req.Body = io.NopCloser(bytes.NewBuffer(data))
			printBody(data)
		}

		rsp, err := rt.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		data, err = httputil.DumpResponse(rsp, false)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))

		if rsp.Body != nil {
			data, err = io.ReadAll(rsp.Body)
			if err != nil {
				panic(err)
			}
			rsp.Body.Close()
			rsp.Body = io.NopCloser(bytes.NewBuffer(data))
			printBody(data)
		}

		return rsp, nil
	}
}

func NewHTTPClient(insecureTLS, debug bool) *http.Client {
	var transport http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecureTLS,
		},
	}
	if debug {
		transport = debugRoundTripper(transport)
	}
	return &http.Client{
		Transport: transport,
	}
}
