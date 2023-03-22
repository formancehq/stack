/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.formance.formance_sdk.utils.HTTPClient;
import com.formance.formance_sdk.utils.HTTPRequest;
import com.formance.formance_sdk.utils.JSON;
import com.formance.formance_sdk.utils.SerializedBody;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import org.apache.http.NameValuePair;

public class Webhooks {
	
	private HTTPClient _defaultClient;
	private HTTPClient _securityClient;
	private String _serverUrl;
	private String _language;
	private String _sdkVersion;
	private String _genVersion;

	public Webhooks(HTTPClient defaultClient, HTTPClient securityClient, String serverUrl, String language, String sdkVersion, String genVersion) {
		this._defaultClient = defaultClient;
		this._securityClient = securityClient;
		this._serverUrl = serverUrl;
		this._language = language;
		this._sdkVersion = sdkVersion;
		this._genVersion = genVersion;
	}

    /**
     * Activate one config
     * Activate a webhooks config by ID, to start receiving webhooks to its endpoint.
     * @param request the request object containing all of the parameters for the API call
     * @return the response from the API call
     * @throws Exception if the API call fails
     */
    public com.formance.formance_sdk.models.operations.ActivateConfigResponse activateConfig(com.formance.formance_sdk.models.operations.ActivateConfigRequest request) throws Exception {
        String baseUrl = this._serverUrl;
        String url = com.formance.formance_sdk.utils.Utils.generateURL(com.formance.formance_sdk.models.operations.ActivateConfigRequest.class, baseUrl, "/api/webhooks/configs/{id}/activate", request, null);
        
        HTTPRequest req = new HTTPRequest();
        req.setMethod("PUT");
        req.setURL(url);

        req.addHeader("Accept", "application/json;q=1, application/json;q=0");
        req.addHeader("user-agent", String.format("speakeasy-sdk/%s %s %s", this._language, this._sdkVersion, this._genVersion));
        
        HTTPClient client = this._securityClient;
        
        HttpResponse<byte[]> httpRes = client.send(req);

        String contentType = httpRes.headers().firstValue("Content-Type").orElse("application/octet-stream");

        com.formance.formance_sdk.models.operations.ActivateConfigResponse res = new com.formance.formance_sdk.models.operations.ActivateConfigResponse(contentType, httpRes.statusCode()) {{
            configResponse = null;
            errorResponse = null;
        }};
        res.rawResponse = httpRes;
        
        if (httpRes.statusCode() == 200) {
            if (com.formance.formance_sdk.utils.Utils.matchContentType(contentType, "application/json")) {
                ObjectMapper mapper = JSON.getMapper();
                com.formance.formance_sdk.models.shared.ConfigResponse out = mapper.readValue(new String(httpRes.body(), StandardCharsets.UTF_8), com.formance.formance_sdk.models.shared.ConfigResponse.class);
                res.configResponse = out;
            }
        }
        else {
            if (com.formance.formance_sdk.utils.Utils.matchContentType(contentType, "application/json")) {
                ObjectMapper mapper = JSON.getMapper();
                com.formance.formance_sdk.models.shared.ErrorResponse out = mapper.readValue(new String(httpRes.body(), StandardCharsets.UTF_8), com.formance.formance_sdk.models.shared.ErrorResponse.class);
                res.errorResponse = out;
            }
        }

        return res;
    }

