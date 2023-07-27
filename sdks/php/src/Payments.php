<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack;

class Payments 
{

	// SDK private variables namespaced with _ to avoid conflicts with API models
	private \GuzzleHttp\ClientInterface $_defaultClient;
	private \GuzzleHttp\ClientInterface $_securityClient;
	private string $_serverUrl;
	private string $_language;
	private string $_sdkVersion;
	private string $_genVersion;	

	/**
	 * @param \GuzzleHttp\ClientInterface $defaultClient
	 * @param \GuzzleHttp\ClientInterface $securityClient
	 * @param string $serverUrl
	 * @param string $language
	 * @param string $sdkVersion
	 * @param string $genVersion
	 */
	public function __construct(\GuzzleHttp\ClientInterface $defaultClient, \GuzzleHttp\ClientInterface $securityClient, string $serverUrl, string $language, string $sdkVersion, string $genVersion)
	{
		$this->_defaultClient = $defaultClient;
		$this->_securityClient = $securityClient;
		$this->_serverUrl = $serverUrl;
		$this->_language = $language;
		$this->_sdkVersion = $sdkVersion;
		$this->_genVersion = $genVersion;
	}
	
    /**
     * Transfer funds between Stripe accounts
     * 
     * Execute a transfer between two Stripe accounts.
     * 
     * @param \formance\stack\Models\Shared\StripeTransferRequest $request
     * @return \formance\stack\Models\Operations\ConnectorsStripeTransferResponse
     */
	public function connectorsStripeTransfer(
        \formance\stack\Models\Shared\StripeTransferRequest $request,
    ): \formance\stack\Models\Operations\ConnectorsStripeTransferResponse
    {
        $baseUrl = $this->_serverUrl;
        $url = Utils\Utils::generateUrl($baseUrl, '/api/payments/connectors/stripe/transfers');
        
        $options = ['http_errors' => false];
        $body = Utils\Utils::serializeRequestBody($request, "request", "json");
        if ($body === null) {
            throw new \Exception('Request body is required');
        }
        $options = array_merge_recursive($options, $body);
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s', $this->_language, $this->_sdkVersion, $this->_genVersion);
        
        $httpResponse = $this->_securityClient->request('POST', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\ConnectorsStripeTransferResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->stripeTransferResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'array<string, mixed>', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Transfer funds between Connector accounts
     * 
     * Execute a transfer between two accounts.
     * 
     * @param \formance\stack\Models\Operations\ConnectorsTransferRequest $request
     * @return \formance\stack\Models\Operations\ConnectorsTransferResponse
     */
	public function connectorsTransfer(
        \formance\stack\Models\Operations\ConnectorsTransferRequest $request,
    ): \formance\stack\Models\Operations\ConnectorsTransferResponse
    {
        $baseUrl = $this->_serverUrl;
        $url = Utils\Utils::generateUrl($baseUrl, '/api/payments/connectors/{connector}/transfers', \formance\stack\Models\Operations\ConnectorsTransferRequest::class, $request);
        
        $options = ['http_errors' => false];
        $body = Utils\Utils::serializeRequestBody($request, "transferRequest", "json");
        if ($body === null) {
            throw new \Exception('Request body is required');
        }
        $options = array_merge_recursive($options, $body);
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s', $this->_language, $this->_sdkVersion, $this->_genVersion);
        
        $httpResponse = $this->_securityClient->request('POST', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\ConnectorsTransferResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->transferResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\TransferResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Get account balances
     * 
     * @param \formance\stack\Models\Operations\GetAccountBalancesRequest $request
     * @return \formance\stack\Models\Operations\GetAccountBalancesResponse
     */
	public function getAccountBalances(
        \formance\stack\Models\Operations\GetAccountBalancesRequest $request,
    ): \formance\stack\Models\Operations\GetAccountBalancesResponse
    {
        $baseUrl = $this->_serverUrl;
        $url = Utils\Utils::generateUrl($baseUrl, '/api/payments/accounts/{accountID}/balances', \formance\stack\Models\Operations\GetAccountBalancesRequest::class, $request);
        
        $options = ['http_errors' => false];
        $options = array_merge_recursive($options, Utils\Utils::getQueryParams(\formance\stack\Models\Operations\GetAccountBalancesRequest::class, $request, null));
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s', $this->_language, $this->_sdkVersion, $this->_genVersion);
        
        $httpResponse = $this->_securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\GetAccountBalancesResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->balancesCursor = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\BalancesCursor', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Read a specific task of the connector
     * 
     * Get a specific task associated to the connector.
     * 
     * @param \formance\stack\Models\Operations\GetConnectorTaskRequest $request
     * @return \formance\stack\Models\Operations\GetConnectorTaskResponse
     */
	public function getConnectorTask(
        \formance\stack\Models\Operations\GetConnectorTaskRequest $request,
    ): \formance\stack\Models\Operations\GetConnectorTaskResponse
    {
        $baseUrl = $this->_serverUrl;
        $url = Utils\Utils::generateUrl($baseUrl, '/api/payments/connectors/{connector}/tasks/{taskId}', \formance\stack\Models\Operations\GetConnectorTaskRequest::class, $request);
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s', $this->_language, $this->_sdkVersion, $this->_genVersion);
        
        $httpResponse = $this->_securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\GetConnectorTaskResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->taskResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\TaskResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Get a payment
     * 
     * @param \formance\stack\Models\Operations\GetPaymentRequest $request
     * @return \formance\stack\Models\Operations\GetPaymentResponse
     */
	public function getPayment(
        \formance\stack\Models\Operations\GetPaymentRequest $request,
    ): \formance\stack\Models\Operations\GetPaymentResponse
    {
        $baseUrl = $this->_serverUrl;
        $url = Utils\Utils::generateUrl($baseUrl, '/api/payments/payments/{paymentId}', \formance\stack\Models\Operations\GetPaymentRequest::class, $request);
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s', $this->_language, $this->_sdkVersion, $this->_genVersion);
        
        $httpResponse = $this->_securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\GetPaymentResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->paymentResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\PaymentResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Install a connector
     * 
     * Install a connector by its name and config.
     * 
     * @param \formance\stack\Models\Operations\InstallConnectorRequest $request
     * @return \formance\stack\Models\Operations\InstallConnectorResponse
     */
	public function installConnector(
        \formance\stack\Models\Operations\InstallConnectorRequest $request,
    ): \formance\stack\Models\Operations\InstallConnectorResponse
    {
        $baseUrl = $this->_serverUrl;
        $url = Utils\Utils::generateUrl($baseUrl, '/api/payments/connectors/{connector}', \formance\stack\Models\Operations\InstallConnectorRequest::class, $request);
        
        $options = ['http_errors' => false];
        $body = Utils\Utils::serializeRequestBody($request, "requestBody", "json");
        if ($body === null) {
            throw new \Exception('Request body is required');
        }
        $options = array_merge_recursive($options, $body);
        $options['headers']['Accept'] = '*/*';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s', $this->_language, $this->_sdkVersion, $this->_genVersion);
        
        $httpResponse = $this->_securityClient->request('POST', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\InstallConnectorResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 204) {
        }

        return $response;
    }
	
    /**
     * List all installed connectors
     * 
     * List all installed connectors.
     * 
     * @return \formance\stack\Models\Operations\ListAllConnectorsResponse
     */
	public function listAllConnectors(
    ): \formance\stack\Models\Operations\ListAllConnectorsResponse
    {
        $baseUrl = $this->_serverUrl;
        $url = Utils\Utils::generateUrl($baseUrl, '/api/payments/connectors');
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s', $this->_language, $this->_sdkVersion, $this->_genVersion);
        
        $httpResponse = $this->_securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\ListAllConnectorsResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->connectorsResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ConnectorsResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * List the configs of each available connector
     * 
     * List the configs of each available connector.
     * 
     * @return \formance\stack\Models\Operations\ListConfigsAvailableConnectorsResponse
     */
	public function listConfigsAvailableConnectors(
    ): \formance\stack\Models\Operations\ListConfigsAvailableConnectorsResponse
    {
        $baseUrl = $this->_serverUrl;
        $url = Utils\Utils::generateUrl($baseUrl, '/api/payments/connectors/configs');
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s', $this->_language, $this->_sdkVersion, $this->_genVersion);
        
        $httpResponse = $this->_securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\ListConfigsAvailableConnectorsResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->connectorsConfigsResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ConnectorsConfigsResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * List tasks from a connector
     * 
     * List all tasks associated with this connector.
     * 
     * @param \formance\stack\Models\Operations\ListConnectorTasksRequest $request
     * @return \formance\stack\Models\Operations\ListConnectorTasksResponse
     */
	public function listConnectorTasks(
        \formance\stack\Models\Operations\ListConnectorTasksRequest $request,
    ): \formance\stack\Models\Operations\ListConnectorTasksResponse
    {
        $baseUrl = $this->_serverUrl;
        $url = Utils\Utils::generateUrl($baseUrl, '/api/payments/connectors/{connector}/tasks', \formance\stack\Models\Operations\ListConnectorTasksRequest::class, $request);
        
        $options = ['http_errors' => false];
        $options = array_merge_recursive($options, Utils\Utils::getQueryParams(\formance\stack\Models\Operations\ListConnectorTasksRequest::class, $request, null));
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s', $this->_language, $this->_sdkVersion, $this->_genVersion);
        
        $httpResponse = $this->_securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\ListConnectorTasksResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->tasksCursor = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\TasksCursor', 'json');
            }
        }

        return $response;
    }
	
    /**
     * List transfers and their statuses
     * 
     * List transfers
     * 
     * @param \formance\stack\Models\Operations\ListConnectorsTransfersRequest $request
     * @return \formance\stack\Models\Operations\ListConnectorsTransfersResponse
     */
	public function listConnectorsTransfers(
        \formance\stack\Models\Operations\ListConnectorsTransfersRequest $request,
    ): \formance\stack\Models\Operations\ListConnectorsTransfersResponse
    {
        $baseUrl = $this->_serverUrl;
        $url = Utils\Utils::generateUrl($baseUrl, '/api/payments/connectors/{connector}/transfers', \formance\stack\Models\Operations\ListConnectorsTransfersRequest::class, $request);
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s', $this->_language, $this->_sdkVersion, $this->_genVersion);
        
        $httpResponse = $this->_securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\ListConnectorsTransfersResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->transfersResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\TransfersResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * List payments
     * 
     * @param \formance\stack\Models\Operations\ListPaymentsRequest $request
     * @return \formance\stack\Models\Operations\ListPaymentsResponse
     */
	public function listPayments(
        \formance\stack\Models\Operations\ListPaymentsRequest $request,
    ): \formance\stack\Models\Operations\ListPaymentsResponse
    {
        $baseUrl = $this->_serverUrl;
        $url = Utils\Utils::generateUrl($baseUrl, '/api/payments/payments');
        
        $options = ['http_errors' => false];
        $options = array_merge_recursive($options, Utils\Utils::getQueryParams(\formance\stack\Models\Operations\ListPaymentsRequest::class, $request, null));
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s', $this->_language, $this->_sdkVersion, $this->_genVersion);
        
        $httpResponse = $this->_securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\ListPaymentsResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->paymentsCursor = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\PaymentsCursor', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Get server info
     * 
     * @return \formance\stack\Models\Operations\PaymentsgetServerInfoResponse
     */
	public function paymentsgetServerInfo(
    ): \formance\stack\Models\Operations\PaymentsgetServerInfoResponse
    {
        $baseUrl = $this->_serverUrl;
        $url = Utils\Utils::generateUrl($baseUrl, '/api/payments/_info');
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s', $this->_language, $this->_sdkVersion, $this->_genVersion);
        
        $httpResponse = $this->_securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\PaymentsgetServerInfoResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->serverInfo = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ServerInfo', 'json');
            }
        }

        return $response;
    }
	
    /**
     * List accounts
     * 
     * @param \formance\stack\Models\Operations\PaymentslistAccountsRequest $request
     * @return \formance\stack\Models\Operations\PaymentslistAccountsResponse
     */
	public function paymentslistAccounts(
        \formance\stack\Models\Operations\PaymentslistAccountsRequest $request,
    ): \formance\stack\Models\Operations\PaymentslistAccountsResponse
    {
        $baseUrl = $this->_serverUrl;
        $url = Utils\Utils::generateUrl($baseUrl, '/api/payments/accounts');
        
        $options = ['http_errors' => false];
        $options = array_merge_recursive($options, Utils\Utils::getQueryParams(\formance\stack\Models\Operations\PaymentslistAccountsRequest::class, $request, null));
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s', $this->_language, $this->_sdkVersion, $this->_genVersion);
        
        $httpResponse = $this->_securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\PaymentslistAccountsResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->accountsCursor = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\AccountsCursor', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Read the config of a connector
     * 
     * Read connector config
     * 
     * @param \formance\stack\Models\Operations\ReadConnectorConfigRequest $request
     * @return \formance\stack\Models\Operations\ReadConnectorConfigResponse
     */
	public function readConnectorConfig(
        \formance\stack\Models\Operations\ReadConnectorConfigRequest $request,
    ): \formance\stack\Models\Operations\ReadConnectorConfigResponse
    {
        $baseUrl = $this->_serverUrl;
        $url = Utils\Utils::generateUrl($baseUrl, '/api/payments/connectors/{connector}/config', \formance\stack\Models\Operations\ReadConnectorConfigRequest::class, $request);
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s', $this->_language, $this->_sdkVersion, $this->_genVersion);
        
        $httpResponse = $this->_securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\ReadConnectorConfigResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->connectorConfigResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ConnectorConfigResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Reset a connector
     * 
     * Reset a connector by its name.
     * It will remove the connector and ALL PAYMENTS generated with it.
     * 
     * 
     * @param \formance\stack\Models\Operations\ResetConnectorRequest $request
     * @return \formance\stack\Models\Operations\ResetConnectorResponse
     */
	public function resetConnector(
        \formance\stack\Models\Operations\ResetConnectorRequest $request,
    ): \formance\stack\Models\Operations\ResetConnectorResponse
    {
        $baseUrl = $this->_serverUrl;
        $url = Utils\Utils::generateUrl($baseUrl, '/api/payments/connectors/{connector}/reset', \formance\stack\Models\Operations\ResetConnectorRequest::class, $request);
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = '*/*';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s', $this->_language, $this->_sdkVersion, $this->_genVersion);
        
        $httpResponse = $this->_securityClient->request('POST', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\ResetConnectorResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 204) {
        }

        return $response;
    }
	
    /**
     * Uninstall a connector
     * 
     * Uninstall a connector by its name.
     * 
     * @param \formance\stack\Models\Operations\UninstallConnectorRequest $request
     * @return \formance\stack\Models\Operations\UninstallConnectorResponse
     */
	public function uninstallConnector(
        \formance\stack\Models\Operations\UninstallConnectorRequest $request,
    ): \formance\stack\Models\Operations\UninstallConnectorResponse
    {
        $baseUrl = $this->_serverUrl;
        $url = Utils\Utils::generateUrl($baseUrl, '/api/payments/connectors/{connector}', \formance\stack\Models\Operations\UninstallConnectorRequest::class, $request);
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = '*/*';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s', $this->_language, $this->_sdkVersion, $this->_genVersion);
        
        $httpResponse = $this->_securityClient->request('DELETE', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\UninstallConnectorResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 204) {
        }

        return $response;
    }
	
    /**
     * Update metadata
     * 
     * @param \formance\stack\Models\Operations\UpdateMetadataRequest $request
     * @return \formance\stack\Models\Operations\UpdateMetadataResponse
     */
	public function updateMetadata(
        \formance\stack\Models\Operations\UpdateMetadataRequest $request,
    ): \formance\stack\Models\Operations\UpdateMetadataResponse
    {
        $baseUrl = $this->_serverUrl;
        $url = Utils\Utils::generateUrl($baseUrl, '/api/payments/payments/{paymentId}/metadata', \formance\stack\Models\Operations\UpdateMetadataRequest::class, $request);
        
        $options = ['http_errors' => false];
        $body = Utils\Utils::serializeRequestBody($request, "paymentMetadata", "json");
        if ($body === null) {
            throw new \Exception('Request body is required');
        }
        $options = array_merge_recursive($options, $body);
        $options['headers']['Accept'] = '*/*';
        $options['headers']['user-agent'] = sprintf('speakeasy-sdk/%s %s %s', $this->_language, $this->_sdkVersion, $this->_genVersion);
        
        $httpResponse = $this->_securityClient->request('PATCH', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $response = new \formance\stack\Models\Operations\UpdateMetadataResponse();
        $response->statusCode = $httpResponse->getStatusCode();
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 204) {
        }

        return $response;
    }
}