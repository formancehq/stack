package com.formance.formance.api;

import com.formance.formance.CollectionFormats.*;

import retrofit2.Call;
import retrofit2.http.*;

import okhttp3.RequestBody;
import okhttp3.ResponseBody;
import okhttp3.MultipartBody;

import com.formance.formance.model.AccountsCursor;
import com.formance.formance.model.Connector;
import com.formance.formance.model.ConnectorConfig;
import com.formance.formance.model.ConnectorConfigResponse;
import com.formance.formance.model.ConnectorsConfigsResponse;
import com.formance.formance.model.ConnectorsResponse;
import com.formance.formance.model.PaymentMetadata;
import com.formance.formance.model.PaymentResponse;
import com.formance.formance.model.PaymentsCursor;
import com.formance.formance.model.StripeTransferRequest;
import com.formance.formance.model.TaskResponse;
import com.formance.formance.model.TasksCursor;
import com.formance.formance.model.TransferRequest;
import com.formance.formance.model.TransferResponse;
import com.formance.formance.model.TransfersResponse;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

public interface PaymentsApi {
  /**
   * Transfer funds between Stripe accounts
   * Execute a transfer between two Stripe accounts.
   * @param stripeTransferRequest  (required)
   * @return Call&lt;Object&gt;
   */
  @Headers({
    "Content-Type:application/json"
  })
  @POST("api/payments/connectors/stripe/transfers")
  Call<Object> connectorsStripeTransfer(
    @retrofit2.http.Body StripeTransferRequest stripeTransferRequest
  );

  /**
   * Transfer funds between Connector accounts
   * Execute a transfer between two accounts.
   * @param connector The name of the connector. (required)
   * @param transferRequest  (required)
   * @return Call&lt;TransferResponse&gt;
   */
  @Headers({
    "Content-Type:application/json"
  })
  @POST("api/payments/connectors/{connector}/transfers")
  Call<TransferResponse> connectorsTransfer(
    @retrofit2.http.Path("connector") Connector connector, @retrofit2.http.Body TransferRequest transferRequest
  );

  /**
   * Read a specific task of the connector
   * Get a specific task associated to the connector.
   * @param connector The name of the connector. (required)
   * @param taskId The task ID. (required)
   * @return Call&lt;TaskResponse&gt;
   */
  @GET("api/payments/connectors/{connector}/tasks/{taskId}")
  Call<TaskResponse> getConnectorTask(
    @retrofit2.http.Path("connector") Connector connector, @retrofit2.http.Path("taskId") String taskId
  );

  /**
   * Get a payment
   * 
   * @param paymentId The payment ID. (required)
   * @return Call&lt;PaymentResponse&gt;
   */
  @GET("api/payments/payments/{paymentId}")
  Call<PaymentResponse> getPayment(
    @retrofit2.http.Path("paymentId") String paymentId
  );

  /**
   * Install a connector
   * Install a connector by its name and config.
   * @param connector The name of the connector. (required)
   * @param connectorConfig  (required)
   * @return Call&lt;Void&gt;
   */
  @Headers({
    "Content-Type:application/json"
  })
  @POST("api/payments/connectors/{connector}")
  Call<Void> installConnector(
    @retrofit2.http.Path("connector") Connector connector, @retrofit2.http.Body ConnectorConfig connectorConfig
  );

  /**
   * List all installed connectors
   * List all installed connectors.
   * @return Call&lt;ConnectorsResponse&gt;
   */
  @GET("api/payments/connectors")
  Call<ConnectorsResponse> listAllConnectors();
    

  /**
   * List the configs of each available connector
   * List the configs of each available connector.
   * @return Call&lt;ConnectorsConfigsResponse&gt;
   */
  @GET("api/payments/connectors/configs")
  Call<ConnectorsConfigsResponse> listConfigsAvailableConnectors();
    

