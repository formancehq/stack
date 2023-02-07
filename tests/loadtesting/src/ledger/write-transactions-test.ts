import {TransactionData} from "../../libs/client/1_3";
import {withClient} from "../../libs/client";
import {fromServerVersion} from "../../libs/client/base";
import uniqid from 'uniqid';
// @ts-ignore
import {describe} from 'https://jslib.k6.io/k6chaijs/4.3.4.1/index.js';
// @ts-ignore
import {randomString} from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

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
    const trade = (wallet: string, amountFiat: number, asset: string, amountAsset: number) : any => {
        const tradeId = uniqid().toUpperCase();

        return {
            postings: [
                {
                    amount: amountFiat,
                    asset: 'EUR/2',
                    source: wallet,
                    destination: `trades:${tradeId}`,
                },
                {
                    amount: amountFiat,
                    asset: 'EUR/2',
                    source: `trades:${tradeId}`,
                    destination: 'formance:fiat:holdings',
                },
                {
                    amount: amountAsset,
                    asset,
                    source: 'teller:otc:nyse',
                    destination: `trades:${tradeId}`,
                },
                {
                    amount: amountAsset,
                    asset,
                    source: `trades:${tradeId}`,
                    destination: wallet,
                },
            ],
        };
    };
    fromServerVersion('1.3.0', () => {
        withClient('1.3', client => {
            describe('Create transaction', () => {
                client.createTransaction(__ENV.LEDGER, {
                    postings: [{
                        source: 'world',
                        destination: tx.postings[0].destination,
                        amount: 100,
                        asset: 'USD',
                    }]
                });
            });
            describe('Batch transaction', () => {
                let txs: TransactionData[] = [];
                for (let i = 0 ; i < 100 ; i++) {
                    txs.push({
                        postings: [{
                            source: 'world',
                            destination: 'bank',
                            amount: 100,
                            asset: 'USD'
                        }]
                    });
                }
                client.createTransactions(__ENV.LEDGER, txs);
            });
            describe('Run script', () => {
                client.runScript(__ENV.LEDGER, {
                    plain: `send [COIN 100] (
                      source = @world
                      destination = @centralbank
                    )`,
                });
            });
            describe('Batch transaction Formance', () => {
                const userId = randomString(20, 'aeioubcdfghijpqrstuv');
                const wallet = `users:${userId}:wallet`;
                const id = uniqid().toUpperCase();
                const txs: TransactionData[] = [];
                txs.push({
                    postings: [
                        {
                            amount: 100e2,
                            asset: 'EUR/2',
                            source: `world`,
                            destination: `payments:adyen:${id}`,
                        },
                    ],
                });

                txs.push({
                    postings: [
                        {
                            amount: 100e2,
                            asset: 'EUR/2',
                            source: `payments:adyen:${id}`,
                            destination: wallet,
                        }
                    ],
                });

                txs.push({
                    postings: [
                        {
                            amount: 0.35e6,
                            asset: 'RBLX/6',
                            source: 'world',
                            destination: 'teller:otc:nyse',
                        },
                        {
                            amount: 1.84e6,
                            asset: 'SNAP/6',
                            source: 'world',
                            destination: 'teller:otc:nyse',
                        }
                    ],
                });

                txs.push(trade(wallet, 15e2, 'RBLX/6', 0.35e6));
                txs.push(trade(wallet, 42.3e2, 'SNAP/6', 1.84e6));

                const withdrawal = `users:${userId}:withdrawals:${uniqid()}`;

                txs.push({
                    postings: [
                        {
                            amount: 22.7e2,
                            asset: 'EUR/2',
                            source: wallet,
                            destination: withdrawal,
                        },
                    ],
                });

                txs.push({
                    postings: [
                        {
                            amount: 22.7e2,
                            asset: 'EUR/2',
                            source: withdrawal,
                            destination: `payments:${uniqid()}`,
                        },
                    ],
                });

                client.createTransactions(__ENV.LEDGER, txs);
            });
        });
    });
};
