"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const ledger_nodejs_1 = require("@numaryhq/ledger-nodejs");
const main = () => {
    const configuration = (0, ledger_nodejs_1.createConfiguration)();
    const serverApi = new ledger_nodejs_1.ServerApi(configuration);
    // TODO: Add some testing
};
main();
