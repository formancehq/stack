import {Client as Client1_7, newClient as newClient1_7} from "./1_7";
import {Client as Client1_5, newClient as newClient1_5} from "./1_5";
import {Client as Client1_4, newClient as newClient1_4} from "./1_4";
import {Client as Client1_3, newClient as newClient1_3} from "./1_3";

export type ClientMapping = {
    '1.3': Client1_3,
    '1.4': Client1_4,
    '1.5': Client1_5,
    '1.7': Client1_7,
};

export const withClient = <V extends keyof ClientMapping>(version: V, callback: (client: ClientMapping[V]) => any) => {
    let client: ClientMapping[V];
    switch(version) {
        case "1.3":
            client = newClient1_3();
            break;
        case "1.4":
            client = newClient1_4();
            break;
        case "1.5":
            client = newClient1_5();
            break;
        case "1.7":
            client = newClient1_7();
            break;
        default:
            throw new Error("Invalid ledger version");
    }
    return callback(client);
};
