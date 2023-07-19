import {Trend} from "k6/metrics";
import {check, group} from "k6";
import http from "k6/http";
import { URL } from 'https://jslib.k6.io/url/1.0.0/index.js';
import {generateTransactions, generateTransactionsNumscript} from "../generators/transactions";

const writeTransactionsWaitingTime = new Trend('waiting_time_transactions_write');
const readTransactionsWaitingTime = new Trend('waiting_time_transactions_read');


export function WriteTransactions(BASE_URL) {
    const url = new URL(`${BASE_URL}/transactions`);

    group('Transactions (WRITE)', function () {
        url.searchParams.append('async', "true");
        const res = http.post(url.toString(), JSON.stringify(generateTransactions()), {
            headers: { 'Content-Type': 'application/json' },
        });
        writeTransactionsWaitingTime.add(res.timings.waiting);

        check(res, {
            'is status 200': (r) => r.status === 200,
        });
    });

    group('Transactions - Numscript (WRITE)', function () {
        url.searchParams.append('async', "true");

        const res = http.post(url.toString(), JSON.stringify(generateTransactionsNumscript()), {
            headers: { 'Content-Type': 'application/json' },
        });
        writeTransactionsWaitingTime.add(res.timings.waiting);

        check(res, {
            'is status 200': (r) => r.status === 200,
        });
    });
}

export function ReadTransactions(BASE_URL) {
    const url = new URL(`${BASE_URL}/transactions`);
    const WALLET = "wallets:99:main";

    group('Transactions - NoParams (READ)', function () {
        const res = http.get(url.toString());
        readTransactionsWaitingTime.add(res.timings.waiting);
    
        check(res, {
            'is status 200': (r) => r.status === 200,
        });
    });
    group('Transactions - Specific Account (READ)', function () {
        url.searchParams.append('account', WALLET);

        const res = http.get(url.toString());
        readTransactionsWaitingTime.add(res.timings.waiting);
        check(res, {
            'is status 200': (r) => r.status === 200,
        });
    });
    group('Transactions - Specific Source (READ)', function () {
        url.searchParams.append('source', WALLET);

        const res = http.get(url.toString());
        readTransactionsWaitingTime.add(res.timings.waiting);
        check(res, {
            'is status 200': (r) => r.status === 200,
        });
    });
    group('Transactions - Specific Destination (READ)', function () {
        url.searchParams.append('destination', WALLET);

        const res = http.get(url.toString());
        readTransactionsWaitingTime.add(res.timings.waiting);
        check(res, {
            'is status 200': (r) => r.status === 200,
        });
    });
    group('Transactions - Pattern Account (READ)', function () {
        url.searchParams.append('account', 'wallet:*:main');

        const res = http.get(url.toString());
        readTransactionsWaitingTime.add(res.timings.waiting);
        check(res, {
            'is status 200': (r) => r.status === 200,
        });
    });
}