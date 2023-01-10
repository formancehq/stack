package com.formance.formance.api;

import com.formance.formance.CollectionFormats.*;

import retrofit2.Call;
import retrofit2.http.*;

import okhttp3.RequestBody;
import okhttp3.ResponseBody;
import okhttp3.MultipartBody;

import com.formance.formance.model.ConnectorConfig;
import com.formance.formance.model.Connectors;
import com.formance.formance.model.ListAccountsResponse;
import com.formance.formance.model.ListConnectorTasks200ResponseInner;
import com.formance.formance.model.ListConnectorsConfigsResponse;
import com.formance.formance.model.ListConnectorsResponse;
import com.formance.formance.model.ListPaymentsResponse;
import com.formance.formance.model.Payment;
import com.formance.formance.model.StripeTransferRequest;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public interface PaymentsApi {
  /**
   * Transfer funds between Stripe accounts
   * Execute a transfer between two Stripe accounts
   * @param stripeTransferRequest  (required)
   * @return Call&lt;Void&gt;
   */
  @Headers({
    "Content-Type:application/json"
  })
  @POST("api/payments/connectors/stripe/transfer")
  Call<Void> connectorsStripeTransfer(
    @retrofit2.http.Body StripeTransferRequest stripeTransferRequest
  );

  /**
   * Get all installed connectors
   * Get all installed connectors
   * @return Call&lt;ListConnectorsResponse&gt;
   */
  @GET("api/payments/connectors")
  Call<ListConnectorsResponse> getAllConnectors();
    

  /**
   * Get all available connectors configs
   * Get all available connectors configs
   * @return Call&lt;ListConnectorsConfigsResponse&gt;
   */
  @GET("api/payments/connectors/configs")
  Call<ListConnectorsConfigsResponse> getAllConnectorsConfigs();
    

  /**
   * Read a specific task of the connector
   * Get a specific task associated to the connector
   * @param connector The connector code (required)
   * @param taskId The task id (required)
   * @return Call&lt;ListConnectorTasks200ResponseInner&gt;
   */
  @GET("api/payments/connectors/{connector}/tasks/{taskId}")
  Call<ListConnectorTasks200ResponseInner> getConnectorTask(
    @retrofit2.http.Path("connector") Connectors connector, @retrofit2.http.Path("taskId") String taskId
  );

  /**
   * Returns a payment.
   * 
   * @param paymentId The payment id (required)
   * @return Call&lt;Payment&gt;
   */
  @GET("api/payments/payments/{paymentId}")
  Call<Payment> getPayment(
    @retrofit2.http.Path("paymentId") String paymentId
  );

  /**
   * Install connector
   * Install connector
   * @param connector The connector code (required)
   * @param connectorConfig  (required)
   * @return Call&lt;Void&gt;
   */
  @Headers({
    "Content-Type:application/json"
  })
  @POST("api/payments/connectors/{connector}")
  Call<Void> installConnector(
    @retrofit2.http.Path("connector") Connectors connector, @retrofit2.http.Body ConnectorConfig connectorConfig
  );

  /**
   * List connector tasks
   * List all tasks associated with this connector.
   * @param connector The connector code (required)
   * @return Call&lt;List&lt;ListConnectorTasks200ResponseInner&gt;&gt;
   */
  @GET("api/payments/connectors/{connector}/tasks")
  Call<List<ListConnectorTasks200ResponseInner>> listConnectorTasks(
    @retrofit2.http.Path("connector") Connectors connector
  );

  /**
   * Returns a list of payments.
   * 
   * @param limit Limit the number of payments to return, pagination can be achieved in conjunction with &#39;skip&#39; parameter. (optional)
   * @param skip How many payments to skip, pagination can be achieved in conjunction with &#39;limit&#39; parameter. (optional)
   * @param sort Field used to sort payments (Default is by date). (optional)
   * @return Call&lt;ListPaymentsResponse&gt;
   */
  @GET("api/payments/payments")
  Call<ListPaymentsResponse> listPayments(
    @retrofit2.http.Query("limit") Integer limit, @retrofit2.http.Query("skip") Integer skip, @retrofit2.http.Query("sort") List<String> sort
  );

  /**
   * Returns a list of accounts.
   * 
   * @param limit Limit the number of accounts to return, pagination can be achieved in conjunction with &#39;skip&#39; parameter. (optional)
   * @param skip How many accounts to skip, pagination can be achieved in conjunction with &#39;limit&#39; parameter. (optional)
   * @param sort Field used to sort payments (Default is by date). (optional)
   * @return Call&lt;ListAccountsResponse&gt;
   */
  @GET("api/payments/accounts")
  Call<ListAccountsResponse> paymentslistAccounts(
    @retrofit2.http.Query("limit") Integer limit, @retrofit2.http.Query("skip") Integer skip, @retrofit2.http.Query("sort") List<String> sort
  );

  /**
   * Read connector config
   * Read connector config
   * @param connector The connector code (required)
   * @return Call&lt;ConnectorConfig&gt;
   */
  @GET("api/payments/connectors/{connector}/config")
  Call<ConnectorConfig> readConnectorConfig(
    @retrofit2.http.Path("connector") Connectors connector
  );

  /**
   * Reset connector
   * Reset connector. Will remove the connector and ALL PAYMENTS generated with it.
   * @param connector The connector code (required)
   * @return Call&lt;Void&gt;
   */
  @POST("api/payments/connectors/{connector}/reset")
  Call<Void> resetConnector(
    @retrofit2.http.Path("connector") Connectors connector
  );

  /**
   * Uninstall connector
   * Uninstall  connector
   * @param connector The connector code (required)
   * @return Call&lt;Void&gt;
   */
  @DELETE("api/payments/connectors/{connector}")
  Call<Void> uninstallConnector(
    @retrofit2.http.Path("connector") Connectors connector
  );

}
