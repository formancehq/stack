/*
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

import { SpeakeasyBase, SpeakeasyMetadata } from "../../../internal/utils";
import * as shared from "../shared";
import { AxiosResponse } from "axios";

export class GetTransferInitiationRequest extends SpeakeasyBase {
  /**
   * The transfer ID.
   */
  @SpeakeasyMetadata({
    data: "pathParam, style=simple;explode=false;name=transferId",
  })
  transferId: string;
}

export class GetTransferInitiationResponse extends SpeakeasyBase {
  @SpeakeasyMetadata()
  contentType: string;

  @SpeakeasyMetadata()
  statusCode: number;

  @SpeakeasyMetadata()
  rawResponse?: AxiosResponse;

  /**
   * OK
   */
  @SpeakeasyMetadata()
  transferInitiationResponse?: shared.TransferInitiationResponse;
}
