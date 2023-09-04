/*
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

import { SpeakeasyBase, SpeakeasyMetadata } from "../../../internal/utils";
import { WorkflowInstanceHistoryStage } from "./workflowinstancehistorystage";
import { Expose, Type } from "class-transformer";

/**
 * The workflow instance stage history
 */
export class GetWorkflowInstanceHistoryStageResponse extends SpeakeasyBase {
    @SpeakeasyMetadata({ elemType: WorkflowInstanceHistoryStage })
    @Expose({ name: "data" })
    @Type(() => WorkflowInstanceHistoryStage)
    data: WorkflowInstanceHistoryStage[];
}
