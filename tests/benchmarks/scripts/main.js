import http from 'k6/http';
import { check } from 'k6';
import extension from 'k6/x/formancehq/benchmarks';
import { URL } from 'https://jslib.k6.io/url/1.0.0/index.js';

export function setup() {
    return extension.startLedger({
        //version: 'v1.10.3', // Can be passed using "LEDGER_VERSION" env var
        version: '13644f2fe711feb83948aeec5732a4d9e47389d5'
    });
}

export default function (ledger) {
    const url = new URL(`${ledger.url}/default/transactions`);
    const res = http.post(url.toString(), JSON.stringify({
        postings: [{
            source: 'world',
            destination: 'bank',
            amount: 100,
            asset: 'USD/2',
        }]
    }), {
        headers: { 'Content-Type': 'application/json' },
    });
    check(res, {
        'is status 200': (r) => r.status === 200,
    });
}

export function teardown(data) {
    extension.stopLedger();
}
