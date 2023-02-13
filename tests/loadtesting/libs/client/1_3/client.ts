import http, {RefinedParams, RefinedResponse} from "k6/http";
import {BaseClient} from "../base";
import {loadConfig} from "../../config";
import {Trend} from "k6/metrics";
// @ts-ignore
import {Validator} from 'k6/x/openapi';
// @ts-ignore
import {URL} from 'https://jslib.k6.io/url/1.0.0/index.js';

const listTransactionsTrend = new Trend('list_transactions', true);

export const DefaultLimit = 15;

interface Info {
    server: string;
    version: string;
    servers: any[];
    config: {
        storage: {
            driver: string
            ledgers: string[]
        }
    };
}

export interface Posting {
    source: string;
    destination: string;
    asset: string;
    amount: number;
}

export interface Metadata {
    [key: string]: any;
}

export interface TransactionData {
    postings: Posting[];
    reference?: string;
    metadata?: Metadata;
}

export interface Script {
    plain: string;
    vars?: {[name: string]: any};
}

export interface Transaction extends TransactionData {
    txid: number;
    timestamp: string;
}

export interface Account {
    address: string;
}

export interface Stats {
    transactions: number;
    accounts: number;
}

export interface Headers { [name: string]: string; }

export interface Config {
    headers: Headers;
    endpoint: string;
}

export interface ListTransactionQuery {
    after: number;
    reference: string;
    account: string;
    source: string;
    destination: string;
    page_size: number;
}

export interface ListAccountQuery {
    after: string;
    address: string;
    metadata: { [path: string]: any };
}

export const deepObjectSerialize = (v: any): { [path: string]: any } => {
    const encodedParams: {[path: string]: any} = {};
    for(const key in v) {
        if(v.hasOwnProperty(key)) {
            encodedParams[`metadata[${key}]`] = v[key];
        }
    }
    return encodedParams;
};

export class Client implements BaseClient {

    version = '1.3';
    private validator: Validator;

    constructor(
        protected config: Config
    ) {
        const swagger = this.getSwagger();
        swagger.servers = [];
        this.validator = Validator(swagger);
    }

    request(name: string, parameters: {
        method: string,
        path: string,
        body?: any,
        queryParams?: { [path: string]: any }
        tags?: { [name: string]: string },
    }): RefinedResponse<'text'> {
        const headers = this.config.headers;
        headers['Content-Type'] = 'application/json';

        const url = new URL(this.config.endpoint + parameters.path);
        if(parameters.queryParams) {
            for (const key in parameters.queryParams) {
                if(parameters.queryParams.hasOwnProperty(key)) {
                    url.searchParams.set(key, parameters.queryParams[key]);
                }
            }
        }

        const params: RefinedParams<'text'> = {
            headers,
            tags: {
                name,
                ...parameters.tags,
            },
        };
        if (__ENV.HTTP_REQUEST_TIMEOUT && __ENV.HTTP_REQUEST_TIMEOUT !== '') {
            params.timeout = __ENV.HTTP_REQUEST_TIMEOUT;
        }

        const requestBody = JSON.stringify(parameters.body);
        const response = http.request<"text">(
            parameters.method,
            url.toString(),
            requestBody,
            params
        );

        if (this.validator) {
            const requestHeadersToValidate: { [key: string]: string[] } = {};
            for(const key in headers) {
                if(headers.hasOwnProperty(key)) {
                    requestHeadersToValidate[key] = [headers[key]];
                }
            }

            const responseHeadersToValidate: { [key: string]: string[] } = {};
            for(const key in response.headers) {
                if(response.headers.hasOwnProperty(key)) {
                    responseHeadersToValidate[key] = [response.headers[key]];
                }
            }

            const error = this.validator.validate({
                url: url.toString(),
                method: parameters.method,
                headers: requestHeadersToValidate,
                body: requestBody
            }, {
                statusCode: response.status,
                headers: responseHeadersToValidate,
                body: response.body,
            });
            if (error) {
                throw error;
            }
        }

        return response;
    }

    getInfo(): Info {
        const res = this.request("GetInfo", {
            method: 'GET',
            path: '/_info',
        });
        if(res.status !== 200) {
            throw new Error("unexpected status code: " + res.status);
        }
        return res.json('data') as any;
    }

    getSwagger(): Info {
        const res = this.request("GetSwagger", {
            method: 'GET',
            path: '/swagger.json',
        });
        if(res.status !== 200) {
            throw new Error("unexpected status code: " + res.status);
        }
        return res.json() as any;
    }

