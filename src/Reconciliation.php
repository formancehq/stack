<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack;

class Reconciliation 
{

	private SDKConfiguration $sdkConfiguration;

	/**
	 * @param SDKConfiguration $sdkConfig
	 */
	public function __construct(SDKConfiguration $sdkConfig)
	{
		$this->sdkConfiguration = $sdkConfig;
	}
	
    /**
     * Create a policy
     * 
     * Create a policy
     * 
     * @param \formance\stack\Models\Shared\PolicyRequest $request
     * @return \formance\stack\Models\Operations\CreatePolicyResponse
     */
	public function createPolicy(
        \formance\stack\Models\Shared\PolicyRequest $request,
    ): \formance\stack\Models\Operations\CreatePolicyResponse
    {
        $baseUrl = $this->sdkConfiguration->getServerUrl();
        $url = Utils\Utils::generateUrl($baseUrl, '/api/reconciliation/policies');
        
        $options = ['http_errors' => false];
        $body = Utils\Utils::serializeRequestBody($request, "request", "json");
        if ($body === null) {
            throw new \Exception('Request body is required');
        }
        $options = array_merge_recursive($options, $body);
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = $this->sdkConfiguration->userAgent;
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('POST', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $statusCode = $httpResponse->getStatusCode();

        $response = new \formance\stack\Models\Operations\CreatePolicyResponse();
        $response->statusCode = $statusCode;
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 201) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->policyResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\PolicyResponse', 'json');
            }
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->reconciliationErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ReconciliationErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Delete a policy
     * 
     * Delete a policy by its id.
     * 
     * @param \formance\stack\Models\Operations\DeletePolicyRequest $request
     * @return \formance\stack\Models\Operations\DeletePolicyResponse
     */
	public function deletePolicy(
        ?\formance\stack\Models\Operations\DeletePolicyRequest $request,
    ): \formance\stack\Models\Operations\DeletePolicyResponse
    {
        $baseUrl = $this->sdkConfiguration->getServerUrl();
        $url = Utils\Utils::generateUrl($baseUrl, '/api/reconciliation/policies/{policyID}', \formance\stack\Models\Operations\DeletePolicyRequest::class, $request);
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = $this->sdkConfiguration->userAgent;
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('DELETE', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $statusCode = $httpResponse->getStatusCode();

        $response = new \formance\stack\Models\Operations\DeletePolicyResponse();
        $response->statusCode = $statusCode;
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 204) {
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->reconciliationErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ReconciliationErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Get a policy
     * 
     * @param \formance\stack\Models\Operations\GetPolicyRequest $request
     * @return \formance\stack\Models\Operations\GetPolicyResponse
     */
	public function getPolicy(
        ?\formance\stack\Models\Operations\GetPolicyRequest $request,
    ): \formance\stack\Models\Operations\GetPolicyResponse
    {
        $baseUrl = $this->sdkConfiguration->getServerUrl();
        $url = Utils\Utils::generateUrl($baseUrl, '/api/reconciliation/policies/{policyID}', \formance\stack\Models\Operations\GetPolicyRequest::class, $request);
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = $this->sdkConfiguration->userAgent;
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $statusCode = $httpResponse->getStatusCode();

        $response = new \formance\stack\Models\Operations\GetPolicyResponse();
        $response->statusCode = $statusCode;
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->policyResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\PolicyResponse', 'json');
            }
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->reconciliationErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ReconciliationErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Get a reconciliation
     * 
     * @param \formance\stack\Models\Operations\GetReconciliationRequest $request
     * @return \formance\stack\Models\Operations\GetReconciliationResponse
     */
	public function getReconciliation(
        ?\formance\stack\Models\Operations\GetReconciliationRequest $request,
    ): \formance\stack\Models\Operations\GetReconciliationResponse
    {
        $baseUrl = $this->sdkConfiguration->getServerUrl();
        $url = Utils\Utils::generateUrl($baseUrl, '/api/reconciliation/reconciliations/{reconciliationID}', \formance\stack\Models\Operations\GetReconciliationRequest::class, $request);
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = $this->sdkConfiguration->userAgent;
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $statusCode = $httpResponse->getStatusCode();

        $response = new \formance\stack\Models\Operations\GetReconciliationResponse();
        $response->statusCode = $statusCode;
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->reconciliationResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ReconciliationResponse', 'json');
            }
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->reconciliationErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ReconciliationErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * List policies
     * 
     * @param \formance\stack\Models\Operations\ListPoliciesRequest $request
     * @return \formance\stack\Models\Operations\ListPoliciesResponse
     */
	public function listPolicies(
        ?\formance\stack\Models\Operations\ListPoliciesRequest $request,
    ): \formance\stack\Models\Operations\ListPoliciesResponse
    {
        $baseUrl = $this->sdkConfiguration->getServerUrl();
        $url = Utils\Utils::generateUrl($baseUrl, '/api/reconciliation/policies');
        
        $options = ['http_errors' => false];
        $options = array_merge_recursive($options, Utils\Utils::getQueryParams(\formance\stack\Models\Operations\ListPoliciesRequest::class, $request, null));
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = $this->sdkConfiguration->userAgent;
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $statusCode = $httpResponse->getStatusCode();

        $response = new \formance\stack\Models\Operations\ListPoliciesResponse();
        $response->statusCode = $statusCode;
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->policiesCursorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\PoliciesCursorResponse', 'json');
            }
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->reconciliationErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ReconciliationErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * List reconciliations
     * 
     * @param \formance\stack\Models\Operations\ListReconciliationsRequest $request
     * @return \formance\stack\Models\Operations\ListReconciliationsResponse
     */
	public function listReconciliations(
        ?\formance\stack\Models\Operations\ListReconciliationsRequest $request,
    ): \formance\stack\Models\Operations\ListReconciliationsResponse
    {
        $baseUrl = $this->sdkConfiguration->getServerUrl();
        $url = Utils\Utils::generateUrl($baseUrl, '/api/reconciliation/reconciliations');
        
        $options = ['http_errors' => false];
        $options = array_merge_recursive($options, Utils\Utils::getQueryParams(\formance\stack\Models\Operations\ListReconciliationsRequest::class, $request, null));
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = $this->sdkConfiguration->userAgent;
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $statusCode = $httpResponse->getStatusCode();

        $response = new \formance\stack\Models\Operations\ListReconciliationsResponse();
        $response->statusCode = $statusCode;
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->reconciliationsCursorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ReconciliationsCursorResponse', 'json');
            }
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->reconciliationErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ReconciliationErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Reconcile using a policy
     * 
     * Reconcile using a policy
     * 
     * @param \formance\stack\Models\Operations\ReconcileRequest $request
     * @return \formance\stack\Models\Operations\ReconcileResponse
     */
	public function reconcile(
        \formance\stack\Models\Operations\ReconcileRequest $request,
    ): \formance\stack\Models\Operations\ReconcileResponse
    {
        $baseUrl = $this->sdkConfiguration->getServerUrl();
        $url = Utils\Utils::generateUrl($baseUrl, '/api/reconciliation/policies/{policyID}/reconciliation', \formance\stack\Models\Operations\ReconcileRequest::class, $request);
        
        $options = ['http_errors' => false];
        $body = Utils\Utils::serializeRequestBody($request, "reconciliationRequest", "json");
        if ($body === null) {
            throw new \Exception('Request body is required');
        }
        $options = array_merge_recursive($options, $body);
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = $this->sdkConfiguration->userAgent;
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('POST', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $statusCode = $httpResponse->getStatusCode();

        $response = new \formance\stack\Models\Operations\ReconcileResponse();
        $response->statusCode = $statusCode;
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->reconciliationResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ReconciliationResponse', 'json');
            }
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->reconciliationErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ReconciliationErrorResponse', 'json');
            }
        }

        return $response;
    }
	
    /**
     * Get server info
     * 
     * @return \formance\stack\Models\Operations\ReconciliationgetServerInfoResponse
     */
	public function reconciliationgetServerInfo(
    ): \formance\stack\Models\Operations\ReconciliationgetServerInfoResponse
    {
        $baseUrl = $this->sdkConfiguration->getServerUrl();
        $url = Utils\Utils::generateUrl($baseUrl, '/api/reconciliation/_info');
        
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = $this->sdkConfiguration->userAgent;
        
        $httpResponse = $this->sdkConfiguration->securityClient->request('GET', $url, $options);
        
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $statusCode = $httpResponse->getStatusCode();

        $response = new \formance\stack\Models\Operations\ReconciliationgetServerInfoResponse();
        $response->statusCode = $statusCode;
        $response->contentType = $contentType;
        $response->rawResponse = $httpResponse;
        
        if ($httpResponse->getStatusCode() === 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->serverInfo = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ServerInfo', 'json');
            }
        }
        else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $response->reconciliationErrorResponse = $serializer->deserialize((string)$httpResponse->getBody(), 'formance\stack\Models\Shared\ReconciliationErrorResponse', 'json');
            }
        }

        return $response;
    }
}