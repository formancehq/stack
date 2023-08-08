import {ReadTransactions, WriteTransactions} from "../src/steps/transactions";
import {ReadAccounts} from "../src/steps/account";
import {ReadBalances} from "../src/steps/balances";
import {ReadAggregateBalances} from "../src/steps/aggregateBalances";
import {K6Options} from "../src/options";
import {startLedger, stopLedger, exportResults} from 'k6/x/formancehq/benchmarks';
import exec from 'k6/execution';

export function setup() {
    return startLedger({
        version: '88b73587a40e320defa0a12f5a555537458c47a7',
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
