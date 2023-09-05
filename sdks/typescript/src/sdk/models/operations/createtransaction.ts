/*
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

import { SpeakeasyBase, SpeakeasyMetadata } from "../../../internal/utils";
import * as shared from "../shared";
import { AxiosResponse } from "axios";

export class CreateTransactionRequest extends SpeakeasyBase {
    /**
     * Use an idempotency key
     */
    @SpeakeasyMetadata({ data: "header, style=simple;explode=false;name=Idempotency-Key" })
    idempotencyKey?: string;

    /**
     * The request body must contain at least one of the following objects:
     *
     * @remarks
     *   - `postings`: suitable for simple transactions
     *   - `script`: enabling more complex transactions with Numscript
     *
     */
    @SpeakeasyMetadata({ data: "request, media_type=application/json" })
    postTransaction: shared.PostTransaction;

    /**
     * Name of the ledger.
     */
    @SpeakeasyMetadata({ data: "pathParam, style=simple;explode=false;name=ledger" })
    ledger: string;

    /**
     * Set the preview mode. Preview mode doesn't add the logs to the database or publish a message to the message broker.
     */
    @SpeakeasyMetadata({ data: "queryParam, style=form;explode=true;name=preview" })
    preview?: boolean;
}

export class CreateTransactionResponse extends SpeakeasyBase {
    @SpeakeasyMetadata()
    contentType: string;

    /**
     * Error
     */
    @SpeakeasyMetadata()
    errorResponse?: shared.ErrorResponse;

    @SpeakeasyMetadata()
    statusCode: number;

    @SpeakeasyMetadata()
    rawResponse?: AxiosResponse;

    /**
     * OK
     */
    @SpeakeasyMetadata()
    transactionsResponse?: shared.TransactionsResponse;
}
