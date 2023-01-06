export * from '../models/Account';
export * from '../models/AccountWithVolumesAndBalances';
export * from '../models/AddMetadataToAccount409Response';
export * from '../models/AssetHolder';
export * from '../models/Attempt';
export * from '../models/AttemptResponse';
export * from '../models/Balance';
export * from '../models/BalanceWithAssets';
export * from '../models/BankingCircleConfig';
export * from '../models/Client';
export * from '../models/ClientAllOf';
export * from '../models/ClientOptions';
export * from '../models/ClientSecret';
export * from '../models/Config';
export * from '../models/ConfigChangeSecret';
export * from '../models/ConfigInfo';
export * from '../models/ConfigInfoResponse';
export * from '../models/ConfigResponse';
export * from '../models/ConfigUser';
export * from '../models/ConfigsResponse';
export * from '../models/ConfirmHoldRequest';
export * from '../models/ConnectorBaseInfo';
export * from '../models/ConnectorConfig';
export * from '../models/Connectors';
export * from '../models/Contract';
export * from '../models/CreateBalanceResponse';
export * from '../models/CreateClientResponse';
export * from '../models/CreateScopeResponse';
export * from '../models/CreateSecretResponse';
export * from '../models/CreateTransaction400Response';
export * from '../models/CreateTransaction409Response';
export * from '../models/CreateTransactions400Response';
export * from '../models/CreateWalletRequest';
export * from '../models/CreateWalletResponse';
export * from '../models/CreditWalletRequest';
export * from '../models/CurrencyCloudConfig';
export * from '../models/Cursor';
export * from '../models/DebitWalletRequest';
export * from '../models/DebitWalletResponse';
export * from '../models/DummyPayConfig';
export * from '../models/ErrorCode';
export * from '../models/ErrorResponse';
export * from '../models/ExpandedDebitHold';
export * from '../models/ExpandedDebitHoldAllOf';
export * from '../models/GetAccount200Response';
export * from '../models/GetAccount400Response';
export * from '../models/GetBalanceResponse';
export * from '../models/GetBalances200Response';
export * from '../models/GetBalances200ResponseCursor';
export * from '../models/GetBalances200ResponseCursorAllOf';
export * from '../models/GetBalancesAggregated200Response';
export * from '../models/GetBalancesAggregated400Response';
export * from '../models/GetHoldResponse';
export * from '../models/GetHoldsResponse';
export * from '../models/GetHoldsResponseCursor';
export * from '../models/GetHoldsResponseCursorAllOf';
export * from '../models/GetPaymentResponse';
export * from '../models/GetTransaction400Response';
export * from '../models/GetTransaction404Response';
export * from '../models/GetTransactionsResponse';
export * from '../models/GetTransactionsResponseCursor';
export * from '../models/GetTransactionsResponseCursorAllOf';
export * from '../models/GetWalletResponse';
export * from '../models/Hold';
export * from '../models/LedgerAccountSubject';
export * from '../models/LedgerStorage';
export * from '../models/ListAccounts200Response';
export * from '../models/ListAccounts200ResponseCursor';
export * from '../models/ListAccounts200ResponseCursorAllOf';
export * from '../models/ListAccounts400Response';
export * from '../models/ListBalancesResponse';
export * from '../models/ListBalancesResponseCursor';
export * from '../models/ListBalancesResponseCursorAllOf';
export * from '../models/ListClientsResponse';
export * from '../models/ListConnectorTasks200ResponseInner';
export * from '../models/ListConnectorsConfigsResponse';
export * from '../models/ListConnectorsConfigsResponseConnector';
export * from '../models/ListConnectorsConfigsResponseConnectorKey';
export * from '../models/ListConnectorsResponse';
export * from '../models/ListPaymentsResponse';
export * from '../models/ListScopesResponse';
export * from '../models/ListTransactions200Response';
export * from '../models/ListTransactions200ResponseCursor';
export * from '../models/ListTransactions200ResponseCursorAllOf';
export * from '../models/ListUsersResponse';
export * from '../models/ListWalletsResponse';
export * from '../models/ListWalletsResponseCursor';
export * from '../models/ListWalletsResponseCursorAllOf';
export * from '../models/Mapping';
export * from '../models/MappingResponse';
export * from '../models/ModulrConfig';
export * from '../models/Monetary';
export * from '../models/Payment';
export * from '../models/Posting';
export * from '../models/Query';
export * from '../models/ReadClientResponse';
export * from '../models/ReadUserResponse';
export * from '../models/Response';
export * from '../models/RunScript400Response';
export * from '../models/Scope';
export * from '../models/ScopeAllOf';
export * from '../models/ScopeOptions';
export * from '../models/Script';
export * from '../models/ScriptResult';
export * from '../models/Secret';
export * from '../models/SecretAllOf';
export * from '../models/SecretOptions';
export * from '../models/ServerInfo';
export * from '../models/Stats';
export * from '../models/StatsResponse';
export * from '../models/StripeConfig';
export * from '../models/StripeTask';
export * from '../models/StripeTransferRequest';
export * from '../models/Subject';
export * from '../models/TaskDescriptorBankingCircle';
export * from '../models/TaskDescriptorBankingCircleDescriptor';
export * from '../models/TaskDescriptorCurrencyCloud';
export * from '../models/TaskDescriptorCurrencyCloudDescriptor';
export * from '../models/TaskDescriptorDummyPay';
export * from '../models/TaskDescriptorDummyPayDescriptor';
export * from '../models/TaskDescriptorModulr';
export * from '../models/TaskDescriptorModulrDescriptor';
export * from '../models/TaskDescriptorStripe';
export * from '../models/TaskDescriptorStripeDescriptor';
export * from '../models/TaskDescriptorWise';
export * from '../models/TaskDescriptorWiseDescriptor';
export * from '../models/Total';
export * from '../models/Transaction';
export * from '../models/TransactionData';
export * from '../models/TransactionResponse';
export * from '../models/Transactions';
export * from '../models/TransactionsResponse';
export * from '../models/UpdateWalletRequest';
export * from '../models/User';
export * from '../models/Volume';
export * from '../models/Wallet';
export * from '../models/WalletSubject';
export * from '../models/WalletWithBalances';
export * from '../models/WalletWithBalancesBalances';
export * from '../models/WalletsCursor';
export * from '../models/WalletsErrorResponse';
export * from '../models/WalletsPosting';
export * from '../models/WalletsTransaction';
export * from '../models/WalletsVolume';
export * from '../models/WebhooksConfig';
export * from '../models/WebhooksCursor';
export * from '../models/WiseConfig';

