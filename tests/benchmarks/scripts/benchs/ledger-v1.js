import {ReadTransactions, WriteTransactions} from "../src/steps/transactions";
import {ReadAccounts} from "../src/steps/account";
import {ReadBalances} from "../src/steps/balances";
import {ReadAggregateBalances} from "../src/steps/aggregateBalances";
import {startLedger, stopLedger, exportResults} from 'k6/x/formancehq/benchmarks';
import exec from 'k6/execution';
import {K6Options} from "../src/options";

export function setup() {
    return startLedger({
        version: 'v1.10.4',
    });
}

export let options = K6Options();

const ledgerName = `/tests`;

export function read(ledger) {
    const url = ledger.url + ledgerName

    ReadTransactions(url);
    ReadAccounts(url);
    ReadBalances(url);
    ReadAggregateBalances(url);
}

export function write(ledger) {
    const url = ledger.url + ledgerName

    WriteTransactions(url, exec.scenario.iterationInTest);
}

export function teardown(data) {
    stopLedger();
    exportResults();
}
