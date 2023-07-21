import {startLedger, stopLedger, checkPrometheusQuery} from 'k6/x/formancehq/benchmarks';
import {ReadTransactions, WriteTransactions} from "../src/steps/transactions";
import {ReadAccounts} from "../src/steps/account";
import {ReadBalances} from "../src/steps/balances";


export function setup() {
    return startLedger({
        // version: 'v1.10.3', // Can be passed using "LEDGER_VERSION" env var
        version: '840ebff87373cf8f78dcb3a691619f8db0430baa',
    });
}

export let options = {
    scenarios: {
        write_constant: {
            executor: 'constant-vus',
            exec: 'write',
            vus: 1,
            duration: '1m',
            tags: { testid: __ENV.TEST_ID}
        },
        read_timer: {
            executor: 'ramping-vus',
            exec: 'read',
            startVUs: 0,
            stages: [
                { duration: '2s', target: 1 },
                { duration: '5m', target: 1 },
            ],
            gracefulRampDown: '0s',
            tags: { testid: __ENV.TEST_ID}
        }
    }
};

export function write(ledger) {
    const url = ledger.url + "/tests01"
    WriteTransactions(url);
}

export function read(ledger) {
    const url = ledger.url + "/tests01"

    ReadTransactions(url);
}

export function teardown(data) {
    checkPrometheusQuery(`(sum(ledger_query_inbound_logs{testid="${__ENV.TEST_ID}"}) - sum(ledger_query_processed_logs{testid="${__ENV.TEST_ID}"}))`, 0);
    stopLedger();
}
