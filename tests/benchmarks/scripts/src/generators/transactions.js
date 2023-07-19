import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

const vendor = `vendor:${randomIntBetween(1, 1000)}:main`;
const wallet = `wallets:${randomIntBetween(1, 1000)}:main`;
const paymentIn = `payments:in:${randomIntBetween(1, 1000)}`;
const paymentOut = `payments:out:${randomIntBetween(1, 1000)}`;
const num = open('./../src/generators/simple.numscript', 'r');
export const generateTransactions = () => (
{
    postings: [
        {
            amount: 100e2,
            asset: 'EUR/2',
            source: `world`,
            destination: paymentIn,
        },
        {
            amount: 100e2,
            asset: 'EUR/2',
            source: paymentIn,
            destination: wallet,
        },
    ],
});

export const generateTransactionsNumscript = () => (
    {
        script: {
            plain: num,
            vars: {
                "wallet": wallet,
                "payment_in": paymentIn,
            }
        },
    }
);