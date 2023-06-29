/*
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

import { SpeakeasyBase, SpeakeasyMetadata } from "../../../internal/utils";
import * as shared from "../shared";
import { AxiosResponse } from "axios";

export class AddMetadataToAccountRequest extends SpeakeasyBase {
  /**
   * Use an idempotency key
   */
  @SpeakeasyMetadata({
    data: "header, style=simple;explode=false;name=Idempotency-Key",
  })
  idempotencyKey?: string;

  /**
   * metadata
   */
  @SpeakeasyMetadata({ data: "request, media_type=application/json" })
  requestBody: Record<string, string>;

  /**
   * Exact address of the account. It must match the following regular expressions pattern:
   *
   * @remarks
   * ```
   * ^\w+(:\w+)*$
   * ```
   *
   */
  @SpeakeasyMetadata({
    data: "pathParam, style=simple;explode=false;name=address",
  })
  address: string;

  /**
   * Set async mode.
   */
  @SpeakeasyMetadata({ data: "queryParam, style=form;explode=true;name=async" })
  async?: boolean;

  /**
   * Set the dry run mode. Dry run mode doesn't add the logs to the database or publish a message to the message broker.
   */
  @SpeakeasyMetadata({
    data: "queryParam, style=form;explode=true;name=dryRun",
  })
  dryRun?: boolean;

  /**
   * Name of the ledger.
   */
  @SpeakeasyMetadata({
    data: "pathParam, style=simple;explode=false;name=ledger",
  })
  ledger: string;
}

export class AddMetadataToAccountResponse extends SpeakeasyBase {
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
}
