package elastictesting

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/testing/docker"
	"github.com/olivere/elastic/v7"
	"github.com/ory/dockertest/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type Server struct {
	elasticsearchEndpoint string
	t                     docker.T
}

func (s *Server) Endpoint() string {
	return s.elasticsearchEndpoint
}

func (s *Server) NewClient() *elastic.Client {
	ret, err := elastic.NewClient(elastic.SetURL(s.elasticsearchEndpoint))
	require.NoError(s.t, err)
	return ret
}

func CreateServer(pool *docker.Pool) *Server {

	resource := pool.Run(docker.Configuration{
		RunOptions: &dockertest.RunOptions{
			Repository: "elasticsearch",
			Tag:        "8.14.3",
			Env: []string{
				"discovery.type=single-node",
				"xpack.security.enabled=false",
				"xpack.security.enrollment.enabled=false",
			},
		},
		CheckFn: func(ctx context.Context, resource *dockertest.Resource) error {
			client, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:" + resource.GetPort("9200/tcp")))
			if err != nil {
				return errors.Wrap(err, "connecting to server")
			}
			client.Stop()
			return nil
		},
	})

	return &Server{
		t:                     pool.T(),
		elasticsearchEndpoint: "http://127.0.0.1:" + resource.GetPort("9200/tcp"),
	}
}
