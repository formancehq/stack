/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package org.openapis.openapi.models.shared;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * AccountsCursor - OK
 */
public class AccountsCursor {
    @JsonProperty("cursor")
    public AccountsCursorCursor cursor;

    public AccountsCursor withCursor(AccountsCursorCursor cursor) {
        this.cursor = cursor;
        return this;
    }
    
    public AccountsCursor(@JsonProperty("cursor") AccountsCursorCursor cursor) {
        this.cursor = cursor;
  }
}
