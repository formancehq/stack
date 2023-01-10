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

import { PaymentStatus } from '../models/PaymentStatus';
import { TaskStripeDescriptor } from '../models/TaskStripeDescriptor';
import { HttpFile } from '../http/http';

export class TaskStripe {
    'id'?: string;
    'connectorID'?: string;
    'createdAt'?: Date;
    'updatedAt'?: Date;
    'descriptor'?: TaskStripeDescriptor;
    'status'?: PaymentStatus;
    'state'?: any;
    'error'?: string;

    static readonly discriminator: string | undefined = undefined;

    static readonly attributeTypeMap: Array<{name: string, baseName: string, type: string, format: string}> = [
        {
            "name": "id",
            "baseName": "id",
            "type": "string",
            "format": "uuid"
        },
        {
            "name": "connectorID",
            "baseName": "connectorID",
            "type": "string",
            "format": "uuid"
        },
        {
            "name": "createdAt",
            "baseName": "createdAt",
            "type": "Date",
            "format": "date-time"
        },
        {
            "name": "updatedAt",
            "baseName": "updatedAt",
            "type": "Date",
            "format": "date-time"
        },
        {
            "name": "descriptor",
            "baseName": "descriptor",
            "type": "TaskStripeDescriptor",
            "format": ""
        },
        {
            "name": "status",
            "baseName": "status",
            "type": "PaymentStatus",
            "format": ""
        },
        {
            "name": "state",
            "baseName": "state",
            "type": "any",
            "format": ""
        },
        {
            "name": "error",
            "baseName": "error",
            "type": "string",
            "format": ""
        }    ];

    static getAttributeTypeMap() {
        return TaskStripe.attributeTypeMap;
    }

    public constructor() {
    }
}