    /**
     * Change the signing secret of a config
     * Change the signing secret of the endpoint of a webhooks config.
     * 
     * If not passed or empty, a secret is automatically generated.
     * The format is a random string of bytes of size 24, base64 encoded. (larger size after encoding)
     * 
     * @param request the request object containing all of the parameters for the API call
     * @return the response from the API call
     * @throws Exception if the API call fails
     */
    public com.formance.formance_sdk.models.operations.ChangeConfigSecretResponse changeConfigSecret(com.formance.formance_sdk.models.operations.ChangeConfigSecretRequest request) throws Exception {
        String baseUrl = this._serverUrl;
        String url = com.formance.formance_sdk.utils.Utils.generateURL(com.formance.formance_sdk.models.operations.ChangeConfigSecretRequest.class, baseUrl, "/api/webhooks/configs/{id}/secret/change", request, null);
        
        HTTPRequest req = new HTTPRequest();
        req.setMethod("PUT");
        req.setURL(url);
        SerializedBody serializedRequestBody = com.formance.formance_sdk.utils.Utils.serializeRequestBody(request, "configChangeSecret", "json");
        req.setBody(serializedRequestBody);

        req.addHeader("Accept", "application/json;q=1, application/json;q=0");
        req.addHeader("user-agent", String.format("speakeasy-sdk/%s %s %s", this._language, this._sdkVersion, this._genVersion));
        
        HTTPClient client = this._securityClient;
        
        HttpResponse<byte[]> httpRes = client.send(req);

        String contentType = httpRes.headers().firstValue("Content-Type").orElse("application/octet-stream");

        com.formance.formance_sdk.models.operations.ChangeConfigSecretResponse res = new com.formance.formance_sdk.models.operations.ChangeConfigSecretResponse(contentType, httpRes.statusCode()) {{
            configResponse = null;
            errorResponse = null;
        }};
        res.rawResponse = httpRes;
        
        if (httpRes.statusCode() == 200) {
            if (com.formance.formance_sdk.utils.Utils.matchContentType(contentType, "application/json")) {
                ObjectMapper mapper = JSON.getMapper();
                com.formance.formance_sdk.models.shared.ConfigResponse out = mapper.readValue(new String(httpRes.body(), StandardCharsets.UTF_8), com.formance.formance_sdk.models.shared.ConfigResponse.class);
                res.configResponse = out;
            }
        }
        else {
            if (com.formance.formance_sdk.utils.Utils.matchContentType(contentType, "application/json")) {
                ObjectMapper mapper = JSON.getMapper();
                com.formance.formance_sdk.models.shared.ErrorResponse out = mapper.readValue(new String(httpRes.body(), StandardCharsets.UTF_8), com.formance.formance_sdk.models.shared.ErrorResponse.class);
                res.errorResponse = out;
            }
        }

        return res;
    }

    /**
     * Deactivate one config
     * Deactivate a webhooks config by ID, to stop receiving webhooks to its endpoint.
     * @param request the request object containing all of the parameters for the API call
     * @return the response from the API call
     * @throws Exception if the API call fails
     */
    public com.formance.formance_sdk.models.operations.DeactivateConfigResponse deactivateConfig(com.formance.formance_sdk.models.operations.DeactivateConfigRequest request) throws Exception {
        String baseUrl = this._serverUrl;
        String url = com.formance.formance_sdk.utils.Utils.generateURL(com.formance.formance_sdk.models.operations.DeactivateConfigRequest.class, baseUrl, "/api/webhooks/configs/{id}/deactivate", request, null);
        
        HTTPRequest req = new HTTPRequest();
        req.setMethod("PUT");
        req.setURL(url);

        req.addHeader("Accept", "application/json;q=1, application/json;q=0");
        req.addHeader("user-agent", String.format("speakeasy-sdk/%s %s %s", this._language, this._sdkVersion, this._genVersion));
        
        HTTPClient client = this._securityClient;
        
        HttpResponse<byte[]> httpRes = client.send(req);

        String contentType = httpRes.headers().firstValue("Content-Type").orElse("application/octet-stream");

        com.formance.formance_sdk.models.operations.DeactivateConfigResponse res = new com.formance.formance_sdk.models.operations.DeactivateConfigResponse(contentType, httpRes.statusCode()) {{
            configResponse = null;
            errorResponse = null;
        }};
        res.rawResponse = httpRes;
        
        if (httpRes.statusCode() == 200) {
            if (com.formance.formance_sdk.utils.Utils.matchContentType(contentType, "application/json")) {
                ObjectMapper mapper = JSON.getMapper();
                com.formance.formance_sdk.models.shared.ConfigResponse out = mapper.readValue(new String(httpRes.body(), StandardCharsets.UTF_8), com.formance.formance_sdk.models.shared.ConfigResponse.class);
                res.configResponse = out;
            }
        }
        else {
            if (com.formance.formance_sdk.utils.Utils.matchContentType(contentType, "application/json")) {
                ObjectMapper mapper = JSON.getMapper();
                com.formance.formance_sdk.models.shared.ErrorResponse out = mapper.readValue(new String(httpRes.body(), StandardCharsets.UTF_8), com.formance.formance_sdk.models.shared.ErrorResponse.class);
                res.errorResponse = out;
            }
        }

        return res;
    }

