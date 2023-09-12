import {ReadTransactions, WriteTransactions} from "../src/steps/transactions";
import {ReadAccounts} from "../src/steps/account";
import {ReadBalances} from "../src/steps/balances";
import {ReadAggregateBalances} from "../src/steps/aggregateBalances";
import {K6Options} from "../src/options";
import {startLedger, stopLedger} from 'k6/x/formancehq/benchmarks';
import exec from 'k6/execution';

export function setup() {
    return startLedger({
        version: 'latest',
    });
}

export let options = K6Options();

const ledgerName = `/tests`;

export function readTransactions(ledger) {
    const url = ledger.url + ledgerName
    ReadTransactions(url);
}
export function readAccounts(ledger) {
    const url = ledger.url + ledgerName
    ReadAccounts(url);
}
export function readBalances(ledger) {
    const url = ledger.url + ledgerName
    ReadBalances(url);
}
export function readAggregatedBalances(ledger) {
    const url = ledger.url + ledgerName
    ReadAggregateBalances(url);
}

export function write(ledger) {
    const url = ledger.url + ledgerName

    WriteTransactions(url, exec.scenario.iterationInTest);
}

export function teardown(data) {
    stopLedger();
}
