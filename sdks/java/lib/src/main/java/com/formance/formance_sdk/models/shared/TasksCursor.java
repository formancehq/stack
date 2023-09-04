/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package com.formance.formance_sdk.models.shared;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * TasksCursor - OK
 */

public class TasksCursor {
    @JsonProperty("cursor")
    public TasksCursorCursor cursor;

    public TasksCursor withCursor(TasksCursorCursor cursor) {
        this.cursor = cursor;
        return this;
    }
    
    public TasksCursor(@JsonProperty("cursor") TasksCursorCursor cursor) {
        this.cursor = cursor;
  }
}