    /**
     * Delete one config
     * Delete a webhooks config by ID.
     * @param request the request object containing all of the parameters for the API call
     * @return the response from the API call
     * @throws Exception if the API call fails
     */
    public com.formance.formance_sdk.models.operations.DeleteConfigResponse deleteConfig(com.formance.formance_sdk.models.operations.DeleteConfigRequest request) throws Exception {
        String baseUrl = this._serverUrl;
        String url = com.formance.formance_sdk.utils.Utils.generateURL(com.formance.formance_sdk.models.operations.DeleteConfigRequest.class, baseUrl, "/api/webhooks/configs/{id}", request, null);
        
        HTTPRequest req = new HTTPRequest();
        req.setMethod("DELETE");
        req.setURL(url);

        req.addHeader("Accept", "application/json");
        req.addHeader("user-agent", String.format("speakeasy-sdk/%s %s %s", this._language, this._sdkVersion, this._genVersion));
        
        HTTPClient client = this._securityClient;
        
        HttpResponse<byte[]> httpRes = client.send(req);

        String contentType = httpRes.headers().firstValue("Content-Type").orElse("application/octet-stream");

        com.formance.formance_sdk.models.operations.DeleteConfigResponse res = new com.formance.formance_sdk.models.operations.DeleteConfigResponse(contentType, httpRes.statusCode()) {{
            errorResponse = null;
        }};
        res.rawResponse = httpRes;
        
        if (httpRes.statusCode() == 200) {
        }
        else {
            if (com.formance.formance_sdk.utils.Utils.matchContentType(contentType, "application/json")) {
                ObjectMapper mapper = JSON.getMapper();
                com.formance.formance_sdk.models.shared.ErrorResponse out = mapper.readValue(new String(httpRes.body(), StandardCharsets.UTF_8), com.formance.formance_sdk.models.shared.ErrorResponse.class);
                res.errorResponse = out;
            }
        }

        return res;
    }

    /**
     * Get many configs
     * Sorted by updated date descending
     * @param request the request object containing all of the parameters for the API call
     * @return the response from the API call
     * @throws Exception if the API call fails
     */
    public com.formance.formance_sdk.models.operations.GetManyConfigsResponse getManyConfigs(com.formance.formance_sdk.models.operations.GetManyConfigsRequest request) throws Exception {
        String baseUrl = this._serverUrl;
        String url = com.formance.formance_sdk.utils.Utils.generateURL(baseUrl, "/api/webhooks/configs");
        
        HTTPRequest req = new HTTPRequest();
        req.setMethod("GET");
        req.setURL(url);

        req.addHeader("Accept", "application/json;q=1, application/json;q=0");
        req.addHeader("user-agent", String.format("speakeasy-sdk/%s %s %s", this._language, this._sdkVersion, this._genVersion));
        java.util.List<NameValuePair> queryParams = com.formance.formance_sdk.utils.Utils.getQueryParams(com.formance.formance_sdk.models.operations.GetManyConfigsRequest.class, request, null);
        if (queryParams != null) {
            for (NameValuePair queryParam : queryParams) {
                req.addQueryParam(queryParam);
            }
        }
        
        HTTPClient client = this._securityClient;
        
        HttpResponse<byte[]> httpRes = client.send(req);

        String contentType = httpRes.headers().firstValue("Content-Type").orElse("application/octet-stream");

        com.formance.formance_sdk.models.operations.GetManyConfigsResponse res = new com.formance.formance_sdk.models.operations.GetManyConfigsResponse(contentType, httpRes.statusCode()) {{
            configsResponse = null;
            errorResponse = null;
        }};
        res.rawResponse = httpRes;
        
        if (httpRes.statusCode() == 200) {
            if (com.formance.formance_sdk.utils.Utils.matchContentType(contentType, "application/json")) {
                ObjectMapper mapper = JSON.getMapper();
                com.formance.formance_sdk.models.shared.ConfigsResponse out = mapper.readValue(new String(httpRes.body(), StandardCharsets.UTF_8), com.formance.formance_sdk.models.shared.ConfigsResponse.class);
                res.configsResponse = out;
            }
        }
        else {
            if (com.formance.formance_sdk.utils.Utils.matchContentType(contentType, "application/json")) {
                ObjectMapper mapper = JSON.getMapper();
                com.formance.formance_sdk.models.shared.ErrorResponse out = mapper.readValue(new String(httpRes.body(), StandardCharsets.UTF_8), com.formance.formance_sdk.models.shared.ErrorResponse.class);
                res.errorResponse = out;
            }
        }

        return res;
    }