import { Account } from '../models/Account';
import { AccountWithVolumesAndBalances } from '../models/AccountWithVolumesAndBalances';
import { AddMetadataToAccount409Response } from '../models/AddMetadataToAccount409Response';
import { AssetHolder } from '../models/AssetHolder';
import { Attempt } from '../models/Attempt';
import { AttemptResponse } from '../models/AttemptResponse';
import { Balance } from '../models/Balance';
import { BalanceWithAssets } from '../models/BalanceWithAssets';
import { BankingCircleConfig } from '../models/BankingCircleConfig';
import { Client } from '../models/Client';
import { ClientAllOf } from '../models/ClientAllOf';
import { ClientOptions } from '../models/ClientOptions';
import { ClientSecret } from '../models/ClientSecret';
import { Config } from '../models/Config';
import { ConfigChangeSecret } from '../models/ConfigChangeSecret';
import { ConfigInfo } from '../models/ConfigInfo';
import { ConfigInfoResponse } from '../models/ConfigInfoResponse';
import { ConfigResponse } from '../models/ConfigResponse';
import { ConfigUser } from '../models/ConfigUser';
import { ConfigsResponse } from '../models/ConfigsResponse';
import { ConfirmHoldRequest } from '../models/ConfirmHoldRequest';
import { ConnectorBaseInfo } from '../models/ConnectorBaseInfo';
import { ConnectorConfig } from '../models/ConnectorConfig';
import { Connectors } from '../models/Connectors';
import { Contract } from '../models/Contract';
import { CreateBalanceResponse } from '../models/CreateBalanceResponse';
import { CreateClientResponse } from '../models/CreateClientResponse';
import { CreateScopeResponse } from '../models/CreateScopeResponse';
import { CreateSecretResponse } from '../models/CreateSecretResponse';
import { CreateTransaction400Response } from '../models/CreateTransaction400Response';
import { CreateTransaction409Response } from '../models/CreateTransaction409Response';
import { CreateTransactions400Response } from '../models/CreateTransactions400Response';
import { CreateWalletRequest } from '../models/CreateWalletRequest';
import { CreateWalletResponse } from '../models/CreateWalletResponse';
import { CreditWalletRequest } from '../models/CreditWalletRequest';
import { CurrencyCloudConfig } from '../models/CurrencyCloudConfig';
import { Cursor } from '../models/Cursor';
import { DebitWalletRequest } from '../models/DebitWalletRequest';
import { DebitWalletResponse } from '../models/DebitWalletResponse';
import { DummyPayConfig } from '../models/DummyPayConfig';
import { ErrorCode } from '../models/ErrorCode';
import { ErrorResponse } from '../models/ErrorResponse';
import { ExpandedDebitHold } from '../models/ExpandedDebitHold';
import { ExpandedDebitHoldAllOf } from '../models/ExpandedDebitHoldAllOf';
import { GetAccount200Response } from '../models/GetAccount200Response';
import { GetAccount400Response } from '../models/GetAccount400Response';
import { GetBalanceResponse } from '../models/GetBalanceResponse';
import { GetBalances200Response } from '../models/GetBalances200Response';
import { GetBalances200ResponseCursor } from '../models/GetBalances200ResponseCursor';
import { GetBalances200ResponseCursorAllOf } from '../models/GetBalances200ResponseCursorAllOf';
import { GetBalancesAggregated200Response } from '../models/GetBalancesAggregated200Response';
import { GetBalancesAggregated400Response } from '../models/GetBalancesAggregated400Response';
import { GetHoldResponse } from '../models/GetHoldResponse';
import { GetHoldsResponse } from '../models/GetHoldsResponse';
import { GetHoldsResponseCursor } from '../models/GetHoldsResponseCursor';
import { GetHoldsResponseCursorAllOf } from '../models/GetHoldsResponseCursorAllOf';
import { GetPaymentResponse } from '../models/GetPaymentResponse';
import { GetTransaction400Response } from '../models/GetTransaction400Response';
import { GetTransaction404Response } from '../models/GetTransaction404Response';
import { GetTransactionsResponse } from '../models/GetTransactionsResponse';
import { GetTransactionsResponseCursor } from '../models/GetTransactionsResponseCursor';
import { GetTransactionsResponseCursorAllOf } from '../models/GetTransactionsResponseCursorAllOf';
import { GetWalletResponse } from '../models/GetWalletResponse';
import { Hold } from '../models/Hold';
import { LedgerAccountSubject } from '../models/LedgerAccountSubject';
import { LedgerStorage } from '../models/LedgerStorage';
import { ListAccounts200Response } from '../models/ListAccounts200Response';
import { ListAccounts200ResponseCursor } from '../models/ListAccounts200ResponseCursor';
import { ListAccounts200ResponseCursorAllOf } from '../models/ListAccounts200ResponseCursorAllOf';
import { ListAccounts400Response } from '../models/ListAccounts400Response';
import { ListBalancesResponse } from '../models/ListBalancesResponse';
import { ListBalancesResponseCursor } from '../models/ListBalancesResponseCursor';
import { ListBalancesResponseCursorAllOf } from '../models/ListBalancesResponseCursorAllOf';
import { ListClientsResponse } from '../models/ListClientsResponse';
import { ListConnectorTasks200ResponseInner  , ListConnectorTasks200ResponseInnerStatusEnum      } from '../models/ListConnectorTasks200ResponseInner';
import { ListConnectorsConfigsResponse } from '../models/ListConnectorsConfigsResponse';
import { ListConnectorsConfigsResponseConnector } from '../models/ListConnectorsConfigsResponseConnector';
import { ListConnectorsConfigsResponseConnectorKey } from '../models/ListConnectorsConfigsResponseConnectorKey';
import { ListConnectorsResponse } from '../models/ListConnectorsResponse';
import { ListPaymentsResponse } from '../models/ListPaymentsResponse';
import { ListScopesResponse } from '../models/ListScopesResponse';
import { ListTransactions200Response } from '../models/ListTransactions200Response';
import { ListTransactions200ResponseCursor } from '../models/ListTransactions200ResponseCursor';
import { ListTransactions200ResponseCursorAllOf } from '../models/ListTransactions200ResponseCursorAllOf';
import { ListUsersResponse } from '../models/ListUsersResponse';
import { ListWalletsResponse } from '../models/ListWalletsResponse';
import { ListWalletsResponseCursor } from '../models/ListWalletsResponseCursor';
import { ListWalletsResponseCursorAllOf } from '../models/ListWalletsResponseCursorAllOf';
import { Mapping } from '../models/Mapping';
import { MappingResponse } from '../models/MappingResponse';
import { ModulrConfig } from '../models/ModulrConfig';
import { Monetary } from '../models/Monetary';
import { Payment  , PaymentSchemeEnum   , PaymentTypeEnum        } from '../models/Payment';
import { Posting } from '../models/Posting';
import { Query } from '../models/Query';
import { ReadClientResponse } from '../models/ReadClientResponse';
import { ReadUserResponse } from '../models/ReadUserResponse';
import { Response } from '../models/Response';
import { RunScript400Response } from '../models/RunScript400Response';
import { Scope } from '../models/Scope';
import { ScopeAllOf } from '../models/ScopeAllOf';
import { ScopeOptions } from '../models/ScopeOptions';
import { Script } from '../models/Script';
import { ScriptResult , ScriptResultErrorCodeEnum     } from '../models/ScriptResult';
import { Secret } from '../models/Secret';
import { SecretAllOf } from '../models/SecretAllOf';
import { SecretOptions } from '../models/SecretOptions';
import { ServerInfo } from '../models/ServerInfo';
import { Stats } from '../models/Stats';
import { StatsResponse } from '../models/StatsResponse';
import { StripeConfig } from '../models/StripeConfig';
import { StripeTask } from '../models/StripeTask';
import { StripeTransferRequest } from '../models/StripeTransferRequest';
import { Subject } from '../models/Subject';
import { TaskDescriptorBankingCircle  , TaskDescriptorBankingCircleStatusEnum      } from '../models/TaskDescriptorBankingCircle';
import { TaskDescriptorBankingCircleDescriptor } from '../models/TaskDescriptorBankingCircleDescriptor';
import { TaskDescriptorCurrencyCloud  , TaskDescriptorCurrencyCloudStatusEnum      } from '../models/TaskDescriptorCurrencyCloud';
import { TaskDescriptorCurrencyCloudDescriptor } from '../models/TaskDescriptorCurrencyCloudDescriptor';
import { TaskDescriptorDummyPay  , TaskDescriptorDummyPayStatusEnum      } from '../models/TaskDescriptorDummyPay';
import { TaskDescriptorDummyPayDescriptor } from '../models/TaskDescriptorDummyPayDescriptor';
import { TaskDescriptorModulr  , TaskDescriptorModulrStatusEnum      } from '../models/TaskDescriptorModulr';
import { TaskDescriptorModulrDescriptor } from '../models/TaskDescriptorModulrDescriptor';
import { TaskDescriptorStripe  , TaskDescriptorStripeStatusEnum      } from '../models/TaskDescriptorStripe';
import { TaskDescriptorStripeDescriptor } from '../models/TaskDescriptorStripeDescriptor';
import { TaskDescriptorWise  , TaskDescriptorWiseStatusEnum      } from '../models/TaskDescriptorWise';
import { TaskDescriptorWiseDescriptor } from '../models/TaskDescriptorWiseDescriptor';
import { Total } from '../models/Total';
import { Transaction } from '../models/Transaction';
import { TransactionData } from '../models/TransactionData';
import { TransactionResponse } from '../models/TransactionResponse';
import { Transactions } from '../models/Transactions';
import { TransactionsResponse } from '../models/TransactionsResponse';
import { UpdateWalletRequest } from '../models/UpdateWalletRequest';
import { User } from '../models/User';
import { Volume } from '../models/Volume';
import { Wallet } from '../models/Wallet';
import { WalletSubject } from '../models/WalletSubject';
import { WalletWithBalances } from '../models/WalletWithBalances';
import { WalletWithBalancesBalances } from '../models/WalletWithBalancesBalances';
import { WalletsCursor } from '../models/WalletsCursor';
import { WalletsErrorResponse, WalletsErrorResponseErrorCodeEnum    } from '../models/WalletsErrorResponse';
import { WalletsPosting } from '../models/WalletsPosting';
import { WalletsTransaction } from '../models/WalletsTransaction';
import { WalletsVolume } from '../models/WalletsVolume';
import { WebhooksConfig } from '../models/WebhooksConfig';
import { WebhooksCursor } from '../models/WebhooksCursor';
import { WiseConfig } from '../models/WiseConfig';

