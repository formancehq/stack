/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.shared;

import com.fasterxml.jackson.annotation.JsonInclude.Include;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.databind.annotation.JsonDeserialize;
import com.fasterxml.jackson.databind.annotation.JsonSerialize;
import com.formance.formance_sdk.utils.DateTimeDeserializer;
import com.formance.formance_sdk.utils.DateTimeSerializer;
import java.time.OffsetDateTime;

public class WorkflowInstance {
    @JsonSerialize(using = DateTimeSerializer.class)
    @JsonDeserialize(using = DateTimeDeserializer.class)
    @JsonProperty("createdAt")
    public OffsetDateTime createdAt;

    public WorkflowInstance withCreatedAt(OffsetDateTime createdAt) {
        this.createdAt = createdAt;
        return this;
    }
    
    @JsonInclude(Include.NON_ABSENT)
    @JsonProperty("error")
    public String error;

    public WorkflowInstance withError(String error) {
        this.error = error;
        return this;
    }
    
    @JsonProperty("id")
    public String id;

    public WorkflowInstance withId(String id) {
        this.id = id;
        return this;
    }
    
    @JsonInclude(Include.NON_ABSENT)
    @JsonProperty("status")
    public StageStatus[] status;

    public WorkflowInstance withStatus(StageStatus[] status) {
        this.status = status;
        return this;
    }
    
    @JsonProperty("terminated")
    public Boolean terminated;

    public WorkflowInstance withTerminated(Boolean terminated) {
        this.terminated = terminated;
        return this;
    }
    
    @JsonInclude(Include.NON_ABSENT)
    @JsonSerialize(using = DateTimeSerializer.class)
    @JsonDeserialize(using = DateTimeDeserializer.class)
    @JsonProperty("terminatedAt")
    public OffsetDateTime terminatedAt;

    public WorkflowInstance withTerminatedAt(OffsetDateTime terminatedAt) {
        this.terminatedAt = terminatedAt;
        return this;
    }
    
    @JsonSerialize(using = DateTimeSerializer.class)
    @JsonDeserialize(using = DateTimeDeserializer.class)
    @JsonProperty("updatedAt")
    public OffsetDateTime updatedAt;

    public WorkflowInstance withUpdatedAt(OffsetDateTime updatedAt) {
        this.updatedAt = updatedAt;
        return this;
    }
    
    @JsonProperty("workflowID")
    public String workflowID;

    public WorkflowInstance withWorkflowID(String workflowID) {
        this.workflowID = workflowID;
        return this;
    }
    
    public WorkflowInstance(@JsonProperty("createdAt") OffsetDateTime createdAt, @JsonProperty("id") String id, @JsonProperty("terminated") Boolean terminated, @JsonProperty("updatedAt") OffsetDateTime updatedAt, @JsonProperty("workflowID") String workflowID) {
        this.createdAt = createdAt;
        this.id = id;
        this.terminated = terminated;
        this.updatedAt = updatedAt;
        this.workflowID = workflowID;
  }
}
