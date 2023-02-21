package com.formance.formance.api;

import com.formance.formance.ApiClient;
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
import org.junit.Before;
import org.junit.Test;

import java.time.LocalDate;
import java.time.OffsetDateTime;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * API tests for PaymentsApi
 */
public class PaymentsApiTest {

    private PaymentsApi api;

    @Before
    public void setup() {
        api = new ApiClient().createService(PaymentsApi.class);
    }

    /**
     * Transfer funds between Stripe accounts
     *
     * Execute a transfer between two Stripe accounts.
     */
    @Test
    public void connectorsStripeTransferTest() {
        StripeTransferRequest stripeTransferRequest = null;
        // Object response = api.connectorsStripeTransfer(stripeTransferRequest);

        // TODO: test validations
    }
    /**
     * Transfer funds between Connector accounts
     *
     * Execute a transfer between two accounts.
     */
    @Test
    public void connectorsTransferTest() {
        Connector connector = null;
        TransferRequest transferRequest = null;
        // TransferResponse response = api.connectorsTransfer(connector, transferRequest);

        // TODO: test validations
    }
    /**
     * Read a specific task of the connector
     *
     * Get a specific task associated to the connector.
     */
    @Test
    public void getConnectorTaskTest() {
        Connector connector = null;
        String taskId = null;
        // TaskResponse response = api.getConnectorTask(connector, taskId);

        // TODO: test validations
    }
    /**
     * Get a payment
     *
     * 
     */
    @Test
    public void getPaymentTest() {
        String paymentId = null;
        // PaymentResponse response = api.getPayment(paymentId);

        // TODO: test validations
    }
    /**
     * Install a connector
     *
     * Install a connector by its name and config.
     */
    @Test
    public void installConnectorTest() {
        Connector connector = null;
        ConnectorConfig connectorConfig = null;
        // api.installConnector(connector, connectorConfig);

        // TODO: test validations
    }
    /**
     * List all installed connectors
     *
     * List all installed connectors.
     */
    @Test
    public void listAllConnectorsTest() {
        // ConnectorsResponse response = api.listAllConnectors();

        // TODO: test validations
    }
    /**
     * List the configs of each available connector
     *
     * List the configs of each available connector.
     */
    @Test
    public void listConfigsAvailableConnectorsTest() {
        // ConnectorsConfigsResponse response = api.listConfigsAvailableConnectors();

        // TODO: test validations
    }
    /**
     * List tasks from a connector
     *
     * List all tasks associated with this connector.
     */
    @Test
    public void listConnectorTasksTest() {
        Connector connector = null;
        Long pageSize = null;
        String cursor = null;
        // TasksCursor response = api.listConnectorTasks(connector, pageSize, cursor);

        // TODO: test validations
    }
    /**
     * List transfers and their statuses
     *
     * List transfers
     */
    @Test
    public void listConnectorsTransfersTest() {
        Connector connector = null;
        // TransfersResponse response = api.listConnectorsTransfers(connector);

        // TODO: test validations
    }
    /**
     * List payments
     *
     * 
     */
    @Test
    public void listPaymentsTest() {
        Long pageSize = null;
        String cursor = null;
        List<String> sort = null;
        // PaymentsCursor response = api.listPayments(pageSize, cursor, sort);

        // TODO: test validations
    }
    /**
     * List accounts
     *
     * 
     */
    @Test
    public void paymentslistAccountsTest() {
        Long pageSize = null;
        String cursor = null;
        List<String> sort = null;
        // AccountsCursor response = api.paymentslistAccounts(pageSize, cursor, sort);

        // TODO: test validations
    }
    /**
     * Read the config of a connector
     *
     * Read connector config
     */
    @Test
    public void readConnectorConfigTest() {
        Connector connector = null;
        // ConnectorConfigResponse response = api.readConnectorConfig(connector);

        // TODO: test validations
    }
    /**
     * Reset a connector
     *
     * Reset a connector by its name. It will remove the connector and ALL PAYMENTS generated with it. 
     */
    @Test
    public void resetConnectorTest() {
        Connector connector = null;
        // api.resetConnector(connector);

        // TODO: test validations
    }
    /**
     * Uninstall a connector
     *
     * Uninstall a connector by its name.
     */
    @Test
    public void uninstallConnectorTest() {
        Connector connector = null;
        // api.uninstallConnector(connector);

        // TODO: test validations
    }
    /**
     * Update metadata
     *
     * 
     */
    @Test
    public void updateMetadataTest() {
        String paymentId = null;
        PaymentMetadata paymentMetadata = null;
        // api.updateMetadata(paymentId, paymentMetadata);

        // TODO: test validations
    }
}
