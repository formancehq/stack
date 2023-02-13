import {Client as Client1_3, Config, Headers} from "../1_3";
import {loadConfig} from "../../config";

export class Client extends Client1_3 {

    version = '1.4';

    constructor(
        config: Config
    ) {
        super(config);
    }
}

export const newClient = () => {
    const config = loadConfig();
    const headers: Headers = {};
    return new Client({
        endpoint: config.ledgerUrl,
        headers
    });
};
