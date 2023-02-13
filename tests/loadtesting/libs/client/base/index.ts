import http from "k6/http";
import {parseVersion, Version} from "../../core";
import {loadConfig} from "../../config";

export interface BaseClient {
    readonly version: string;
}

export const discoverServerVersion = (): Version => {
    const config = loadConfig();
    const ret = http.get(config.ledgerUrl + '/_info');

    return parseVersion(ret.json("data.version") as string);
};

export const fromServerVersion = (version: string, callback: () => void) => {
    const serverVersion = discoverServerVersion();
    if (serverVersion.lt(parseVersion(version))) {
        console.info(`Skip test, support is after version: ${version}`);
        return;
    }
    callback();
};
