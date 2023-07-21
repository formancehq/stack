import {ReadTransactions, WriteTransactions} from "../src/steps/transactions";
import {ReadAccounts} from "../src/steps/account";
import {ReadBalances} from "../src/steps/balances";
import {ReadAggregateBalances} from "../src/steps/aggregateBalances";
import {startLedger, stopLedger, exportResults} from 'k6/x/formancehq/benchmarks';

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
            vus: 50,
            duration: '2m',
            tags: { testid: __ENV.TEST_ID}
        },
        write_timer: {
            executor: 'ramping-vus',
            exec: 'write',
            startVUs: 0,
            stages: [
                { duration: '10m', target: 0 },
                { duration: '1s', target: 1 },
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

export function teardown(data) {
    stopLedger();
    exportResults();
}
