/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.operations;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.formance.formance_sdk.utils.SpeakeasyMetadata;

public class RunWorkflowRequest {
    @SpeakeasyMetadata("request:mediaType=application/json")
    public java.util.Map<String, String> requestBody;

    public RunWorkflowRequest withRequestBody(java.util.Map<String, String> requestBody) {
        this.requestBody = requestBody;
        return this;
    }
    
    /**
     * Wait end of the workflow before return
     */
    @SpeakeasyMetadata("queryParam:style=form,explode=true,name=wait")
    public Boolean wait;

    public RunWorkflowRequest withWait(Boolean wait) {
        this.wait = wait;
        return this;
    }
    
    /**
     * The flow id
     */
    @SpeakeasyMetadata("pathParam:style=simple,explode=false,name=workflowID")
    public String workflowID;

    public RunWorkflowRequest withWorkflowID(String workflowID) {
        this.workflowID = workflowID;
        return this;
    }
    
    public RunWorkflowRequest(@JsonProperty("workflowID") String workflowID) {
        this.workflowID = workflowID;
  }
}
