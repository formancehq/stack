/*
Formance Stack API

Testing OrchestrationApiService

*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech);

package formance

import (
    "context"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "testing"
    client "./openapi"
)

func Test_formance_OrchestrationApiService(t *testing.T) {

    configuration := client.NewConfiguration()
    apiClient := client.NewAPIClient(configuration)

    t.Run("Test OrchestrationApiService CreateWorkflow", func(t *testing.T) {

        t.Skip("skip test")  // remove to run test

        resp, httpRes, err := apiClient.OrchestrationApi.CreateWorkflow(context.Background()).Execute()

        require.Nil(t, err)
        require.NotNil(t, resp)
        assert.Equal(t, 200, httpRes.StatusCode)

    })

    t.Run("Test OrchestrationApiService GetFlow", func(t *testing.T) {

        t.Skip("skip test")  // remove to run test

        var flowId string

        resp, httpRes, err := apiClient.OrchestrationApi.GetFlow(context.Background(), flowId).Execute()

        require.Nil(t, err)
        require.NotNil(t, resp)
        assert.Equal(t, 200, httpRes.StatusCode)

    })

    t.Run("Test OrchestrationApiService GetWorkflowOccurrence", func(t *testing.T) {

        t.Skip("skip test")  // remove to run test

        var flowId string
        var runId string

        resp, httpRes, err := apiClient.OrchestrationApi.GetWorkflowOccurrence(context.Background(), flowId, runId).Execute()

        require.Nil(t, err)
        require.NotNil(t, resp)
        assert.Equal(t, 200, httpRes.StatusCode)

    })

    t.Run("Test OrchestrationApiService ListFlows", func(t *testing.T) {

        t.Skip("skip test")  // remove to run test

        resp, httpRes, err := apiClient.OrchestrationApi.ListFlows(context.Background()).Execute()

        require.Nil(t, err)
        require.NotNil(t, resp)
        assert.Equal(t, 200, httpRes.StatusCode)

    })

    t.Run("Test OrchestrationApiService ListRuns", func(t *testing.T) {

        t.Skip("skip test")  // remove to run test

        var flowId string

        resp, httpRes, err := apiClient.OrchestrationApi.ListRuns(context.Background(), flowId).Execute()

        require.Nil(t, err)
        require.NotNil(t, resp)
        assert.Equal(t, 200, httpRes.StatusCode)

    })

    t.Run("Test OrchestrationApiService OrchestrationgetServerInfo", func(t *testing.T) {

        t.Skip("skip test")  // remove to run test

        resp, httpRes, err := apiClient.OrchestrationApi.OrchestrationgetServerInfo(context.Background()).Execute()

        require.Nil(t, err)
        require.NotNil(t, resp)
        assert.Equal(t, 200, httpRes.StatusCode)

    })

    t.Run("Test OrchestrationApiService RunWorkflow", func(t *testing.T) {

        t.Skip("skip test")  // remove to run test

        var flowId string

        resp, httpRes, err := apiClient.OrchestrationApi.RunWorkflow(context.Background(), flowId).Execute()

        require.Nil(t, err)
        require.NotNil(t, resp)
        assert.Equal(t, 200, httpRes.StatusCode)

    })

}
