// @ts-ignore
import {Options} from 'k6/options';
import testInfo from "./../src/ledger/infos-test";
import testStats from "./../src/ledger/stats-test";
import testAccounts from "./../src/ledger/accounts-test";
import testTransactions from "./../src/ledger/transactions-test";
import testWriteTransactions from "./../src/ledger/write-transactions-test";

export const projectId = 3590983;

export let options: Options = {
    vus: 1,
    iterations: 1,
    duration: '1m',
    throw: true,
    ext: {
        loadimpact: {
            name: 'Performance tests',
            projectID: projectId,
            distribution: {
                "amazon:fr:paris": { loadZone: "amazon:fr:paris", percent: 100 },
            },
        },
    },
    thresholds: {
        http_req_failed: ['rate<0.01'], // http errors should be less than 1%
        http_req_duration: ['p(95)<200'], // 95% of requests should be below 200ms
        http_reqs: ['count>=20'], // at least 20 requests should be made
    },
};

export default () => {
    testInfo();
    testStats();
    testAccounts();
    testTransactions();
    testWriteTransactions();
};
