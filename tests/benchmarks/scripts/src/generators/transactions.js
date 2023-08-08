const num = open('./../src/generators/simple.numscript', 'r');

export function generateTransactions(iterationInTest) {
    const paymentIn = `payments:in:${iterationInTest}`;
    const wallet = `wallets:${iterationInTest}:main`;
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
                source: paymentIn,
                destination: wallet,
            },
        ],
    }
}

export function generateTransactionsNumscript(iterationInTest) {
    const paymentIn = `payments:in:${iterationInTest}`;
    const wallet = `wallets:${iterationInTest}:main`;
    return {
        script: {
            plain: num,
            vars: {
                "payment_in": paymentIn,
                "wallets": wallet,
                "world": `world:${iterationInTest}`,
            }
        },
    }
}