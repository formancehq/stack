/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package org.openapis.openapi.models.shared;

import com.fasterxml.jackson.annotation.JsonValue;

public enum ErrorErrorCode {
    VALIDATION("VALIDATION"),
    NOT_FOUND("NOT_FOUND");

    @JsonValue
    public final String value;

    private ErrorErrorCode(String value) {
        this.value = value;
    }
}
