import {Account, Config, Headers, ListTransactionQuery as ListTransactionQuery1_3, Transaction} from "../1_3";
import {Client as Client1_4} from '../1_4';
import {loadConfig} from "../../config";

export interface ListTransactionQuery extends ListTransactionQuery1_3 {
    start_time: Date;
    end_time: Date;
}

export class Client extends Client1_4 {

    version = '1.5';

    constructor(
        config: Config
    ) {
        super(config);
    }

    listTransactions(ledger: string, query?: Partial<ListTransactionQuery>, tags?: { [name: string]: string }): Transaction[] {
        return super.listTransactions(ledger, {
            ...query,
            // @ts-ignore
            end_time: query?.end_time ? query.end_time.toISOString() : '',
            start_time: query?.start_time ? query.start_time.toISOString() : '',
        }, tags);
    }

    countTransactions(ledger: string, query?: Partial<ListTransactionQuery>, tags?: { [p: string]: string }): number {
        return super.countTransactions(ledger, {
            ...query,
            // @ts-ignore
            end_time: query?.end_time ? query.end_time.toISOString() : '',
            start_time: query?.start_time ? query.start_time.toISOString() : '',
        }, tags);
    }

    listAccounts(ledger: string): Account[] {
        return super.listAccounts(ledger);
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
