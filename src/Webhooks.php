<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack;

use formance\stack\Models\Operations;
use formance\stack\Models\Shared;
use JMS\Serializer\DeserializationContext;

class Webhooks
{
    private SDKConfiguration $sdkConfiguration;

    /**
     * @param  SDKConfiguration  $sdkConfig
     */
    public function __construct(SDKConfiguration $sdkConfig)
    {
        $this->sdkConfiguration = $sdkConfig;
    }

    /**
     * Activate one config
     *
     * Activate a webhooks config by ID, to start receiving webhooks to its endpoint.
     *
     * @param  Operations\ActivateConfigRequest  $request
     * @return Operations\ActivateConfigResponse
     * @throws \formance\stack\Models\Errors\SDKException
     */
    public function activateConfig(
        ?Operations\ActivateConfigRequest $request,
    ): Operations\ActivateConfigResponse {
        $baseUrl = $this->sdkConfiguration->getServerUrl();
        $url = Utils\Utils::generateUrl($baseUrl, '/api/webhooks/configs/{id}/activate', Operations\ActivateConfigRequest::class, $request);
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = $this->sdkConfiguration->userAgent;
        $httpRequest = new \GuzzleHttp\Psr7\Request('PUT', $url);


        $httpResponse = $this->sdkConfiguration->securityClient->send($httpRequest, $options);
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $statusCode = $httpResponse->getStatusCode();
        if ($statusCode == 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $obj = $serializer->deserialize((string) $httpResponse->getBody(), '\formance\stack\Models\Shared\ConfigResponse', 'json', DeserializationContext::create()->setRequireAllRequiredProperties(true));
                $response = new Operations\ActivateConfigResponse(
                    statusCode: $statusCode,
                    contentType: $contentType,
                    rawResponse: $httpResponse,
                    configResponse: $obj);

                return $response;
            } else {
                throw new \formance\stack\Models\Errors\SDKException('Unknown content type received', $statusCode, $httpResponse->getBody()->getContents(), $httpResponse);
            }
        } else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $obj = $serializer->deserialize((string) $httpResponse->getBody(), '\formance\stack\Models\Errors\WebhooksErrorResponse', 'json', DeserializationContext::create()->setRequireAllRequiredProperties(true));
                throw $obj;
            } else {
                throw new \formance\stack\Models\Errors\SDKException('Unknown content type received', $statusCode, $httpResponse->getBody()->getContents(), $httpResponse);
            }
        }
    }

    /**
     * Change the signing secret of a config
     *
     * Change the signing secret of the endpoint of a webhooks config.
     *
     * If not passed or empty, a secret is automatically generated.
     * The format is a random string of bytes of size 24, base64 encoded. (larger size after encoding)
     *
     *
     * @param  Operations\ChangeConfigSecretRequest  $request
     * @return Operations\ChangeConfigSecretResponse
     * @throws \formance\stack\Models\Errors\SDKException
     */
    public function changeConfigSecret(
        ?Operations\ChangeConfigSecretRequest $request,
    ): Operations\ChangeConfigSecretResponse {
        $baseUrl = $this->sdkConfiguration->getServerUrl();
        $url = Utils\Utils::generateUrl($baseUrl, '/api/webhooks/configs/{id}/secret/change', Operations\ChangeConfigSecretRequest::class, $request);
        $options = ['http_errors' => false];
        $body = Utils\Utils::serializeRequestBody($request, 'configChangeSecret', 'json');
        if ($body !== null) {
            $options = array_merge_recursive($options, $body);
        }
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = $this->sdkConfiguration->userAgent;
        $httpRequest = new \GuzzleHttp\Psr7\Request('PUT', $url);


        $httpResponse = $this->sdkConfiguration->securityClient->send($httpRequest, $options);
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $statusCode = $httpResponse->getStatusCode();
        if ($statusCode == 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $obj = $serializer->deserialize((string) $httpResponse->getBody(), '\formance\stack\Models\Shared\ConfigResponse', 'json', DeserializationContext::create()->setRequireAllRequiredProperties(true));
                $response = new Operations\ChangeConfigSecretResponse(
                    statusCode: $statusCode,
                    contentType: $contentType,
                    rawResponse: $httpResponse,
                    configResponse: $obj);

                return $response;
            } else {
                throw new \formance\stack\Models\Errors\SDKException('Unknown content type received', $statusCode, $httpResponse->getBody()->getContents(), $httpResponse);
            }
        } else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $obj = $serializer->deserialize((string) $httpResponse->getBody(), '\formance\stack\Models\Errors\WebhooksErrorResponse', 'json', DeserializationContext::create()->setRequireAllRequiredProperties(true));
                throw $obj;
            } else {
                throw new \formance\stack\Models\Errors\SDKException('Unknown content type received', $statusCode, $httpResponse->getBody()->getContents(), $httpResponse);
            }
        }
    }

    /**
     * Deactivate one config
     *
     * Deactivate a webhooks config by ID, to stop receiving webhooks to its endpoint.
     *
     * @param  Operations\DeactivateConfigRequest  $request
     * @return Operations\DeactivateConfigResponse
     * @throws \formance\stack\Models\Errors\SDKException
     */
    public function deactivateConfig(
        ?Operations\DeactivateConfigRequest $request,
    ): Operations\DeactivateConfigResponse {
        $baseUrl = $this->sdkConfiguration->getServerUrl();
        $url = Utils\Utils::generateUrl($baseUrl, '/api/webhooks/configs/{id}/deactivate', Operations\DeactivateConfigRequest::class, $request);
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = $this->sdkConfiguration->userAgent;
        $httpRequest = new \GuzzleHttp\Psr7\Request('PUT', $url);


        $httpResponse = $this->sdkConfiguration->securityClient->send($httpRequest, $options);
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $statusCode = $httpResponse->getStatusCode();
        if ($statusCode == 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $obj = $serializer->deserialize((string) $httpResponse->getBody(), '\formance\stack\Models\Shared\ConfigResponse', 'json', DeserializationContext::create()->setRequireAllRequiredProperties(true));
                $response = new Operations\DeactivateConfigResponse(
                    statusCode: $statusCode,
                    contentType: $contentType,
                    rawResponse: $httpResponse,
                    configResponse: $obj);

                return $response;
            } else {
                throw new \formance\stack\Models\Errors\SDKException('Unknown content type received', $statusCode, $httpResponse->getBody()->getContents(), $httpResponse);
            }
        } else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $obj = $serializer->deserialize((string) $httpResponse->getBody(), '\formance\stack\Models\Errors\WebhooksErrorResponse', 'json', DeserializationContext::create()->setRequireAllRequiredProperties(true));
                throw $obj;
            } else {
                throw new \formance\stack\Models\Errors\SDKException('Unknown content type received', $statusCode, $httpResponse->getBody()->getContents(), $httpResponse);
            }
        }
    }

    /**
     * Delete one config
     *
     * Delete a webhooks config by ID.
     *
     * @param  Operations\DeleteConfigRequest  $request
     * @return Operations\DeleteConfigResponse
     * @throws \formance\stack\Models\Errors\SDKException
     */
    public function deleteConfig(
        ?Operations\DeleteConfigRequest $request,
    ): Operations\DeleteConfigResponse {
        $baseUrl = $this->sdkConfiguration->getServerUrl();
        $url = Utils\Utils::generateUrl($baseUrl, '/api/webhooks/configs/{id}', Operations\DeleteConfigRequest::class, $request);
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = $this->sdkConfiguration->userAgent;
        $httpRequest = new \GuzzleHttp\Psr7\Request('DELETE', $url);


        $httpResponse = $this->sdkConfiguration->securityClient->send($httpRequest, $options);
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $statusCode = $httpResponse->getStatusCode();
        if ($statusCode == 200) {
            return new Operations\DeleteConfigResponse(
                statusCode: $statusCode,
                contentType: $contentType,
                rawResponse: $httpResponse
            );
        } else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $obj = $serializer->deserialize((string) $httpResponse->getBody(), '\formance\stack\Models\Errors\WebhooksErrorResponse', 'json', DeserializationContext::create()->setRequireAllRequiredProperties(true));
                throw $obj;
            } else {
                throw new \formance\stack\Models\Errors\SDKException('Unknown content type received', $statusCode, $httpResponse->getBody()->getContents(), $httpResponse);
            }
        }
    }

    /**
     * Get many configs
     *
     * Sorted by updated date descending
     *
     * @param  Operations\GetManyConfigsRequest  $request
     * @return Operations\GetManyConfigsResponse
     * @throws \formance\stack\Models\Errors\SDKException
     */
    public function getManyConfigs(
        ?Operations\GetManyConfigsRequest $request,
    ): Operations\GetManyConfigsResponse {
        $baseUrl = $this->sdkConfiguration->getServerUrl();
        $url = Utils\Utils::generateUrl($baseUrl, '/api/webhooks/configs');
        $options = ['http_errors' => false];
        $options = array_merge_recursive($options, Utils\Utils::getQueryParams(Operations\GetManyConfigsRequest::class, $request));
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = $this->sdkConfiguration->userAgent;
        $httpRequest = new \GuzzleHttp\Psr7\Request('GET', $url);


        $httpResponse = $this->sdkConfiguration->securityClient->send($httpRequest, $options);
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $statusCode = $httpResponse->getStatusCode();
        if ($statusCode == 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $obj = $serializer->deserialize((string) $httpResponse->getBody(), '\formance\stack\Models\Shared\ConfigsResponse', 'json', DeserializationContext::create()->setRequireAllRequiredProperties(true));
                $response = new Operations\GetManyConfigsResponse(
                    statusCode: $statusCode,
                    contentType: $contentType,
                    rawResponse: $httpResponse,
                    configsResponse: $obj);

                return $response;
            } else {
                throw new \formance\stack\Models\Errors\SDKException('Unknown content type received', $statusCode, $httpResponse->getBody()->getContents(), $httpResponse);
            }
        } else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $obj = $serializer->deserialize((string) $httpResponse->getBody(), '\formance\stack\Models\Errors\WebhooksErrorResponse', 'json', DeserializationContext::create()->setRequireAllRequiredProperties(true));
                throw $obj;
            } else {
                throw new \formance\stack\Models\Errors\SDKException('Unknown content type received', $statusCode, $httpResponse->getBody()->getContents(), $httpResponse);
            }
        }
    }

    /**
     * Insert a new config
     *
     * Insert a new webhooks config.
     *
     * The endpoint should be a valid https URL and be unique.
     *
     * The secret is the endpoint's verification secret.
     * If not passed or empty, a secret is automatically generated.
     * The format is a random string of bytes of size 24, base64 encoded. (larger size after encoding)
     *
     * All eventTypes are converted to lower-case when inserted.
     *
     *
     * @param  Shared\ConfigUser  $request
     * @return Operations\InsertConfigResponse
     * @throws \formance\stack\Models\Errors\SDKException
     */
    public function insertConfig(
        Shared\ConfigUser $request,
    ): Operations\InsertConfigResponse {
        $baseUrl = $this->sdkConfiguration->getServerUrl();
        $url = Utils\Utils::generateUrl($baseUrl, '/api/webhooks/configs');
        $options = ['http_errors' => false];
        $body = Utils\Utils::serializeRequestBody($request, 'request', 'json');
        if ($body === null) {
            throw new \Exception('Request body is required');
        }
        $options = array_merge_recursive($options, $body);
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = $this->sdkConfiguration->userAgent;
        $httpRequest = new \GuzzleHttp\Psr7\Request('POST', $url);


        $httpResponse = $this->sdkConfiguration->securityClient->send($httpRequest, $options);
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $statusCode = $httpResponse->getStatusCode();
        if ($statusCode == 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $obj = $serializer->deserialize((string) $httpResponse->getBody(), '\formance\stack\Models\Shared\ConfigResponse', 'json', DeserializationContext::create()->setRequireAllRequiredProperties(true));
                $response = new Operations\InsertConfigResponse(
                    statusCode: $statusCode,
                    contentType: $contentType,
                    rawResponse: $httpResponse,
                    configResponse: $obj);

                return $response;
            } else {
                throw new \formance\stack\Models\Errors\SDKException('Unknown content type received', $statusCode, $httpResponse->getBody()->getContents(), $httpResponse);
            }
        } else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $obj = $serializer->deserialize((string) $httpResponse->getBody(), '\formance\stack\Models\Errors\WebhooksErrorResponse', 'json', DeserializationContext::create()->setRequireAllRequiredProperties(true));
                throw $obj;
            } else {
                throw new \formance\stack\Models\Errors\SDKException('Unknown content type received', $statusCode, $httpResponse->getBody()->getContents(), $httpResponse);
            }
        }
    }

    /**
     * Test one config
     *
     * Test a config by sending a webhook to its endpoint.
     *
     * @param  Operations\TestConfigRequest  $request
     * @return Operations\TestConfigResponse
     * @throws \formance\stack\Models\Errors\SDKException
     */
    public function testConfig(
        ?Operations\TestConfigRequest $request,
    ): Operations\TestConfigResponse {
        $baseUrl = $this->sdkConfiguration->getServerUrl();
        $url = Utils\Utils::generateUrl($baseUrl, '/api/webhooks/configs/{id}/test', Operations\TestConfigRequest::class, $request);
        $options = ['http_errors' => false];
        $options['headers']['Accept'] = 'application/json';
        $options['headers']['user-agent'] = $this->sdkConfiguration->userAgent;
        $httpRequest = new \GuzzleHttp\Psr7\Request('GET', $url);


        $httpResponse = $this->sdkConfiguration->securityClient->send($httpRequest, $options);
        $contentType = $httpResponse->getHeader('Content-Type')[0] ?? '';

        $statusCode = $httpResponse->getStatusCode();
        if ($statusCode == 200) {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $obj = $serializer->deserialize((string) $httpResponse->getBody(), '\formance\stack\Models\Shared\AttemptResponse', 'json', DeserializationContext::create()->setRequireAllRequiredProperties(true));
                $response = new Operations\TestConfigResponse(
                    statusCode: $statusCode,
                    contentType: $contentType,
                    rawResponse: $httpResponse,
                    attemptResponse: $obj);

                return $response;
            } else {
                throw new \formance\stack\Models\Errors\SDKException('Unknown content type received', $statusCode, $httpResponse->getBody()->getContents(), $httpResponse);
            }
        } else {
            if (Utils\Utils::matchContentType($contentType, 'application/json')) {
                $serializer = Utils\JSON::createSerializer();
                $obj = $serializer->deserialize((string) $httpResponse->getBody(), '\formance\stack\Models\Errors\WebhooksErrorResponse', 'json', DeserializationContext::create()->setRequireAllRequiredProperties(true));
                throw $obj;
            } else {
                throw new \formance\stack\Models\Errors\SDKException('Unknown content type received', $statusCode, $httpResponse->getBody()->getContents(), $httpResponse);
            }
        }
    }
}