package testserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/formancehq/ingester/ingesterclient"
	"github.com/formancehq/stack/ee/ingester/internal"
	"github.com/formancehq/stack/libs/go-libs/api"
	. "github.com/formancehq/stack/libs/go-libs/testing/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
)

func UseNewTestServer(configurationProvider func() Configuration) *Deferred[*Server] {
	d := NewDeferred[*Server]()
	BeforeEach(func() {
		d.Reset()
		d.SetValue(New(GinkgoT(), configurationProvider()))
	})
	return d
}

func WithConnector(name string, connectorFactory func() Connector, fn func(p *Deferred[Connector])) {
	Context(fmt.Sprintf("with connector '%s'", name), func() {
		ret := NewDeferred[Connector]()
		BeforeEach(func() {
			ret.Reset()
			ret.SetValue(connectorFactory())
		})
		fn(ret)
	})
}

func CreateConnector(ctx context.Context, srv *Server, request ingesterclient.ConnectorConfiguration) (*ingesterclient.Connector, error) {
	response, httpResponse, err := srv.Client().ConnectorsApi.CreateConnector(ctx).Body(request).Execute()
	if err != nil {
		return nil, mapSDKError(err, httpResponse)
	}

	return &response.Data, nil
}

func MustCreateConnector(srv *Deferred[*Server], connectorConfiguration ingesterclient.ConnectorConfiguration) ingesterclient.Connector {
	GinkgoHelper()

	connector, err := CreateConnector(context.Background(), srv.GetValue(), connectorConfiguration)
	Expect(err).To(BeNil())

	return *connector
}

func DeleteConnector(ctx context.Context, srv *Server, id string) error {
	response, err := srv.Client().ConnectorsApi.DeleteConnector(ctx, id).Execute()
	return mapSDKError(err, response)
}

func MustDeleteConnector(srv *Deferred[*Server], id string) {
	GinkgoHelper()

	err := DeleteConnector(context.Background(), srv.GetValue(), id)
	Expect(err).To(BeNil())
}

func ListConnectors(ctx context.Context, srv *Server) *ingesterclient.ListConnectors200ResponseCursor {
	GinkgoHelper()

	listResponse, _, err := srv.Client().ConnectorsApi.ListConnectors(ctx).Execute()
	Expect(err).To(BeNil())

	return listResponse.Cursor
}

func GetConnector(ctx context.Context, srv *Server, id string) (*ingesterclient.Connector, error) {
	response, httpResponse, err := srv.Client().ConnectorsApi.GetConnectorState(ctx, id).Execute()
	if err != nil {
		return nil, mapSDKError(err, httpResponse)
	}

	return &response.Data, nil
}

func MustGetConnector(srv *Deferred[*Server], id string) ingesterclient.Connector {
	GinkgoHelper()

	connector, err := GetConnector(context.Background(), srv.GetValue(), id)
	Expect(err).To(BeNil())

	return *connector
}

func pipelineAction(
	ctx context.Context,
	srv *Server,
	id string,
	fn func(*ingesterclient.PipelinesApiService, context.Context, string) interface {
		Execute() (*http.Response, error)
	},
) error {
	response, err := fn(srv.Client().PipelinesApi, ctx, id).Execute()
	if err != nil {
		return mapSDKError(err, response)
	}
	return nil
}

