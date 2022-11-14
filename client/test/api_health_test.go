/*
Webhooks server

Testing HealthApiService

*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech);

package client

import (
    "context"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "testing"
    openapiclient "github.com/formancehq/webhooks/client"
)

func Test_client_HealthApiService(t *testing.T) {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)

    t.Run("Test HealthApiService HealthCheck", func(t *testing.T) {

        t.Skip("skip test")  // remove to run test

        resp, httpRes, err := apiClient.HealthApi.HealthCheck(context.Background()).Execute()

        require.Nil(t, err)
        require.NotNil(t, resp)
        assert.Equal(t, 200, httpRes.StatusCode)

    })

}