    /**
     * Insert a new config
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
     * @param request the request object containing all of the parameters for the API call
     * @return the response from the API call
     * @throws Exception if the API call fails
     */
    public com.formance.formance_sdk.models.operations.InsertConfigResponse insertConfig(com.formance.formance_sdk.models.shared.ConfigUser request) throws Exception {
        String baseUrl = this._serverUrl;
        String url = com.formance.formance_sdk.utils.Utils.generateURL(baseUrl, "/api/webhooks/configs");
        
        HTTPRequest req = new HTTPRequest();
        req.setMethod("POST");
        req.setURL(url);
        SerializedBody serializedRequestBody = com.formance.formance_sdk.utils.Utils.serializeRequestBody(request, "request", "json");
        if (serializedRequestBody == null) {
            throw new Exception("Request body is required");
        }
        req.setBody(serializedRequestBody);

        req.addHeader("Accept", "application/json;q=1, application/json;q=0");
        req.addHeader("user-agent", String.format("speakeasy-sdk/%s %s %s", this._language, this._sdkVersion, this._genVersion));
        
        HTTPClient client = this._securityClient;
        
        HttpResponse<byte[]> httpRes = client.send(req);

        String contentType = httpRes.headers().firstValue("Content-Type").orElse("application/octet-stream");

        com.formance.formance_sdk.models.operations.InsertConfigResponse res = new com.formance.formance_sdk.models.operations.InsertConfigResponse(contentType, httpRes.statusCode()) {{
            configResponse = null;
            errorResponse = null;
        }};
        res.rawResponse = httpRes;
        
        if (httpRes.statusCode() == 200) {
            if (com.formance.formance_sdk.utils.Utils.matchContentType(contentType, "application/json")) {
                ObjectMapper mapper = JSON.getMapper();
                com.formance.formance_sdk.models.shared.ConfigResponse out = mapper.readValue(new String(httpRes.body(), StandardCharsets.UTF_8), com.formance.formance_sdk.models.shared.ConfigResponse.class);
                res.configResponse = out;
            }
        }
        else {
            if (com.formance.formance_sdk.utils.Utils.matchContentType(contentType, "application/json")) {
                ObjectMapper mapper = JSON.getMapper();
                com.formance.formance_sdk.models.shared.ErrorResponse out = mapper.readValue(new String(httpRes.body(), StandardCharsets.UTF_8), com.formance.formance_sdk.models.shared.ErrorResponse.class);
                res.errorResponse = out;
            }
        }

        return res;
    }

    /**
     * Test one config
     * Test a config by sending a webhook to its endpoint.
     * @param request the request object containing all of the parameters for the API call
     * @return the response from the API call
     * @throws Exception if the API call fails
     */
    public com.formance.formance_sdk.models.operations.TestConfigResponse testConfig(com.formance.formance_sdk.models.operations.TestConfigRequest request) throws Exception {
        String baseUrl = this._serverUrl;
        String url = com.formance.formance_sdk.utils.Utils.generateURL(com.formance.formance_sdk.models.operations.TestConfigRequest.class, baseUrl, "/api/webhooks/configs/{id}/test", request, null);
        
        HTTPRequest req = new HTTPRequest();
        req.setMethod("GET");
        req.setURL(url);

        req.addHeader("Accept", "application/json;q=1, application/json;q=0");
        req.addHeader("user-agent", String.format("speakeasy-sdk/%s %s %s", this._language, this._sdkVersion, this._genVersion));
        
        HTTPClient client = this._securityClient;
        
        HttpResponse<byte[]> httpRes = client.send(req);

        String contentType = httpRes.headers().firstValue("Content-Type").orElse("application/octet-stream");

        com.formance.formance_sdk.models.operations.TestConfigResponse res = new com.formance.formance_sdk.models.operations.TestConfigResponse(contentType, httpRes.statusCode()) {{
            attemptResponse = null;
            errorResponse = null;
        }};
        res.rawResponse = httpRes;
        
        if (httpRes.statusCode() == 200) {
            if (com.formance.formance_sdk.utils.Utils.matchContentType(contentType, "application/json")) {
                ObjectMapper mapper = JSON.getMapper();
                com.formance.formance_sdk.models.shared.AttemptResponse out = mapper.readValue(new String(httpRes.body(), StandardCharsets.UTF_8), com.formance.formance_sdk.models.shared.AttemptResponse.class);
                res.attemptResponse = out;
            }
        }
        else {
            if (com.formance.formance_sdk.utils.Utils.matchContentType(contentType, "application/json")) {
                ObjectMapper mapper = JSON.getMapper();
                com.formance.formance_sdk.models.shared.ErrorResponse out = mapper.readValue(new String(httpRes.body(), StandardCharsets.UTF_8), com.formance.formance_sdk.models.shared.ErrorResponse.class);
                res.errorResponse = out;
            }
        }

        return res;
    }
}