func CreatePipeline(ctx context.Context, srv *Server, pipelineConfiguration ingesterclient.PipelineConfiguration) (*ingesterclient.Pipeline, error) {
	res, _, err := srv.Client().PipelinesApi.CreatePipeline(ctx).CreatePipelineRequest(ingesterclient.CreatePipelineRequest{
		Module:      pipelineConfiguration.Module,
		ConnectorID: pipelineConfiguration.ConnectorID,
	}).Execute()
	if err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func MustCreatePipeline(srv *Deferred[*Server], pipelineConfiguration ingesterclient.PipelineConfiguration) ingesterclient.Pipeline {
	GinkgoHelper()

	pipeline, err := CreatePipeline(context.Background(), srv.GetValue(), pipelineConfiguration)
	Expect(err).To(BeNil())
	Eventually(pipeline.Id).Should(HaveState(srv.GetValue(), ingester.StateLabelReady))

	return *pipeline
}

func PausePipeline(ctx context.Context, srv *Server, id string) error {
	return pipelineAction(ctx, srv, id, func(service *ingesterclient.PipelinesApiService, ctx context.Context, id string) interface {
		Execute() (*http.Response, error)
	} {
		return service.PausePipeline(ctx, id)
	})
}

func MustPausePipeline(srv *Deferred[*Server], id string) {
	GinkgoHelper()

	Expect(PausePipeline(context.Background(), srv.GetValue(), id)).To(Succeed())
	Eventually(id).Should(HaveState(srv.GetValue(), ingester.StateLabelPause))
}

func ResumePipeline(ctx context.Context, srv *Server, id string) error {
	return pipelineAction(ctx, srv, id, func(service *ingesterclient.PipelinesApiService, ctx context.Context, id string) interface {
		Execute() (*http.Response, error)
	} {
		return service.ResumePipeline(ctx, id)
	})
}

func MustResumePipeline(srv *Deferred[*Server], id string) {
	GinkgoHelper()

	Expect(ResumePipeline(context.Background(), srv.GetValue(), id)).To(Succeed())
	Eventually(id).Should(HaveState(srv.GetValue(), ingester.StateLabelReady))
}

func ResetPipeline(ctx context.Context, srv *Server, id string) error {
	return pipelineAction(ctx, srv, id, func(service *ingesterclient.PipelinesApiService, ctx context.Context, id string) interface {
		Execute() (*http.Response, error)
	} {
		return service.ResetPipeline(ctx, id)
	})
}

func MustResetPipeline(srv *Deferred[*Server], id string) {
	GinkgoHelper()

	Expect(ResetPipeline(context.Background(), srv.GetValue(), id)).To(Succeed())
	Eventually(id).Should(HaveState(srv.GetValue(), ingester.StateLabelReady))
}

func StopPipeline(ctx context.Context, srv *Server, id string) error {
	return pipelineAction(ctx, srv, id, func(service *ingesterclient.PipelinesApiService, ctx context.Context, id string) interface {
		Execute() (*http.Response, error)
	} {
		return service.StopPipeline(ctx, id)
	})
}

func MustStopPipeline(srv *Deferred[*Server], id string) {
	GinkgoHelper()

	Expect(StopPipeline(context.Background(), srv.GetValue(), id)).To(Succeed())
	Eventually(id).Should(HaveState(srv.GetValue(), ingester.StateLabelStop))
}

func StartPipeline(ctx context.Context, srv *Server, id string) error {
	return pipelineAction(ctx, srv, id, func(service *ingesterclient.PipelinesApiService, ctx context.Context, id string) interface {
		Execute() (*http.Response, error)
	} {
		return service.StartPipeline(ctx, id)
	})
}

func MustStartPipeline(srv *Deferred[*Server], id string) {
	GinkgoHelper()

	Expect(StartPipeline(context.Background(), srv.GetValue(), id)).To(Succeed())
	Eventually(id).Should(HaveState(srv.GetValue(), ingester.StateLabelReady))
}

func DeletePipeline(ctx context.Context, srv *Server, id string) error {
	response, err := srv.Client().PipelinesApi.DeletePipeline(ctx, id).Execute()
	return mapSDKError(err, response)
}

func ListPipelines(ctx context.Context, srv *Server) *ingesterclient.ListPipelines200ResponseCursor {
	listResponse, _, err := srv.Client().PipelinesApi.ListPipelines(ctx).Execute()
	Expect(err).To(BeNil())

	return listResponse.Cursor
}

func GetPipeline(ctx context.Context, srv *Server, id string) (*ingesterclient.Pipeline, error) {
	response, httpResponse, err := srv.Client().PipelinesApi.GetPipelineState(ctx, id).Execute()
	if err != nil {
		return nil, mapSDKError(err, httpResponse)
	}

	return &response.Data, nil
}

func MustGetPipeline(srv *Deferred[*Server], id string) ingesterclient.Pipeline {
	GinkgoHelper()

	pipeline, err := GetPipeline(context.Background(), srv.GetValue(), id)
	Expect(err).To(BeNil())

	return *pipeline
}

func mapSDKError(err error, response *http.Response) error {
	switch err.(type) {
	case *ingesterclient.GenericOpenAPIError:
		errorResponse := api.ErrorResponse{}
		if err := json.NewDecoder(response.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf(
				"got unexpected status code %d (unable to parse the body: %s) ",
				response.StatusCode,
				err,
			)
		}
		return errorResponse
	default:
		return err
	}
}
