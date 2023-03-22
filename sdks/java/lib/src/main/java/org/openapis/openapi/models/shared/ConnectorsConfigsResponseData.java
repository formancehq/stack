/* 
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

package org.openapis.openapi.models.shared;

import com.fasterxml.jackson.annotation.JsonProperty;

public class ConnectorsConfigsResponseData {
    @JsonProperty("connector")
    public ConnectorsConfigsResponseDataConnector connector;

    public ConnectorsConfigsResponseData withConnector(ConnectorsConfigsResponseDataConnector connector) {
        this.connector = connector;
        return this;
    }
    
    public ConnectorsConfigsResponseData(@JsonProperty("connector") ConnectorsConfigsResponseDataConnector connector) {
        this.connector = connector;
  }
}
