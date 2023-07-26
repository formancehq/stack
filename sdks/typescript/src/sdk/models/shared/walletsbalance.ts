/*
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

import { SpeakeasyBase, SpeakeasyMetadata } from "../../../internal/utils";
import { Expose, Transform } from "class-transformer";

export class WalletsBalance extends SpeakeasyBase {
  @SpeakeasyMetadata()
  @Expose({ name: "expiresAt" })
  @Transform(({ value }) => new Date(value), { toClassOnly: true })
  expiresAt?: Date;

  @SpeakeasyMetadata()
  @Expose({ name: "name" })
  name: string;

  @SpeakeasyMetadata()
  @Expose({ name: "priority" })
  priority?: number;
}
