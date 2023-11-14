/*
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

import { SpeakeasyBase, SpeakeasyMetadata } from "../../../internal/utils";
import { Expose } from "class-transformer";

export class ActivityStripeTransfer extends SpeakeasyBase {
  @SpeakeasyMetadata()
  @Expose({ name: "amount" })
  amount?: number;

  @SpeakeasyMetadata()
  @Expose({ name: "asset" })
  asset?: string;

  @SpeakeasyMetadata()
  @Expose({ name: "connectorID" })
  connectorID?: string;

  @SpeakeasyMetadata()
  @Expose({ name: "destination" })
  destination?: string;

  /**
   * A set of key/value pairs that you can attach to a transfer object.
   *
   * @remarks
   * It can be useful for storing additional information about the transfer in a structured format.
   *
   */
  @SpeakeasyMetadata()
  @Expose({ name: "metadata" })
  metadata?: Record<string, any>;
}