/* tslint:disable:no-unused-variable */
let primitives = [
                    "string",
                    "boolean",
                    "double",
                    "integer",
                    "long",
                    "float",
                    "number",
                    "any"
                 ];

const supportedMediaTypes: { [mediaType: string]: number } = {
  "application/json": Infinity,
  "application/octet-stream": 0,
  "application/x-www-form-urlencoded": 0
}


let enumsMap: Set<string> = new Set<string>([
    "Connectors",
    "ErrorCode",
    "ListConnectorTasks200ResponseInnerStatusEnum",
    "PaymentSchemeEnum",
    "PaymentTypeEnum",
    "ScriptResultErrorCodeEnum",
    "TaskDescriptorBankingCircleStatusEnum",
    "TaskDescriptorCurrencyCloudStatusEnum",
    "TaskDescriptorDummyPayStatusEnum",
    "TaskDescriptorModulrStatusEnum",
    "TaskDescriptorStripeStatusEnum",
    "TaskDescriptorWiseStatusEnum",
    "WalletsErrorResponseErrorCodeEnum",
]);

let typeMap: {[index: string]: any} = {
    "Account": Account,
    "AccountWithVolumesAndBalances": AccountWithVolumesAndBalances,
    "AddMetadataToAccount409Response": AddMetadataToAccount409Response,
    "AssetHolder": AssetHolder,
    "Attempt": Attempt,
    "AttemptResponse": AttemptResponse,
    "Balance": Balance,
    "BalanceWithAssets": BalanceWithAssets,
    "BankingCircleConfig": BankingCircleConfig,
    "Client": Client,
    "ClientAllOf": ClientAllOf,
    "ClientOptions": ClientOptions,
    "ClientSecret": ClientSecret,
    "Config": Config,
    "ConfigChangeSecret": ConfigChangeSecret,
    "ConfigInfo": ConfigInfo,
    "ConfigInfoResponse": ConfigInfoResponse,
    "ConfigResponse": ConfigResponse,
    "ConfigUser": ConfigUser,
    "ConfigsResponse": ConfigsResponse,
    "ConfirmHoldRequest": ConfirmHoldRequest,
    "ConnectorBaseInfo": ConnectorBaseInfo,
    "ConnectorConfig": ConnectorConfig,
    "Contract": Contract,
    "CreateBalanceResponse": CreateBalanceResponse,
    "CreateClientResponse": CreateClientResponse,
    "CreateScopeResponse": CreateScopeResponse,
    "CreateSecretResponse": CreateSecretResponse,
    "CreateTransaction400Response": CreateTransaction400Response,
    "CreateTransaction409Response": CreateTransaction409Response,
    "CreateTransactions400Response": CreateTransactions400Response,
    "CreateWalletRequest": CreateWalletRequest,
    "CreateWalletResponse": CreateWalletResponse,
    "CreditWalletRequest": CreditWalletRequest,
    "CurrencyCloudConfig": CurrencyCloudConfig,
    "Cursor": Cursor,
    "DebitWalletRequest": DebitWalletRequest,
    "DebitWalletResponse": DebitWalletResponse,
    "DummyPayConfig": DummyPayConfig,
    "ErrorResponse": ErrorResponse,
    "ExpandedDebitHold": ExpandedDebitHold,
    "ExpandedDebitHoldAllOf": ExpandedDebitHoldAllOf,
    "GetAccount200Response": GetAccount200Response,
    "GetAccount400Response": GetAccount400Response,
    "GetBalanceResponse": GetBalanceResponse,
    "GetBalances200Response": GetBalances200Response,
    "GetBalances200ResponseCursor": GetBalances200ResponseCursor,
    "GetBalances200ResponseCursorAllOf": GetBalances200ResponseCursorAllOf,
    "GetBalancesAggregated200Response": GetBalancesAggregated200Response,
    "GetBalancesAggregated400Response": GetBalancesAggregated400Response,
    "GetHoldResponse": GetHoldResponse,
    "GetHoldsResponse": GetHoldsResponse,
    "GetHoldsResponseCursor": GetHoldsResponseCursor,
    "GetHoldsResponseCursorAllOf": GetHoldsResponseCursorAllOf,
    "GetPaymentResponse": GetPaymentResponse,
    "GetTransaction400Response": GetTransaction400Response,
    "GetTransaction404Response": GetTransaction404Response,
    "GetTransactionsResponse": GetTransactionsResponse,
    "GetTransactionsResponseCursor": GetTransactionsResponseCursor,
    "GetTransactionsResponseCursorAllOf": GetTransactionsResponseCursorAllOf,
    "GetWalletResponse": GetWalletResponse,
    "Hold": Hold,
    "LedgerAccountSubject": LedgerAccountSubject,
    "LedgerStorage": LedgerStorage,
    "ListAccounts200Response": ListAccounts200Response,
    "ListAccounts200ResponseCursor": ListAccounts200ResponseCursor,
    "ListAccounts200ResponseCursorAllOf": ListAccounts200ResponseCursorAllOf,
    "ListAccounts400Response": ListAccounts400Response,
    "ListBalancesResponse": ListBalancesResponse,
    "ListBalancesResponseCursor": ListBalancesResponseCursor,
    "ListBalancesResponseCursorAllOf": ListBalancesResponseCursorAllOf,
    "ListClientsResponse": ListClientsResponse,
    "ListConnectorTasks200ResponseInner": ListConnectorTasks200ResponseInner,
    "ListConnectorsConfigsResponse": ListConnectorsConfigsResponse,
    "ListConnectorsConfigsResponseConnector": ListConnectorsConfigsResponseConnector,
    "ListConnectorsConfigsResponseConnectorKey": ListConnectorsConfigsResponseConnectorKey,
    "ListConnectorsResponse": ListConnectorsResponse,
    "ListPaymentsResponse": ListPaymentsResponse,
    "ListScopesResponse": ListScopesResponse,
    "ListTransactions200Response": ListTransactions200Response,
    "ListTransactions200ResponseCursor": ListTransactions200ResponseCursor,
    "ListTransactions200ResponseCursorAllOf": ListTransactions200ResponseCursorAllOf,
    "ListUsersResponse": ListUsersResponse,
    "ListWalletsResponse": ListWalletsResponse,
    "ListWalletsResponseCursor": ListWalletsResponseCursor,
    "ListWalletsResponseCursorAllOf": ListWalletsResponseCursorAllOf,
    "Mapping": Mapping,
    "MappingResponse": MappingResponse,
    "ModulrConfig": ModulrConfig,
    "Monetary": Monetary,
    "Payment": Payment,
    "Posting": Posting,
    "Query": Query,
    "ReadClientResponse": ReadClientResponse,
    "ReadUserResponse": ReadUserResponse,
    "Response": Response,
    "RunScript400Response": RunScript400Response,
    "Scope": Scope,
    "ScopeAllOf": ScopeAllOf,
    "ScopeOptions": ScopeOptions,
    "Script": Script,
    "ScriptResult": ScriptResult,
    "Secret": Secret,
    "SecretAllOf": SecretAllOf,
    "SecretOptions": SecretOptions,
    "ServerInfo": ServerInfo,
    "Stats": Stats,
    "StatsResponse": StatsResponse,
    "StripeConfig": StripeConfig,
    "StripeTask": StripeTask,
    "StripeTransferRequest": StripeTransferRequest,
    "Subject": Subject,
    "TaskDescriptorBankingCircle": TaskDescriptorBankingCircle,
    "TaskDescriptorBankingCircleDescriptor": TaskDescriptorBankingCircleDescriptor,
    "TaskDescriptorCurrencyCloud": TaskDescriptorCurrencyCloud,
    "TaskDescriptorCurrencyCloudDescriptor": TaskDescriptorCurrencyCloudDescriptor,
    "TaskDescriptorDummyPay": TaskDescriptorDummyPay,
    "TaskDescriptorDummyPayDescriptor": TaskDescriptorDummyPayDescriptor,
    "TaskDescriptorModulr": TaskDescriptorModulr,
    "TaskDescriptorModulrDescriptor": TaskDescriptorModulrDescriptor,
    "TaskDescriptorStripe": TaskDescriptorStripe,
    "TaskDescriptorStripeDescriptor": TaskDescriptorStripeDescriptor,
    "TaskDescriptorWise": TaskDescriptorWise,
    "TaskDescriptorWiseDescriptor": TaskDescriptorWiseDescriptor,
    "Total": Total,
    "Transaction": Transaction,
    "TransactionData": TransactionData,
    "TransactionResponse": TransactionResponse,
    "Transactions": Transactions,
    "TransactionsResponse": TransactionsResponse,
    "UpdateWalletRequest": UpdateWalletRequest,
    "User": User,
    "Volume": Volume,
    "Wallet": Wallet,
    "WalletSubject": WalletSubject,
    "WalletWithBalances": WalletWithBalances,
    "WalletWithBalancesBalances": WalletWithBalancesBalances,
    "WalletsCursor": WalletsCursor,
    "WalletsErrorResponse": WalletsErrorResponse,
    "WalletsPosting": WalletsPosting,
    "WalletsTransaction": WalletsTransaction,
    "WalletsVolume": WalletsVolume,
    "WebhooksConfig": WebhooksConfig,
    "WebhooksCursor": WebhooksCursor,
    "WiseConfig": WiseConfig,
}

