import {Rate, Trend} from "k6/metrics";
import {check, group} from "k6";
import http from "k6/http";
import { URL } from 'https://jslib.k6.io/url/1.0.0/index.js';

const readBalancesWaitingTime = new Trend('waiting_time_balances_read');
const WALLET = "wallets:99:main";

export function ReadBalances(BASE_URL) {
    const url = new URL(`${BASE_URL}/balances`);

    group('Balances - NoParams (READ)', function () {
        const res = http.get(url.toString());
        readBalancesWaitingTime.add(res.timings.waiting);
        check(res, {
            'is status 200': (r) => r.status === 200,
        });
    });
    group('Balances - filter by address (READ)', function () {
        url.searchParams.append('address', WALLET);

        const res = http.get(url.toString());
        readBalancesWaitingTime.add(res.timings.waiting);
        check(res, {
            'is status 200': (r) => r.status === 200,
        });
    });
    group('Balances - filter by address pattern (READ)', function () {
        url.searchParams.append('address', 'wallet:*:main');

        const res = http.get(url.toString());
        readBalancesWaitingTime.add(res.timings.waiting);
        check(res, {
            'is status 200': (r) => r.status === 200,
        });
    });
}