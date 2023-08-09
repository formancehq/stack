import {ReadTransactions, WriteTransactions} from "../src/steps/transactions";
import {ReadAccounts} from "../src/steps/account";
import {ReadBalances} from "../src/steps/balances";
import {ReadAggregateBalances} from "../src/steps/aggregateBalances";
import {K6Options} from "../src/options";
import {startLedger, stopLedger, exportResults} from 'k6/x/formancehq/benchmarks';
import exec from 'k6/execution';

export function setup() {
    return startLedger({
        version: '15a430c40e95d38d864599dece235bc3964a3588',
    });
}

export let options = K6Options();

const ledgerName = `/tests`;

export function read_transactions(ledger) {
    const url = ledger.url + ledgerName
    ReadTransactions(url);
}
export function read_accounts(ledger) {
    const url = ledger.url + ledgerName
    ReadAccounts(url);
}
export function read_balances(ledger) {
    const url = ledger.url + ledgerName
    ReadBalances(url);
}
export function read_aggregatebalances(ledger) {
    const url = ledger.url + ledgerName
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