export class ObjectSerializer {
    public static findCorrectType(data: any, expectedType: string) {
        if (data == undefined) {
            return expectedType;
        } else if (primitives.indexOf(expectedType.toLowerCase()) !== -1) {
            return expectedType;
        } else if (expectedType === "Date") {
            return expectedType;
        } else {
            if (enumsMap.has(expectedType)) {
                return expectedType;
            }

            if (!typeMap[expectedType]) {
                return expectedType; // w/e we don't know the type
            }

            // Check the discriminator
            let discriminatorProperty = typeMap[expectedType].discriminator;
            if (discriminatorProperty == null) {
                return expectedType; // the type does not have a discriminator. use it.
            } else {
                if (data[discriminatorProperty]) {
                    var discriminatorType = data[discriminatorProperty];
                    if(typeMap[discriminatorType]){
                        return discriminatorType; // use the type given in the discriminator
                    } else {
                        return expectedType; // discriminator did not map to a type
                    }
                } else {
                    return expectedType; // discriminator was not present (or an empty string)
                }
            }
        }
    }

    public static serialize(data: any, type: string, format: string) {
        if (data == undefined) {
            return data;
        } else if (primitives.indexOf(type.toLowerCase()) !== -1) {
            return data;
        } else if (type.lastIndexOf("Array<", 0) === 0) { // string.startsWith pre es6
            let subType: string = type.replace("Array<", ""); // Array<Type> => Type>
            subType = subType.substring(0, subType.length - 1); // Type> => Type
            let transformedData: any[] = [];
            for (let index in data) {
                let date = data[index];
                transformedData.push(ObjectSerializer.serialize(date, subType, format));
            }
            return transformedData;
        } else if (type === "Date") {
            if (format == "date") {
                let month = data.getMonth()+1
                month = month < 10 ? "0" + month.toString() : month.toString()
                let day = data.getDate();
                day = day < 10 ? "0" + day.toString() : day.toString();

                return data.getFullYear() + "-" + month + "-" + day;
            } else {
                return data.toISOString();
            }
        } else {
            if (enumsMap.has(type)) {
                return data;
            }
            if (!typeMap[type]) { // in case we dont know the type
                return data;
            }

            // Get the actual type of this object
            type = this.findCorrectType(data, type);

            // get the map for the correct type.
            let attributeTypes = typeMap[type].getAttributeTypeMap();
            let instance: {[index: string]: any} = {};
            for (let index in attributeTypes) {
                let attributeType = attributeTypes[index];
                instance[attributeType.baseName] = ObjectSerializer.serialize(data[attributeType.name], attributeType.type, attributeType.format);
            }
            return instance;
        }
    }

