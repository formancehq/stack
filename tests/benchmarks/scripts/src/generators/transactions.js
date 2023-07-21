import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

const vendor = `vendor:${randomIntBetween(1, 1000)}:main`;
const paymentOut = `payments:out:${randomIntBetween(1, 1000)}`;
const num = open('./../src/generators/simple.numscript', 'r');

export function generateTransactions() {
    const paymentIn = `payments:in:${randomIntBetween(1, 1000)}`;
    const wallet = `wallets:${randomIntBetween(1, 1000)}:main`;
    return {
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
                source:`world`,
                destination: wallet,
            },
        ],
    }
}

export function generateTransactionsNumscript() {
    const paymentIn = `payments:in:${randomIntBetween(1, 1000)}`;
    const wallet = `wallets:${randomIntBetween(1, 1000)}:main`;
    return {
        script: {
            plain: num,
            vars: {
                "payment_in": paymentIn,
                "wallets": wallet,
                "world": `world:${randomIntBetween(1, 100000)}`,
            }
        },
    }
}