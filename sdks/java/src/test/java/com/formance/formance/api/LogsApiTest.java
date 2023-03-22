package com.formance.formance.api;

import com.formance.formance.ApiClient;
import com.formance.formance.model.ErrorResponse;
import com.formance.formance.model.LogsCursorResponse;
import java.time.OffsetDateTime;
import org.junit.Before;
import org.junit.Test;

import java.time.LocalDate;
import java.time.OffsetDateTime;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * API tests for LogsApi
 */
public class LogsApiTest {

    private LogsApi api;

    @Before
    public void setup() {
        api = new ApiClient().createService(LogsApi.class);
    }

    /**
     * List the logs from a ledger
     *
     * List the logs from a ledger, sorted by ID in descending order.
     */
    @Test
    public void listLogsTest() {
        String ledger = null;
        Long pageSize = null;
        String after = null;
        OffsetDateTime startTime = null;
        OffsetDateTime endTime = null;
        String cursor = null;
        // LogsCursorResponse response = api.listLogs(ledger, pageSize, after, startTime, endTime, cursor);

        // TODO: test validations
    }
}
