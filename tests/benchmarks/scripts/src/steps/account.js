import {Trend} from "k6/metrics";
import {check, group} from "k6";
import http from "k6/http";
import { URL } from 'https://jslib.k6.io/url/1.0.0/index.js';

const readAccountsWaitingTime = new Trend('waiting_time_accounts_read');
const WALLET = "wallets:99:main";

export function ReadAccounts(BASE_URL) {
    const url = new URL(`${BASE_URL}/accounts`);

    group('Accounts - NoParams (READ)', function () {
        const res = http.get(url.toString());
        readAccountsWaitingTime.add(res.timings.waiting);
        check(res, {
            'is status 200': (r) => r.status === 200,
        });
    });
    group('Accounts - filter by address (READ)', function () {
        url.searchParams.append('address', WALLET);

        const res = http.get(url.toString());
        readAccountsWaitingTime.add(res.timings.waiting);
        check(res, {
            'is status 200': (r) => r.status === 200,
        });
    });
    group('Accounts - filter by address pattern (READ)', function () {
        url.searchParams.append('address', 'wallet::main');

        const res = http.get(url.toString());
        readAccountsWaitingTime.add(res.timings.waiting);
        check(res, {
            'is status 200': (r) => r.status === 200,
        });
    });
}
