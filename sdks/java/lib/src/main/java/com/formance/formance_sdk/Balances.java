/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.formance.formance_sdk.utils.HTTPClient;
import com.formance.formance_sdk.utils.HTTPRequest;
import com.formance.formance_sdk.utils.JSON;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import org.apache.http.NameValuePair;

public class Balances {
	
	private HTTPClient _defaultClient;
	private HTTPClient _securityClient;
	private String _serverUrl;
	private String _language;
	private String _sdkVersion;
	private String _genVersion;

	public Balances(HTTPClient defaultClient, HTTPClient securityClient, String serverUrl, String language, String sdkVersion, String genVersion) {
		this._defaultClient = defaultClient;
		this._securityClient = securityClient;
		this._serverUrl = serverUrl;
		this._language = language;
		this._sdkVersion = sdkVersion;
		this._genVersion = genVersion;
	}

    /**
     * Get the balances from a ledger's account
     * @param request the request object containing all of the parameters for the API call
     * @return the response from the API call
     * @throws Exception if the API call fails
     */
    public com.formance.formance_sdk.models.operations.GetBalancesResponse getBalances(com.formance.formance_sdk.models.operations.GetBalancesRequest request) throws Exception {
        String baseUrl = this._serverUrl;
        String url = com.formance.formance_sdk.utils.Utils.generateURL(com.formance.formance_sdk.models.operations.GetBalancesRequest.class, baseUrl, "/api/ledger/{ledger}/balances", request, null);
        
        HTTPRequest req = new HTTPRequest();
        req.setMethod("GET");
        req.setURL(url);

        req.addHeader("Accept", "application/json;q=1, application/json;q=0");
        req.addHeader("user-agent", String.format("speakeasy-sdk/%s %s %s", this._language, this._sdkVersion, this._genVersion));
        java.util.List<NameValuePair> queryParams = com.formance.formance_sdk.utils.Utils.getQueryParams(com.formance.formance_sdk.models.operations.GetBalancesRequest.class, request, null);
        if (queryParams != null) {
            for (NameValuePair queryParam : queryParams) {
                req.addQueryParam(queryParam);
            }
        }
        
        HTTPClient client = this._securityClient;
        
        HttpResponse<byte[]> httpRes = client.send(req);

        String contentType = httpRes.headers().firstValue("Content-Type").orElse("application/octet-stream");

        com.formance.formance_sdk.models.operations.GetBalancesResponse res = new com.formance.formance_sdk.models.operations.GetBalancesResponse(contentType, httpRes.statusCode()) {{
            balancesCursorResponse = null;
            errorResponse = null;
        }};
        res.rawResponse = httpRes;
        
        if (httpRes.statusCode() == 200) {
            if (com.formance.formance_sdk.utils.Utils.matchContentType(contentType, "application/json")) {
                ObjectMapper mapper = JSON.getMapper();
                com.formance.formance_sdk.models.shared.BalancesCursorResponse out = mapper.readValue(new String(httpRes.body(), StandardCharsets.UTF_8), com.formance.formance_sdk.models.shared.BalancesCursorResponse.class);
                res.balancesCursorResponse = out;
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
     * Get the aggregated balances from selected accounts
     * @param request the request object containing all of the parameters for the API call
     * @return the response from the API call
     * @throws Exception if the API call fails
     */
    public com.formance.formance_sdk.models.operations.GetBalancesAggregatedResponse getBalancesAggregated(com.formance.formance_sdk.models.operations.GetBalancesAggregatedRequest request) throws Exception {
        String baseUrl = this._serverUrl;
        String url = com.formance.formance_sdk.utils.Utils.generateURL(com.formance.formance_sdk.models.operations.GetBalancesAggregatedRequest.class, baseUrl, "/api/ledger/{ledger}/aggregate/balances", request, null);
        
        HTTPRequest req = new HTTPRequest();
        req.setMethod("GET");
        req.setURL(url);

        req.addHeader("Accept", "application/json;q=1, application/json;q=0");
        req.addHeader("user-agent", String.format("speakeasy-sdk/%s %s %s", this._language, this._sdkVersion, this._genVersion));
        java.util.List<NameValuePair> queryParams = com.formance.formance_sdk.utils.Utils.getQueryParams(com.formance.formance_sdk.models.operations.GetBalancesAggregatedRequest.class, request, null);
        if (queryParams != null) {
            for (NameValuePair queryParam : queryParams) {
                req.addQueryParam(queryParam);
            }
        }
        
        HTTPClient client = this._securityClient;
        
        HttpResponse<byte[]> httpRes = client.send(req);

        String contentType = httpRes.headers().firstValue("Content-Type").orElse("application/octet-stream");

        com.formance.formance_sdk.models.operations.GetBalancesAggregatedResponse res = new com.formance.formance_sdk.models.operations.GetBalancesAggregatedResponse(contentType, httpRes.statusCode()) {{
            aggregateBalancesResponse = null;
            errorResponse = null;
        }};
        res.rawResponse = httpRes;
        
        if (httpRes.statusCode() == 200) {
            if (com.formance.formance_sdk.utils.Utils.matchContentType(contentType, "application/json")) {
                ObjectMapper mapper = JSON.getMapper();
                com.formance.formance_sdk.models.shared.AggregateBalancesResponse out = mapper.readValue(new String(httpRes.body(), StandardCharsets.UTF_8), com.formance.formance_sdk.models.shared.AggregateBalancesResponse.class);
                res.aggregateBalancesResponse = out;
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