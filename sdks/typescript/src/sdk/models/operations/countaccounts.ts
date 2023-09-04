/*
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

import { SpeakeasyBase, SpeakeasyMetadata } from "../../../internal/utils";
import * as shared from "../shared";
import { AxiosResponse } from "axios";

/**
 * Filter accounts by metadata key value pairs. The filter can be used like this -> metadata[key]=value1&metadata[a.nested.key]=value2
 */
export class CountAccountsMetadata extends SpeakeasyBase {}

export class CountAccountsRequest extends SpeakeasyBase {
    /**
     * Filter accounts by address pattern (regular expression placed between ^ and $).
     */
    @SpeakeasyMetadata({ data: "queryParam, style=form;explode=true;name=address" })
    address?: string;

    /**
     * Name of the ledger.
     */
    @SpeakeasyMetadata({ data: "pathParam, style=simple;explode=false;name=ledger" })
    ledger: string;

    /**
     * Filter accounts by metadata key value pairs. The filter can be used like this -> metadata[key]=value1&metadata[a.nested.key]=value2
     */
    @SpeakeasyMetadata({ data: "queryParam, style=deepObject;explode=true;name=metadata" })
    metadata?: CountAccountsMetadata;
}

export class CountAccountsResponse extends SpeakeasyBase {
    @SpeakeasyMetadata()
    contentType: string;

    /**
     * Error
     */
    @SpeakeasyMetadata()
    errorResponse?: shared.ErrorResponse;

    @SpeakeasyMetadata()
    headers?: Record<string, string[]>;

    @SpeakeasyMetadata()
    statusCode: number;

    @SpeakeasyMetadata()
    rawResponse?: AxiosResponse;
}
