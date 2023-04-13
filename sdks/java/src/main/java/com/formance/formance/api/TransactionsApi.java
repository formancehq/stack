package com.formance.formance.api;

import com.formance.formance.CollectionFormats.*;

import retrofit2.Call;
import retrofit2.http.*;

import okhttp3.RequestBody;
import okhttp3.ResponseBody;
import okhttp3.MultipartBody;

import com.formance.formance.model.ErrorResponse;
import java.time.OffsetDateTime;
import com.formance.formance.model.PostTransaction;
import com.formance.formance.model.TransactionResponse;
import com.formance.formance.model.TransactionsCursorResponse;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

public interface TransactionsApi {
  /**
   * Set the metadata of a transaction by its ID
   * 
   * @param ledger Name of the ledger. (required)
   * @param txid Transaction ID. (required)
   * @param requestBody metadata (optional)
   * @return Call&lt;Void&gt;
   */
  @Headers({
    "Content-Type:application/json"
  })
  @POST("api/ledger/{ledger}/transactions/{txid}/metadata")
  Call<Void> addMetadataOnTransaction(
    @retrofit2.http.Path("ledger") String ledger, @retrofit2.http.Path("txid") Long txid, @retrofit2.http.Body Map<String, String> requestBody
  );

  /**
   * Count the transactions from a ledger
   * 
   * @param ledger Name of the ledger. (required)
   * @param reference Filter transactions by reference field. (optional)
   * @param account Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $). (optional)
   * @param source Filter transactions with postings involving given account at source (regular expression placed between ^ and $). (optional)
   * @param destination Filter transactions with postings involving given account at destination (regular expression placed between ^ and $). (optional)
   * @param startTime Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; includes the first second of 4th minute).  (optional)
   * @param endTime Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; excludes the first second of 4th minute).  (optional)
   * @param metadata Filter transactions by metadata key value pairs. Nested objects can be used as seen in the example below. (optional)
   * @return Call&lt;Void&gt;
   */
  @HEAD("api/ledger/{ledger}/transactions")
  Call<Void> countTransactions(
    @retrofit2.http.Path("ledger") String ledger, @retrofit2.http.Query("reference") String reference, @retrofit2.http.Query("account") String account, @retrofit2.http.Query("source") String source, @retrofit2.http.Query("destination") String destination, @retrofit2.http.Query("startTime") OffsetDateTime startTime, @retrofit2.http.Query("endTime") OffsetDateTime endTime, @retrofit2.http.Query("metadata") Map<String, String> metadata
  );

  /**
   * Create a new transaction to a ledger
   * 
   * @param ledger Name of the ledger. (required)
   * @param postTransaction The request body must contain at least one of the following objects:   - &#x60;postings&#x60;: suitable for simple transactions   - &#x60;script&#x60;: enabling more complex transactions with Numscript  (required)
   * @param preview Set the preview mode. Preview mode doesn&#39;t add the logs to the database or publish a message to the message broker. (optional)
   * @return Call&lt;TransactionResponse&gt;
   */
  @Headers({
    "Content-Type:application/json"
  })
  @POST("api/ledger/{ledger}/transactions")
  Call<TransactionResponse> createTransaction(
    @retrofit2.http.Path("ledger") String ledger, @retrofit2.http.Body PostTransaction postTransaction, @retrofit2.http.Query("preview") Boolean preview
  );

  /**
   * Get transaction from a ledger by its ID
   * 
   * @param ledger Name of the ledger. (required)
   * @param txid Transaction ID. (required)
   * @return Call&lt;TransactionResponse&gt;
   */
  @GET("api/ledger/{ledger}/transactions/{txid}")
  Call<TransactionResponse> getTransaction(
    @retrofit2.http.Path("ledger") String ledger, @retrofit2.http.Path("txid") Long txid
  );

  /**
   * List transactions from a ledger
   * List transactions from a ledger, sorted by txid in descending order.
   * @param ledger Name of the ledger. (required)
   * @param pageSize The maximum number of results to return per page.  (optional)
   * @param after Pagination cursor, will return transactions after given txid (in descending order). (optional)
   * @param reference Find transactions by reference field. (optional)
   * @param account Filter transactions with postings involving given account, either as source or destination (regular expression placed between ^ and $). (optional)
   * @param source Filter transactions with postings involving given account at source (regular expression placed between ^ and $). (optional)
   * @param destination Filter transactions with postings involving given account at destination (regular expression placed between ^ and $). (optional)
   * @param startTime Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; includes the first second of 4th minute).  (optional)
   * @param endTime Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; excludes the first second of 4th minute).  (optional)
   * @param cursor Parameter used in pagination requests. Maximum page size is set to 15. Set to the value of next for the next page of results. Set to the value of previous for the previous page of results. No other parameters can be set when this parameter is set.  (optional)
   * @param metadata Filter transactions by metadata key value pairs. (optional)
   * @return Call&lt;TransactionsCursorResponse&gt;
   */
  @GET("api/ledger/{ledger}/transactions")
  Call<TransactionsCursorResponse> listTransactions(
    @retrofit2.http.Path("ledger") String ledger, @retrofit2.http.Query("pageSize") Long pageSize, @retrofit2.http.Query("after") String after, @retrofit2.http.Query("reference") String reference, @retrofit2.http.Query("account") String account, @retrofit2.http.Query("source") String source, @retrofit2.http.Query("destination") String destination, @retrofit2.http.Query("startTime") OffsetDateTime startTime, @retrofit2.http.Query("endTime") OffsetDateTime endTime, @retrofit2.http.Query("cursor") String cursor, @retrofit2.http.Query("metadata") Map<String, String> metadata
  );

  /**
   * Revert a ledger transaction by its ID
   * 
   * @param ledger Name of the ledger. (required)
   * @param txid Transaction ID. (required)
   * @return Call&lt;TransactionResponse&gt;
   */
  @POST("api/ledger/{ledger}/transactions/{txid}/revert")
  Call<TransactionResponse> revertTransaction(
    @retrofit2.http.Path("ledger") String ledger, @retrofit2.http.Path("txid") Long txid
  );

}
