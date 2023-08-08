
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
                    { duration: '10s', target: 0 },
                    { duration: '10m', target: 0 },
                ],
                gracefulRampDown: '0s',
                exec: 'write',
                tags: { testid: __ENV.TEST_ID}
            },
            read_constant: {
                executor: 'ramping-vus',
                startVUs: 0,
                stages: [
                    { duration: '10m', target: 0 },
                    { duration: '10s', target: 0 },
                    { duration: '10s', target: 20 },
                    { duration: '5m', target: 20 },
                ],
                gracefulRampDown: '0s',
                exec: 'read',
                tags: { testid: __ENV.TEST_ID}
            },
        }
    };
}