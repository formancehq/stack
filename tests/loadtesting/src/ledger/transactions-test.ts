import {DefaultLimit, ListTransactionQuery as ListTransactionQuery1_3} from "../../libs/client/1_3";
// @ts-ignore
import {randomString} from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import {withClient} from "../../libs/client";
import {MatrixTest, withMatrix} from "../../libs/core";
import {ListTransactionQuery as ListTransactionQuery1_5} from "../../libs/client/1_5";
import {fromServerVersion} from "../../libs/client/base";
// @ts-ignore
import {describe, expect} from 'https://jslib.k6.io/k6chaijs/4.3.4.1/index.js';

export default () => {

    const reference = randomString(20, 'aeioubcdfghijpqrstuv');
    const destination = randomString(20, 'aeioubcdfghijpqrstuv');
    const tx = withClient('1.3', client => {
        return client.createTransaction(__ENV.LEDGER, {
            postings: [
                {
                    asset: 'USD',
                    amount: 100,
                    source: 'world',
                    destination
                }
            ],
            reference,
        });
    });

    interface InternalTest extends MatrixTest {
        filter: Partial<ListTransactionQuery1_3|ListTransactionQuery1_5>;
        name: string;
    }
    const tests: InternalTest[] = [
        {
            filter: {},
            name: 'all transactions',
            serverVersion: '1.3.0',
            clientVersion: '1.3'
        },
        {
            filter: {
                destination: tx.postings[0].destination
            },
            name: 'by destination',
            serverVersion: '1.3.0',
            clientVersion: '1.3'
        },
        {
            filter: {
                account: tx.postings[0].destination
            },
            name: 'by account',
            serverVersion: '1.3.0',
            clientVersion: '1.3'
        },
        {
            filter: {
                reference: tx.reference
            },
            name: 'by reference',
            serverVersion: '1.3.0',
            clientVersion: '1.3'
        },
        {
            filter: {
                start_time: new Date(tx.timestamp),
                end_time: new Date(new Date(tx.timestamp).getTime() + 1000)
            },
            name: 'by timerange',
            serverVersion: '1.5.0',
            clientVersion: '1.5'
        },
    ];

    describe('List transactions', () => {
        withMatrix(tests, (t, client) => {
            client.listTransactions(__ENV.LEDGER, t.filter, {
                name: `List transactions: ${t.name}`
            });
        });
    });
    describe('Count transactions', () => {
        withMatrix(tests, (t, client) => {
            client.countTransactions(__ENV.LEDGER, t.filter, {
                name: `Count transactions: ${t.name}`
            });
        });
    });

    fromServerVersion('1.3.0', () => {
        withClient('1.3', client => {
            describe('Get transaction', () => {
                client.getTransaction(__ENV.LEDGER, tx.txid);
            });
            describe('Add metadata to transaction', () => {
                client.addMetadataToTransaction(__ENV.LEDGER, tx.txid, {
                    [randomString(20, 'aeioubcdfghijpqrstuv')]: randomString(20, 'aeioubcdfghijpqrstuv'),
                });
            });
            describe('Pagination on transactions', () => {
                let after = client.countTransactions(__ENV.LEDGER, {}, {
                    name: 'Count transactions: all transactions'
                });
                for(let i = 0 ; i < 100 ; i++) {
                    const txs = client.listTransactions(__ENV.LEDGER, {
                        after
                    }, {
                        name: `Pagination on transactions: Page ${i}`,
                    });
                    after = txs[DefaultLimit-1].txid;
                }
            });
            describe('Get transaction with 500 Page Size', () => {
                client.listTransactions(__ENV.LEDGER,{
                    source: `.*:.*:withdrawals:.*`,
                    page_size: 500,
                });
            });
            describe('Get transaction with default Page Size', () => {
                client.listTransactions(__ENV.LEDGER,{
                    source: `.*:.*:withdrawals:.*`,
                });
            });
        });
    });
};
