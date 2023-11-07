import {readAccounts, readAggregatedBalances, readBalances} from "../benchs/ledger-v2";

export function K6Options() {
    return {
        discardResponseBodies: true,
        scenarios: {
            write_constant: {
                executor: 'ramping-vus',
                startVUs: 0,
                stages: [
                    { duration: '10s', target: 20 },
                    { duration: '5m', target: 20 },
                ],
                gracefulRampDown: '0s',
                exec: 'write',
                tags: { testid: __ENV.TEST_ID}
            },
            read_transactions_constant: {
                executor: 'ramping-vus',
                startVUs: 0,
                stages: [
                    { duration: '6m', target: 0 },
                    { duration: '10s', target: 10 },
                    { duration: '5m', target: 10 },
                ],
                gracefulRampDown: '0s',
                exec: 'readTransactions',
                tags: { testid: __ENV.TEST_ID}
            },
            read_accounts_constant: {
                executor: 'ramping-vus',
                startVUs: 0,
                stages: [
                    { duration: '6m', target: 0 },
                    { duration: '10s', target: 10 },
                    { duration: '5m', target: 10 },
                ],
                gracefulRampDown: '0s',
                exec: 'readAccounts',
                tags: { testid: __ENV.TEST_ID}
            },
            read_balances_constant: {
                executor: 'ramping-vus',
                startVUs: 0,
                stages: [
                    { duration: '6m', target: 0 },
                    { duration: '10s', target: 10 },
                    { duration: '5m', target: 10 },
                ],
                gracefulRampDown: '0s',
                exec: 'readBalances',
                tags: { testid: __ENV.TEST_ID}
            },
            read_aggregatebalances_constant: {
                executor: 'ramping-vus',
                startVUs: 0,
                stages: [
                    { duration: '6m', target: 0 },
                    { duration: '10s', target: 10 },
                    { duration: '5m', target: 10 },
                ],
                gracefulRampDown: '0s',
                exec: 'readAggregatedBalances',
                tags: { testid: __ENV.TEST_ID}
            },
        }
    };
}
