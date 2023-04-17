package com.formance.formance.api;

import com.formance.formance.CollectionFormats.*;

import retrofit2.Call;
import retrofit2.http.*;

import okhttp3.RequestBody;
import okhttp3.ResponseBody;
import okhttp3.MultipartBody;

import com.formance.formance.model.ErrorResponse;
import com.formance.formance.model.LogsCursorResponse;
import java.time.OffsetDateTime;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

public interface LogsApi {
  /**
   * List the logs from a ledger
   * List the logs from a ledger, sorted by ID in descending order.
   * @param ledger Name of the ledger. (required)
   * @param pageSize The maximum number of results to return per page.  (optional)
   * @param startTime Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; includes the first second of 4th minute).  (optional)
   * @param endTime Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; excludes the first second of 4th minute).  (optional)
   * @param cursor Parameter used in pagination requests. Maximum page size is set to 15. Set to the value of next for the next page of results. Set to the value of previous for the previous page of results. No other parameters can be set when this parameter is set.  (optional)
   * @return Call&lt;LogsCursorResponse&gt;
   */
  @GET("api/ledger/{ledger}/logs")
  Call<LogsCursorResponse> listLogs(
    @retrofit2.http.Path("ledger") String ledger, @retrofit2.http.Query("pageSize") Long pageSize, @retrofit2.http.Query("startTime") OffsetDateTime startTime, @retrofit2.http.Query("endTime") OffsetDateTime endTime, @retrofit2.http.Query("cursor") String cursor
  );

}