    public static deserialize(data: any, type: string, format: string) {
        // polymorphism may change the actual type.
        type = ObjectSerializer.findCorrectType(data, type);
        if (data == undefined) {
            return data;
        } else if (primitives.indexOf(type.toLowerCase()) !== -1) {
            return data;
        } else if (type.lastIndexOf("Array<", 0) === 0) { // string.startsWith pre es6
            let subType: string = type.replace("Array<", ""); // Array<Type> => Type>
            subType = subType.substring(0, subType.length - 1); // Type> => Type
            let transformedData: any[] = [];
            for (let index in data) {
                let date = data[index];
                transformedData.push(ObjectSerializer.deserialize(date, subType, format));
            }
            return transformedData;
        } else if (type === "Date") {
            return new Date(data);
        } else {
            if (enumsMap.has(type)) {// is Enum
                return data;
            }

            if (!typeMap[type]) { // dont know the type
                return data;
            }
            let instance = new typeMap[type]();
            let attributeTypes = typeMap[type].getAttributeTypeMap();
            for (let index in attributeTypes) {
                let attributeType = attributeTypes[index];
                let value = ObjectSerializer.deserialize(data[attributeType.baseName], attributeType.type, attributeType.format);
                if (value !== undefined) {
                    instance[attributeType.name] = value;
                }
            }
            return instance;
        }
    }


