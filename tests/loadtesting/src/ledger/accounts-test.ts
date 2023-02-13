import {withClient} from "../../libs/client";
import {DefaultLimit, ListAccountQuery} from "../../libs/client/1_3";
import {MatrixTest, withMatrix} from "../../libs/core";
import {fromServerVersion} from "../../libs/client/base";
// @ts-ignore
import {randomString} from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
// @ts-ignore
import {describe} from 'https://jslib.k6.io/k6chaijs/4.3.4.1/index.js';

export default () => {
    const [tx, metadata] = withClient('1.3', client => {
        const reference = randomString(20, 'aeioubcdfghijpqrstuv');
        const destination = randomString(20, 'aeioubcdfghijpqrstuv');
        const metadata = randomString(20, 'aeioubcdfghijpqrstuv');
        const tx = client.createTransaction(__ENV.LEDGER, {
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
        client.addMetadataToAccount(__ENV.LEDGER, destination, {
            [metadata]: 'Testing'
        });
        return [tx, metadata];
    });

    interface InternalTest extends MatrixTest {
        filter: Partial<ListAccountQuery>;
        name: string;
    }
    const tests: InternalTest[] = [
        {
            filter: {},
            name: 'all accounts',
            serverVersion: '1.3.0',
            clientVersion: '1.3'
        },
        {
            filter: {
                address: '.*' + tx.postings[0].destination.substr(3)
            },
            name: 'by address (regexp)',
            serverVersion: '1.3.0',
            clientVersion: '1.3'
        },
        {
            filter: {
                metadata: {
                    [metadata]: 'Testing'
                }
            },
            name: 'by metadata',
            serverVersion: '1.3.0',
            clientVersion: '1.3'
        }
    ];

    describe('List accounts', () => {
        withMatrix(tests, (t, client) => {
            client.listAccounts(__ENV.LEDGER, t.filter, {
                name: `List accounts: ${t.name}`
            });
        });
    });

    describe('Count accounts', () => {
        withMatrix(tests, (t, client) => {
            client.countAccounts(__ENV.LEDGER, t.filter, {
                name: `Count accounts: ${t.name}`
            });
        });
    });

    fromServerVersion('1.3.0', () => {
        withClient('1.3', client => {
            describe('Get account', () => {
                client.getAccount(__ENV.LEDGER, tx.postings[0].destination);
            });
            describe('Add metadata to account', () => {
                client.addMetadataToAccount(__ENV.LEDGER, tx.postings[0].destination, {
                    [randomString(20, 'aeioubcdfghijpqrstuv')]: randomString(20, 'aeioubcdfghijpqrstuv'),
                });
            });
            describe('Pagination on accounts', () => {
                let after: string|undefined;
                for(let i = 0 ; i < 100 ; i++) {
                    const params: Partial<ListAccountQuery> = {};
                    if (after) {
                        params.after = after;
                    }
                    const accounts = client.listAccounts(__ENV.LEDGER, params, {
                        name: `Pagination on accounts: Page ${i}`,
                    });
                    after = accounts[DefaultLimit-1].address;
                }
            });
        });
    });
};
