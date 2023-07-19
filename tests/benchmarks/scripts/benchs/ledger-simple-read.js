import {ReadTransactions} from "../src/steps/transactions";
import {ReadAccounts} from "../src/steps/account";
import {ReadBalances} from "../src/steps/balances";
import {ReadAggregateBalances} from "../src/steps/aggregateBalances";
import extension from 'k6/x/formancehq/benchmarks';

export function setup() {
    return extension.startLedger({
        //version: 'v1.10.3', // Can be passed using "LEDGER_VERSION" env var
        version: '13644f2fe711feb83948aeec5732a4d9e47389d5'
    });
}

export let options = {
    scenarios: {
        contacts: {
            executor: 'per-vu-iterations',
            vus: 1000,
            iterations: 70,
            maxDuration: '10m',
        },
    },
};

export default function (ledger) {
    const url = ledger.url + "/tests01"
    ReadTransactions(url);
    ReadAccounts(url);
    ReadBalances(url);
    ReadAggregateBalances(url);
}

export function teardown(data) {
    extension.stopLedger();
}
