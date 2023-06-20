/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.shared;

import com.fasterxml.jackson.annotation.JsonInclude.Include;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;

public class TaskWiseDescriptor {
    @JsonInclude(Include.NON_ABSENT)
    @JsonProperty("key")
    public String key;

    public TaskWiseDescriptor withKey(String key) {
        this.key = key;
        return this;
    }
    
    @JsonInclude(Include.NON_ABSENT)
    @JsonProperty("name")
    public String name;

    public TaskWiseDescriptor withName(String name) {
        this.name = name;
        return this;
    }
    
    @JsonInclude(Include.NON_ABSENT)
    @JsonProperty("profileID")
    public Long profileID;

    public TaskWiseDescriptor withProfileID(Long profileID) {
        this.profileID = profileID;
        return this;
    }
    
    public TaskWiseDescriptor(){}
}
