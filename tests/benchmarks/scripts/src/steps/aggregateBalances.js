import {Trend} from "k6/metrics";
import {check, group} from "k6";
import http from "k6/http";
import { URL } from 'https://jslib.k6.io/url/1.0.0/index.js';

const readAggregateBalancesWaitingTime = new Trend('waiting_time_aggregateBalances_read');
const WALLET = "wallets:99:main";

export function ReadAggregateBalances(BASE_URL) {
    const url = new URL(`${BASE_URL}/aggregate/balances`);
    // group('AggregateBalances - NoParams (READ)', function () {
    //     const res = http.get(url.toString());
    //     readAggregateBalancesWaitingTime.add(res.timings.waiting);
    //
    //     check(res, {
    //         'is status 200': (r) => r.status === 200,
    //     });
    // });
    group('AggregateBalances - filter by address (READ)', function () {
        url.searchParams.append('address', WALLET);

        const res = http.get(url.toString());
        readAggregateBalancesWaitingTime.add(res.timings.waiting);
        check(res, {
            'is status 200': (r) => r.status === 200,
        });
    });
    group('AggregateBalances - filter by address pattern (READ)', function () {
        url.searchParams.append('address', 'wallet:*:main');

        const res = http.get(url.toString());
        readAggregateBalancesWaitingTime.add(res.timings.waiting);
        check(res, {
            'is status 200': (r) => r.status === 200,
        });
    });
}