    /**
     * Normalize media type
     *
     * We currently do not handle any media types attributes, i.e. anything
     * after a semicolon. All content is assumed to be UTF-8 compatible.
     */
    public static normalizeMediaType(mediaType: string | undefined): string | undefined {
        if (mediaType === undefined) {
            return undefined;
        }
        return mediaType.split(";")[0].trim().toLowerCase();
    }

    /**
     * From a list of possible media types, choose the one we can handle best.
     *
     * The order of the given media types does not have any impact on the choice
     * made.
     */
    public static getPreferredMediaType(mediaTypes: Array<string>): string {
        /** According to OAS 3 we should default to json */
        if (!mediaTypes) {
            return "application/json";
        }

        const normalMediaTypes = mediaTypes.map(this.normalizeMediaType);
        let selectedMediaType: string | undefined = undefined;
        let selectedRank: number = -Infinity;
        for (const mediaType of normalMediaTypes) {
            if (supportedMediaTypes[mediaType!] > selectedRank) {
                selectedMediaType = mediaType;
                selectedRank = supportedMediaTypes[mediaType!];
            }
        }

        if (selectedMediaType === undefined) {
            throw new Error("None of the given media types are supported: " + mediaTypes.join(", "));
        }

        return selectedMediaType!;
    }

    /**
     * Convert data to a string according the given media type
     */
    public static stringify(data: any, mediaType: string): string {
        if (mediaType === "text/plain") {
            return String(data);
        }

        if (mediaType === "application/json") {
            return JSON.stringify(data);
        }

        throw new Error("The mediaType " + mediaType + " is not supported by ObjectSerializer.stringify.");
    }

    /**
     * Parse data from a string according to the given media type
     */
    public static parse(rawData: string, mediaType: string | undefined) {
        if (mediaType === undefined) {
            throw new Error("Cannot parse content. No Content-Type defined.");
        }

        if (mediaType === "text/plain") {
            return rawData;
        }

        if (mediaType === "application/json") {
            return JSON.parse(rawData);
        }

        if (mediaType === "text/html") {
            return rawData;
        }

        throw new Error("The mediaType " + mediaType + " is not supported by ObjectSerializer.parse.");
    }
}
