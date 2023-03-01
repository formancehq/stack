/**
 * Formance Stack API
 * Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions /> 
 *
 * OpenAPI spec version: v1.0.20230301
 * Contact: support@formance.com
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { HttpFile } from '../http/http';

export class StripeConfig {
    /**
    * The frequency at which the connector will try to fetch new BalanceTransaction objects from Stripe API. 
    */
    'pollingPeriod'?: string;
    'apiKey': string;
    /**
    * Number of BalanceTransaction to fetch at each polling interval. 
    */
    'pageSize'?: number;

    static readonly discriminator: string | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "pollingPeriod",
            "baseName": "pollingPeriod",
            "type": "string",
            "format": ""
        },
        {
            "name": "apiKey",
            "baseName": "apiKey",
            "type": "string",
            "format": ""
        },
        {
            "name": "pageSize",
            "baseName": "pageSize",
            "type": "number",
            "format": "int64"
        }    ];

    static getAttributeTypeMap() {
        return StripeConfig.attributeTypeMap;
    }

    public constructor() {
    }
}

