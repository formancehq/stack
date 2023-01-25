/**
 * Formance Stack API
 * Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions /> 
 *
 * OpenAPI spec version: develop
 * Contact: support@formance.com
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { ActivityConfirmHold } from '../models/ActivityConfirmHold';
import { ActivityCreateTransaction } from '../models/ActivityCreateTransaction';
import { ActivityCreditWallet } from '../models/ActivityCreditWallet';
import { ActivityDebitWallet } from '../models/ActivityDebitWallet';
import { ActivityGetAccount } from '../models/ActivityGetAccount';
import { ActivityGetPayment } from '../models/ActivityGetPayment';
import { ActivityGetWallet } from '../models/ActivityGetWallet';
import { ActivityRevertTransaction } from '../models/ActivityRevertTransaction';
import { ActivityVoidHold } from '../models/ActivityVoidHold';
import { DebitWalletRequest } from '../models/DebitWalletRequest';
import { StripeTransferRequest } from '../models/StripeTransferRequest';
import { HttpFile } from '../http/http';

export class WorkflowInstanceHistoryStageInput {
    'id': string;
    'ledger': string;
    'data'?: DebitWalletRequest;
    'amount'?: number;
    'asset'?: string;
    'destination'?: string;
    /**
    * A set of key/value pairs that you can attach to a transfer object. It can be useful for storing additional information about the transfer in a structured format. 
    */
    'metadata'?: any;

    static readonly discriminator: string | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "id",
            "baseName": "id",
            "type": "string",
            "format": ""
        },
        {
            "name": "ledger",
            "baseName": "ledger",
            "type": "string",
            "format": ""
        },
        {
            "name": "data",
            "baseName": "data",
            "type": "DebitWalletRequest",
            "format": ""
        },
        {
            "name": "amount",
            "baseName": "amount",
            "type": "number",
            "format": "int64"
        },
        {
            "name": "asset",
            "baseName": "asset",
            "type": "string",
            "format": ""
        },
        {
            "name": "destination",
            "baseName": "destination",
            "type": "string",
            "format": ""
        },
        {
            "name": "metadata",
            "baseName": "metadata",
            "type": "any",
            "format": ""
        }    ];

    static getAttributeTypeMap() {
        return WorkflowInstanceHistoryStageInput.attributeTypeMap;
    }

    public constructor() {
    }
}

