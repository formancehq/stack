package com.formance.formance.api;

import com.formance.formance.ApiClient;
import com.formance.formance.model.ErrorResponse;
import java.time.OffsetDateTime;
import com.formance.formance.model.PostTransaction;
import com.formance.formance.model.TransactionResponse;
import com.formance.formance.model.TransactionsCursorResponse;
import org.junit.Before;
import org.junit.Test;

import java.time.LocalDate;
import java.time.OffsetDateTime;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * API tests for TransactionsApi
 */
public class TransactionsApiTest {

    private TransactionsApi api;

    @Before
    public void setup() {
        api = new ApiClient().createService(TransactionsApi.class);
    }

    /**
     * Set the metadata of a transaction by its ID
     *
     * 
     */
    @Test
    public void addMetadataOnTransactionTest() {
        String ledger = null;
        Long txid = null;
        Map<String, Object> requestBody = null;
        // api.addMetadataOnTransaction(ledger, txid, requestBody);

        // TODO: test validations
    }
    /**
     * Count the transactions from a ledger
     *
     * 
     */
    @Test
    public void countTransactionsTest() {
        String ledger = null;
        String reference = null;
        String account = null;
        String source = null;
        String destination = null;
        OffsetDateTime startTime = null;
        OffsetDateTime endTime = null;
        Object metadata = null;
        // api.countTransactions(ledger, reference, account, source, destination, startTime, endTime, metadata);

        // TODO: test validations
    }
    /**
     * Create a new transaction to a ledger
     *
     * 
     */
    @Test
    public void createTransactionTest() {
        String ledger = null;
        PostTransaction postTransaction = null;
        Boolean preview = null;
        // TransactionResponse response = api.createTransaction(ledger, postTransaction, preview);

        // TODO: test validations
    }
    /**
     * Get transaction from a ledger by its ID
     *
     * 
     */
    @Test
    public void getTransactionTest() {
        String ledger = null;
        Long txid = null;
        // TransactionResponse response = api.getTransaction(ledger, txid);

        // TODO: test validations
    }
    /**
     * List transactions from a ledger
     *
     * List transactions from a ledger, sorted by txid in descending order.
     */
    @Test
    public void listTransactionsTest() {
        String ledger = null;
        Long pageSize = null;
        String after = null;
        String reference = null;
        String account = null;
        String source = null;
        String destination = null;
        OffsetDateTime startTime = null;
        OffsetDateTime endTime = null;
        String cursor = null;
        Object metadata = null;
        // TransactionsCursorResponse response = api.listTransactions(ledger, pageSize, after, reference, account, source, destination, startTime, endTime, cursor, metadata);

        // TODO: test validations
    }
    /**
     * Revert a ledger transaction by its ID
     *
     * 
     */
    @Test
    public void revertTransactionTest() {
        String ledger = null;
        Long txid = null;
        // TransactionResponse response = api.revertTransaction(ledger, txid);

        // TODO: test validations
    }
}
