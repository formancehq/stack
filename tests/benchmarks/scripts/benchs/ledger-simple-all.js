import {ReadTransactions, WriteTransactions} from "../src/steps/transactions";
import {ReadAccounts} from "../src/steps/account";
import {ReadBalances} from "../src/steps/balances";
import {ReadAggregateBalances} from "../src/steps/aggregateBalances";
import {startLedger, stopLedger, exportResults} from 'k6/x/formancehq/benchmarks';

export function setup() {
    return startLedger({
        //version: 'v1.10.3', // Can be passed using "LEDGER_VERSION" env var
        version: '13644f2fe711feb83948aeec5732a4d9e47389d5',        
    });
}

export let options = {
    scenarios: {
        write_constant: {
            executor: 'constant-vus',
            exec: 'write',
            vus: 5,
            duration: '10m',
            tags: { testid: __ENV.TEST_ID}
        },
        read_constant: {
            executor: 'constant-vus',
            exec: 'read',
            vus: 20,
            duration: '10m',
            tags: { testid: __ENV.TEST_ID}
        },
        // read_spike: {
        //     executor: 'ramping-vus',
        //     exec: 'read',
        //     startVUs: 0,
        //     stages: [
        //         { duration: '2m', target: 0 },
        //         { duration: '10s', target: 20 },
        //         { duration: '2m', target: 0 },
        //         { duration: '20s', target: 20 },
        //         { duration: '2m', target: 0 },
        //         { duration: '10s', target: 20 },
        //         { duration: '2m', target: 0 },
        //         { duration: '20s', target: 20 },
        //         { duration: '1m', target: 0 },
        //     ],
        //     gracefulRampDown: '0s',
        //     tags: { testid: __ENV.TEST_ID}
        // }
    }
    //     read: {
    //         executor: 'per-vu-iterations',
    //         exec: 'read',
    //         vus: 10,
    //         iterations: 100,
    //         maxDuration: '10m',
    //     },
    //     read2: {
    //         exec: 'read',
    //         executor: 'constant-arrival-rate',
    //         rate: 90,
    //         timeUnit: '1m',
    //         duration: '5m',
    //         preAllocatedVUs: 10,
    //     },
    //     write: {
    //         executor: 'per-vu-iterations',
    //         exec: 'write',
    //         vus: 10,
    //         iterations: 100,
    //         maxDuration: '10m',
    //     },
    // },
};

export function read(ledger) {
    const url = ledger.url + "/tests01"

    ReadTransactions(url);
    ReadAccounts(url);
    ReadBalances(url);
    // ReadAggregateBalances(url);
}

export function write(ledger) {
    const url = ledger.url + "/tests01"

    WriteTransactions(url);
}

export function teardown(data) {
    stopLedger();
    exportResults();
}