  /**
   * List tasks from a connector
   * List all tasks associated with this connector.
   * @param connector The name of the connector. (required)
   * @param pageSize The maximum number of results to return per page.  (optional, default to 15)
   * @param cursor Parameter used in pagination requests. Maximum page size is set to 15. Set to the value of next for the next page of results. Set to the value of previous for the previous page of results. No other parameters can be set when this parameter is set.  (optional)
   * @return Call&lt;TasksCursor&gt;
   */
  @GET("api/payments/connectors/{connector}/tasks")
  Call<TasksCursor> listConnectorTasks(
    @retrofit2.http.Path("connector") Connector connector, @retrofit2.http.Query("pageSize") Long pageSize, @retrofit2.http.Query("cursor") String cursor
  );

  /**
   * List transfers and their statuses
   * List transfers
   * @param connector The name of the connector. (required)
   * @return Call&lt;TransfersResponse&gt;
   */
  @GET("api/payments/connectors/{connector}/transfers")
  Call<TransfersResponse> listConnectorsTransfers(
    @retrofit2.http.Path("connector") Connector connector
  );

  /**
   * List payments
   * 
   * @param pageSize The maximum number of results to return per page.  (optional, default to 15)
   * @param cursor Parameter used in pagination requests. Maximum page size is set to 15. Set to the value of next for the next page of results. Set to the value of previous for the previous page of results. No other parameters can be set when this parameter is set.  (optional)
   * @param sort Fields used to sort payments (default is date:desc). (optional)
   * @return Call&lt;PaymentsCursor&gt;
   */
  @GET("api/payments/payments")
  Call<PaymentsCursor> listPayments(
    @retrofit2.http.Query("pageSize") Long pageSize, @retrofit2.http.Query("cursor") String cursor, @retrofit2.http.Query("sort") List<String> sort
  );

  /**
   * List accounts
   * 
   * @param pageSize The maximum number of results to return per page.  (optional, default to 15)
   * @param cursor Parameter used in pagination requests. Maximum page size is set to 15. Set to the value of next for the next page of results. Set to the value of previous for the previous page of results. No other parameters can be set when this parameter is set.  (optional)
   * @param sort Fields used to sort payments (default is date:desc). (optional)
   * @return Call&lt;AccountsCursor&gt;
   */
  @GET("api/payments/accounts")
  Call<AccountsCursor> paymentslistAccounts(
    @retrofit2.http.Query("pageSize") Long pageSize, @retrofit2.http.Query("cursor") String cursor, @retrofit2.http.Query("sort") List<String> sort
  );

  /**
   * Read the config of a connector
   * Read connector config
   * @param connector The name of the connector. (required)
   * @return Call&lt;ConnectorConfigResponse&gt;
   */
  @GET("api/payments/connectors/{connector}/config")
  Call<ConnectorConfigResponse> readConnectorConfig(
    @retrofit2.http.Path("connector") Connector connector
  );

  /**
   * Reset a connector
   * Reset a connector by its name. It will remove the connector and ALL PAYMENTS generated with it. 
   * @param connector The name of the connector. (required)
   * @return Call&lt;Void&gt;
   */
  @POST("api/payments/connectors/{connector}/reset")
  Call<Void> resetConnector(
    @retrofit2.http.Path("connector") Connector connector
  );

  /**
   * Uninstall a connector
   * Uninstall a connector by its name.
   * @param connector The name of the connector. (required)
   * @return Call&lt;Void&gt;
   */
  @DELETE("api/payments/connectors/{connector}")
  Call<Void> uninstallConnector(
    @retrofit2.http.Path("connector") Connector connector
  );

  /**
   * Update metadata
   * 
   * @param paymentId The payment ID. (required)
   * @param paymentMetadata  (required)
   * @return Call&lt;Void&gt;
   */
  @Headers({
    "Content-Type:application/json"
  })
  @PATCH("api/payments/payments/{paymentId}/metadata")
  Call<Void> updateMetadata(
    @retrofit2.http.Path("paymentId") String paymentId, @retrofit2.http.Body PaymentMetadata paymentMetadata
  );

}
