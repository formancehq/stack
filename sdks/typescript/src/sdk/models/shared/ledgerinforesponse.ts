/*
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

import { SpeakeasyBase, SpeakeasyMetadata } from "../../../internal/utils";
import { LedgerInfo } from "./ledgerinfo";
import { Expose, Type } from "class-transformer";

/**
 * OK
 */
export class LedgerInfoResponse extends SpeakeasyBase {
    @SpeakeasyMetadata()
    @Expose({ name: "data" })
    @Type(() => LedgerInfo)
    data?: LedgerInfo;
}
