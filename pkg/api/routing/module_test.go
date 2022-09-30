package routing

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func testWithUrl(t *testing.T, urlStr string) {
	u, err := url.Parse(urlStr)
	require.NoError(t, err)

	ctx := NewContext(context.Background())

	app := fx.New(
		fx.Invoke(fx.Annotate(func(router *mux.Router) {
			router.Path("/subpath-with-prefix").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})
		}, fx.ParamTags(`name:"prefixedRouter"`))),
		fx.Invoke(fx.Annotate(func(router *mux.Router) {
			router.Path("/subpath-to-root").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})
		}, fx.ParamTags(`name:"rootRouter"`))),
		Module(":0", u),
		fx.NopLogger,
	)

	require.NoError(t, app.Start(ctx))
	defer func() {
		require.NoError(t, app.Stop(ctx))
	}()

	serverUrl := fmt.Sprintf("http://localhost:%d", ListeningPort(ctx))
	serverUrlWithPath := fmt.Sprintf("%s%s", serverUrl, u.Path)

	rsp, err := http.Get(fmt.Sprintf("%s/_healthcheck", serverUrlWithPath))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rsp.StatusCode)

	rsp, err = http.Get(fmt.Sprintf("%s/_healthcheck", serverUrl))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rsp.StatusCode)

	rsp, err = http.Get(fmt.Sprintf("%s/subpath-with-prefix", serverUrlWithPath))
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, rsp.StatusCode)

	rsp, err = http.Get(fmt.Sprintf("%s/subpath-to-root", serverUrl))
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, rsp.StatusCode)
}

func TestModule(t *testing.T) {
	testWithUrl(t, "http://localhost")
	testWithUrl(t, "http://localhost/any/sub/path")
}
