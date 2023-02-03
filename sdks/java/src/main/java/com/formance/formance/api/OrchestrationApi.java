package com.formance.formance.api;

import com.formance.formance.CollectionFormats.*;

import retrofit2.Call;
import retrofit2.http.*;

import okhttp3.RequestBody;
import okhttp3.ResponseBody;
import okhttp3.MultipartBody;

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

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

public interface OrchestrationApi {
  /**
   * Cancel a running workflow
   * Cancel a running workflow
   * @param instanceID The instance id (required)
   * @return Call&lt;Void&gt;
   */
  @PUT("api/orchestration/instances/{instanceID}/abort")
  Call<Void> cancelEvent(
    @retrofit2.http.Path("instanceID") String instanceID
  );

  /**
   * Create workflow
   * Create a workflow
   * @param body  (optional)
   * @return Call&lt;CreateWorkflowResponse&gt;
   */
  @Headers({
    "Content-Type:application/json"
  })
  @POST("api/orchestration/workflows")
  Call<CreateWorkflowResponse> createWorkflow(
    @retrofit2.http.Body WorkflowConfig body
  );

  /**
   * Get a workflow instance by id
   * Get a workflow instance by id
   * @param instanceID The instance id (required)
   * @return Call&lt;GetWorkflowInstanceResponse&gt;
   */
  @GET("api/orchestration/instances/{instanceID}")
  Call<GetWorkflowInstanceResponse> getInstance(
    @retrofit2.http.Path("instanceID") String instanceID
  );

  /**
   * Get a workflow instance history by id
   * Get a workflow instance history by id
   * @param instanceID The instance id (required)
   * @return Call&lt;GetWorkflowInstanceHistoryResponse&gt;
   */
  @GET("api/orchestration/instances/{instanceID}/history")
  Call<GetWorkflowInstanceHistoryResponse> getInstanceHistory(
    @retrofit2.http.Path("instanceID") String instanceID
  );

  /**
   * Get a workflow instance stage history
   * Get a workflow instance stage history
   * @param instanceID The instance id (required)
   * @param number The stage number (required)
   * @return Call&lt;GetWorkflowInstanceHistoryStageResponse&gt;
   */
  @GET("api/orchestration/instances/{instanceID}/stages/{number}/history")
  Call<GetWorkflowInstanceHistoryStageResponse> getInstanceStageHistory(
    @retrofit2.http.Path("instanceID") String instanceID, @retrofit2.http.Path("number") Integer number
  );

  /**
   * Get a flow by id
   * Get a flow by id
   * @param flowId The flow id (required)
   * @return Call&lt;GetWorkflowResponse&gt;
   */
  @GET("api/orchestration/workflows/{flowId}")
  Call<GetWorkflowResponse> getWorkflow(
    @retrofit2.http.Path("flowId") String flowId
  );

  /**
   * List instances of a workflow
   * List instances of a workflow
   * @param workflowID A workflow id (optional)
   * @param running Filter running instances (optional)
   * @return Call&lt;ListRunsResponse&gt;
   */
  @GET("api/orchestration/instances")
  Call<ListRunsResponse> listInstances(
    @retrofit2.http.Query("workflowID") String workflowID, @retrofit2.http.Query("running") Boolean running
  );

  /**
   * List registered workflows
   * List registered workflows
   * @return Call&lt;ListWorkflowsResponse&gt;
   */
  @GET("api/orchestration/workflows")
  Call<ListWorkflowsResponse> listWorkflows();
    

  /**
   * Get server info
   * 
   * @return Call&lt;ServerInfo&gt;
   */
  @GET("api/orchestration/_info")
  Call<ServerInfo> orchestrationgetServerInfo();
    

  /**
   * Run workflow
   * Run workflow
   * @param workflowID The flow id (required)
   * @param wait Wait end of the workflow before return (optional)
   * @param requestBody  (optional)
   * @return Call&lt;RunWorkflowResponse&gt;
   */
  @Headers({
    "Content-Type:application/json"
  })
  @POST("api/orchestration/workflows/{workflowID}/instances")
  Call<RunWorkflowResponse> runWorkflow(
    @retrofit2.http.Path("workflowID") String workflowID, @retrofit2.http.Query("wait") Boolean wait, @retrofit2.http.Body Map<String, String> requestBody
  );

  /**
   * Send an event to a running workflow
   * Send an event to a running workflow
   * @param instanceID The instance id (required)
   * @param sendEventRequest  (optional)
   * @return Call&lt;Void&gt;
   */
  @Headers({
    "Content-Type:application/json"
  })
  @POST("api/orchestration/instances/{instanceID}/events")
  Call<Void> sendEvent(
    @retrofit2.http.Path("instanceID") String instanceID, @retrofit2.http.Body SendEventRequest sendEventRequest
  );

}
