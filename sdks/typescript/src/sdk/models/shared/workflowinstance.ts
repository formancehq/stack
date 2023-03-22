/*
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

import { SpeakeasyBase, SpeakeasyMetadata } from "../../../internal/utils";
import { StageStatus } from "./stagestatus";
import { Expose, Transform, Type } from "class-transformer";

export class WorkflowInstance extends SpeakeasyBase {
  @SpeakeasyMetadata()
  @Expose({ name: "createdAt" })
  @Transform(({ value }) => new Date(value), { toClassOnly: true })
  createdAt: Date;

  @SpeakeasyMetadata()
  @Expose({ name: "error" })
  error?: string;

  @SpeakeasyMetadata()
  @Expose({ name: "id" })
  id: string;

  @SpeakeasyMetadata({ elemType: StageStatus })
  @Expose({ name: "status" })
  @Type(() => StageStatus)
  status?: StageStatus[];

  @SpeakeasyMetadata()
  @Expose({ name: "terminated" })
  terminated: boolean;

  @SpeakeasyMetadata()
  @Expose({ name: "terminatedAt" })
  @Transform(({ value }) => new Date(value), { toClassOnly: true })
  terminatedAt?: Date;

  @SpeakeasyMetadata()
  @Expose({ name: "updatedAt" })
  @Transform(({ value }) => new Date(value), { toClassOnly: true })
  updatedAt: Date;

  @SpeakeasyMetadata()
  @Expose({ name: "workflowID" })
  workflowID: string;
}
