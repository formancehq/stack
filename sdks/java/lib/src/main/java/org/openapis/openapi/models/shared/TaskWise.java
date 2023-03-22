/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package org.openapis.openapi.models.shared;

import com.fasterxml.jackson.annotation.JsonInclude.Include;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.databind.annotation.JsonDeserialize;
import com.fasterxml.jackson.databind.annotation.JsonSerialize;
import java.time.OffsetDateTime;
import org.openapis.openapi.utils.DateTimeDeserializer;
import org.openapis.openapi.utils.DateTimeSerializer;

public class TaskWise {
    @JsonProperty("connectorId")
    public String connectorId;

    public TaskWise withConnectorId(String connectorId) {
        this.connectorId = connectorId;
        return this;
    }
    
    @JsonSerialize(using = DateTimeSerializer.class)
    @JsonDeserialize(using = DateTimeDeserializer.class)
    @JsonProperty("createdAt")
    public OffsetDateTime createdAt;

    public TaskWise withCreatedAt(OffsetDateTime createdAt) {
        this.createdAt = createdAt;
        return this;
    }
    
    @JsonProperty("descriptor")
    public TaskWiseDescriptor descriptor;

    public TaskWise withDescriptor(TaskWiseDescriptor descriptor) {
        this.descriptor = descriptor;
        return this;
    }
    
    @JsonInclude(Include.NON_ABSENT)
    @JsonProperty("error")
    public String error;

    public TaskWise withError(String error) {
        this.error = error;
        return this;
    }
    
    @JsonProperty("id")
    public String id;

    public TaskWise withId(String id) {
        this.id = id;
        return this;
    }
    
    @JsonProperty("state")
    public java.util.Map<String, Object> state;

    public TaskWise withState(java.util.Map<String, Object> state) {
        this.state = state;
        return this;
    }
    
    @JsonProperty("status")
    public PaymentStatus status;

    public TaskWise withStatus(PaymentStatus status) {
        this.status = status;
        return this;
    }
    
    @JsonSerialize(using = DateTimeSerializer.class)
    @JsonDeserialize(using = DateTimeDeserializer.class)
    @JsonProperty("updatedAt")
    public OffsetDateTime updatedAt;

    public TaskWise withUpdatedAt(OffsetDateTime updatedAt) {
        this.updatedAt = updatedAt;
        return this;
    }
    
    public TaskWise(@JsonProperty("connectorId") String connectorId, @JsonProperty("createdAt") OffsetDateTime createdAt, @JsonProperty("descriptor") TaskWiseDescriptor descriptor, @JsonProperty("id") String id, @JsonProperty("state") java.util.Map<String, Object> state, @JsonProperty("status") PaymentStatus status, @JsonProperty("updatedAt") OffsetDateTime updatedAt) {
        this.connectorId = connectorId;
        this.createdAt = createdAt;
        this.descriptor = descriptor;
        this.id = id;
        this.state = state;
        this.status = status;
        this.updatedAt = updatedAt;
  }
}
