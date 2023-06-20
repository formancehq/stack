/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package org.openapis.openapi.models.shared;

import com.fasterxml.jackson.annotation.JsonInclude.Include;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * CreateScopeResponse - Created scope
 */
public class CreateScopeResponse {
    @JsonInclude(Include.NON_ABSENT)
    @JsonProperty("data")
    public Scope data;

    public CreateScopeResponse withData(Scope data) {
        this.data = data;
        return this;
    }
    
    public CreateScopeResponse(){}
}
