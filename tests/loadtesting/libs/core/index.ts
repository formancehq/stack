import {ClientMapping, withClient} from "../client";
import {fromServerVersion} from "../client/base";

export * from './versions';
export * from './expectations';

export interface MatrixTest {
    serverVersion: string;
    clientVersion: keyof ClientMapping;
}

export const withMatrix = <T extends MatrixTest>(
    tests: T[],
    callback: (t: T, client: ClientMapping[T['clientVersion']]) => void,
) => {
    for (const test of tests) {
        fromServerVersion(test.serverVersion, () => {
            withClient(test.clientVersion, client => {
                callback(test, client);
            });
        });
    }
};
