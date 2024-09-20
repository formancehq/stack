package internal

import (
	"context"
	"fmt"
	"github.com/formancehq/go-libs/collectionutils"
	"github.com/formancehq/go-libs/httpclient"
	"github.com/xo/dburl"
	"net/http"
	"net/http/httptest"
	"os"
)

type routing struct {
	port        uint16
	routingFunc func(path, method string) bool
}

type Test struct {
	env             *Env
	id              string
	loadedModules   collectionutils.Set[string]
	servicesToRoute map[string][]routing
	httpServer      *httptest.Server
	httpTransport   http.RoundTripper
}

func (test *Test) setup() error {
	gateway := newGateway(test)
	test.httpServer = httptest.NewServer(gateway)
	test.httpTransport = httpclient.NewDebugHTTPTransport(http.DefaultTransport)

	return nil
}

func (test *Test) tearDown() error {
	if test.httpServer != nil {
		test.httpServer.Close()
	}
	return nil
}

func (test *Test) loadModule(ctx context.Context, m *Module) error {
	if test.loadedModules.Contains(m.name) {
		return nil
	}

	if m.createDatabase {
		if err := test.createDatabase(ctx, m.name); err != nil {
			return err
		}
	}

	for _, s := range m.services {
		if err := s.load(ctx, test); err != nil {
			return err
		}
	}

	test.loadedModules.Put(m.name)

	return nil
}

func (test *Test) unloadModule(ctx context.Context, m *Module) error {

	for _, s := range m.services {
		if err := s.unload(ctx, test); err != nil {
			return err
		}
	}

	if m.createDatabase {
		if err := test.dropDatabase(ctx, m.name); err != nil {
			return err
		}
	}
	return nil
}

func (test *Test) createDatabase(ctx context.Context, name string) error {
	_, err := test.env.sqlConn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s-%s";`, test.id, name))
	return err
}

func (test *Test) dropDatabase(ctx context.Context, name string) error {
	if os.Getenv("NO_CLEANUP") != "true" {
		_, err := test.env.sqlConn.Exec(ctx, fmt.Sprintf(`DROP DATABASE "%s-%s";`, test.id, name))
		return err
	}
	return nil
}

func (test *Test) registerServiceToRoute(name string, routing routing) {
	test.servicesToRoute[name] = append(test.servicesToRoute[name], routing)
}

func (test *Test) GetDatabaseSourceName(name string) string {
	dsn, err := dburl.Parse(GetPostgresDSNString())
	if err != nil {
		panic(err)
	}
	dsn.Path = fmt.Sprintf("%s-%s", test.id, name)

	return dsn.String()
}

func (test *Test) ID() string {
	return test.id
}

func (test *Test) Workdir() string {
	return test.env.workdir
}

func (test *Test) GatewayURL() string {
	return test.httpServer.URL
}
