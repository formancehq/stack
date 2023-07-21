import {ReadTransactions, WriteTransactions} from "../src/steps/transactions";
import {ReadAccounts} from "../src/steps/account";
import {ReadBalances} from "../src/steps/balances";
import {ReadAggregateBalances} from "../src/steps/aggregateBalances";
import {startLedger, stopLedger, exportResults} from 'k6/x/formancehq/benchmarks';

export function setup() {
    return startLedger({
        //version: 'v1.10.3', // Can be passed using "LEDGER_VERSION" env var
        version: '840ebff87373cf8f78dcb3a691619f8db0430baa',
    });
}

export let options = {
    scenarios: {
        read_constant: {
            executor: 'constant-vus',
            exec: 'read',
            vus: 20,
            duration: '10m',
            tags: { testid: __ENV.TEST_ID}
        },
    }
};

export function read(ledger) {
    const url = ledger.url + "/tests02"

    ReadTransactions(url);
    ReadAccounts(url);
    ReadBalances(url);
    ReadAggregateBalances(url);
}

export function teardown(data) {
    stopLedger();
    exportResults();
}
