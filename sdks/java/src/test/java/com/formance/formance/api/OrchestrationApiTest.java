package com.formance.formance.api;

import com.formance.formance.ApiClient;
import com.formance.formance.model.CreateWorkflowResponse;
import com.formance.formance.model.Error;
import com.formance.formance.model.GetWorkflowInstanceHistoryResponse;
import com.formance.formance.model.GetWorkflowInstanceHistoryStageResponse;
import com.formance.formance.model.GetWorkflowInstanceResponse;
import com.formance.formance.model.GetWorkflowResponse;
import com.formance.formance.model.ListRunsResponse;
import com.formance.formance.model.ListWorkflowsResponse;
import com.formance.formance.model.RunWorkflowResponse;
import com.formance.formance.model.SendEventRequest;
import com.formance.formance.model.ServerInfo;
import com.formance.formance.model.WorkflowConfig;
import org.junit.Before;
import org.junit.Test;

import java.time.LocalDate;
import java.time.OffsetDateTime;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * API tests for OrchestrationApi
 */
public class OrchestrationApiTest {

    private OrchestrationApi api;

    @Before
    public void setup() {
        api = new ApiClient().createService(OrchestrationApi.class);
    }

    /**
     * Cancel a running workflow
     *
     * Cancel a running workflow
     */
    @Test
    public void cancelEventTest() {
        String instanceID = null;
        // api.cancelEvent(instanceID);

        // TODO: test validations
    }
    /**
     * Create workflow
     *
     * Create a workflow
     */
    @Test
    public void createWorkflowTest() {
        WorkflowConfig body = null;
        // CreateWorkflowResponse response = api.createWorkflow(body);

        // TODO: test validations
    }
    /**
     * Get a workflow instance by id
     *
     * Get a workflow instance by id
     */
    @Test
    public void getInstanceTest() {
        String instanceID = null;
        // GetWorkflowInstanceResponse response = api.getInstance(instanceID);

        // TODO: test validations
    }
    /**
     * Get a workflow instance history by id
     *
     * Get a workflow instance history by id
     */
    @Test
    public void getInstanceHistoryTest() {
        String instanceID = null;
        // GetWorkflowInstanceHistoryResponse response = api.getInstanceHistory(instanceID);

        // TODO: test validations
    }
    /**
     * Get a workflow instance stage history
     *
     * Get a workflow instance stage history
     */
    @Test
    public void getInstanceStageHistoryTest() {
        String instanceID = null;
        Integer number = null;
        // GetWorkflowInstanceHistoryStageResponse response = api.getInstanceStageHistory(instanceID, number);

        // TODO: test validations
    }
    /**
     * Get a flow by id
     *
     * Get a flow by id
     */
    @Test
    public void getWorkflowTest() {
        String flowId = null;
        // GetWorkflowResponse response = api.getWorkflow(flowId);

        // TODO: test validations
    }
    /**
     * List instances of a workflow
     *
     * List instances of a workflow
     */
    @Test
    public void listInstancesTest() {
        String workflowID = null;
        Boolean running = null;
        // ListRunsResponse response = api.listInstances(workflowID, running);

        // TODO: test validations
    }
    /**
     * List registered workflows
     *
     * List registered workflows
     */
    @Test
    public void listWorkflowsTest() {
        // ListWorkflowsResponse response = api.listWorkflows();

        // TODO: test validations
    }
    /**
     * Get server info
     *
     * 
     */
    @Test
    public void orchestrationgetServerInfoTest() {
        // ServerInfo response = api.orchestrationgetServerInfo();

        // TODO: test validations
    }
    /**
     * Run workflow
     *
     * Run workflow
     */
    @Test
    public void runWorkflowTest() {
        String workflowID = null;
        Boolean wait = null;
        Map<String, String> requestBody = null;
        // RunWorkflowResponse response = api.runWorkflow(workflowID, wait, requestBody);

        // TODO: test validations
    }
    /**
     * Send an event to a running workflow
     *
     * Send an event to a running workflow
     */
    @Test
    public void sendEventTest() {
        String instanceID = null;
        SendEventRequest sendEventRequest = null;
        // api.sendEvent(instanceID, sendEventRequest);

        // TODO: test validations
    }
}