    listTransactions(ledger: string, query?: Partial<ListTransactionQuery>, tags?: { [name: string]: string }): Transaction[] {
        const res = this.request("ListTransactions", {
            method: 'GET',
            path: "/" + ledger + "/transactions",
            queryParams: query,
            tags,
        });
        listTransactionsTrend.add(res.timings.duration, tags);
        if(res.status !== 200) {
            throw new Error("unexpected status code: " + res.status);
        }
        return res.json('cursor.data') as any;
    }

    countTransactions(ledger: string, query?: Partial<ListTransactionQuery>, tags?: { [name: string]: string }): number {
        const res = this.request("CountTransactions", {
            method: 'HEAD',
            path: "/" + ledger + "/transactions",
            queryParams: query,
            tags,
        });
        if(res.status !== 200) {
            throw new Error("unexpected status code: " + res.status);
        }
        return parseInt(res.headers.Count, 10);
    }

    listAccounts(ledger: string, query?: Partial<ListAccountQuery>, tags?: { [name: string]: string }): Account[] {
        let encodedParams: any = query;
        if(query?.metadata) {
            encodedParams= {...query, ...deepObjectSerialize(query.metadata)};
            delete encodedParams.metadata;
        }

        const res = this.request("ListAccounts", {
            method: 'GET',
            path: '/' + ledger + '/accounts',
            queryParams: encodedParams,
            tags,
        });
        if(res.status !== 200) {
            throw new Error("unexpected status code: " + res.status);
        }
        return res.json('cursor.data') as any;
    }

    countAccounts(ledger: string, query?: Partial<ListAccountQuery>, tags?: { [name: string]: string }): number {
        let encodedParams: any = query;
        if(query?.metadata) {
            encodedParams= {...query, ...deepObjectSerialize(query.metadata)};
            delete encodedParams.metadata;
        }

        const res = this.request("CountAccounts", {
            method: 'HEAD',
            path: "/" + ledger + "/accounts",
            queryParams: encodedParams,
            tags
        });
        if(res.status !== 200) {
            throw new Error("unexpected status code: " + res.status);
        }
        return parseInt(res.headers.Count, 10);
    }

    getAccount(ledger: string, address: string): Account {
        const res = this.request("GetAccount", {
            method: 'GET',
            path: '/' + ledger + '/accounts/' + address,
        });
        if(res.status !== 200) {
            throw new Error("unexpected status code: " + res.status);
        }
        return res.json('data') as any;
    }

    createTransaction(ledger: string, data: TransactionData): Transaction {
        const res = this.request("CreateTransaction", {
            method: 'POST',
            body: data,
            path: "/" + ledger + "/transactions"
        });
        if(res.status !== 200) {
            throw new Error("unexpected status code: " + res.status);
        }
        return (res.json('data') as any)[0];
    }

    getTransaction(ledger: string, id: number): Transaction {
        const res = this.request("GetTransaction", {
            method: 'GET',
            path: "/" + ledger + "/transactions/" + id
        });
        if(res.status !== 200) {
            throw new Error("unexpected status code: " + res.status);
        }
        return res.json('data') as any;
    }

    addMetadataToAccount(ledger: string, account: string, metadata: any): void {
        const res = this.request("AddMetadataToAccount", {
            method: 'POST',
            body: metadata,
            path: "/" + ledger + "/accounts/" + account + "/metadata",
        });
        if(res.status !== 204) {
            throw new Error("unexpected status code: " + res.status);
        }
    }

    addMetadataToTransaction(ledger: string, transaction: number, metadata: any): void {
        const res = this.request("AddMetadataToTransaction", {
            method: 'POST',
            body: metadata,
            path: "/" + ledger + "/transactions/" + transaction + "/metadata",
        });
        if(res.status !== 204) {
            throw new Error("unexpected status code: " + res.status);
        }
    }

    createTransactions(ledger: string, data: TransactionData[]): Transaction[] {
        const res = this.request("BatchTransactions", {
            method: 'POST',
            path: "/" + ledger + "/transactions/batch",
            body: {
                transactions: data
            }
        });
        if(res.status !== 200) {
            throw new Error("unexpected status code: " + res.status);
        }
        return res.json('data') as any;
    }

    runScript(ledger: string, script: Script): Transaction {
        const res = this.request("RunScript", {
            method: 'POST',
            path: '/' + ledger + '/script',
            body: script,
        });
        if(res.status !== 200) {
            throw new Error("unexpected status code: " + res.status);
        }
        return res.json('transaction') as any;
    }

    getStats(ledger: string): Stats {
        const res = this.request("GetStats", {
            method: 'GET',
            path: "/" + ledger + "/stats"
        });
        if(res.status !== 200) {
            throw new Error("unexpected status code: " + res.status);
        }
        return res.json('data') as any;
    }
}

export const newClient = () => {
    const config = loadConfig();
    const headers: Headers = {};
    return new Client({
        endpoint: config.ledgerUrl,
        headers
    });
};
