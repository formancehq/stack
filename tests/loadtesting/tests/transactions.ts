import {Options} from "k6/options";
import http from "k6/http";
import exec from 'k6/execution';

export const projectId = 3590983;

export let options: Options = {
    vus: 100,
    iterations: 10000,
    duration: '10m',
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
        http_req_duration: ['p(95)<500'], // 95% of requests should be below 200ms
        http_reqs: ['count>=20'], // at least 20 requests should be made
    },
};

export default () => {
    const url = `${__ENV.LEDGER_URL}/default/transactions`
    http.post(url, JSON.stringify({
        "postings": [{
            "source": "world",
            "destination": `account:${exec.scenario.iterationInTest%100}`,
            "amount": 100,
            "asset": "USD/2",
        }]
    }))
}
