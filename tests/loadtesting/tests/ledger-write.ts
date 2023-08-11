// @ts-ignore
import {Options} from 'k6/options';
import testInfo from "./../src/ledger/infos-test";
import testWriteTransactions from "./../src/ledger/write-transactions-test";

export let options: Options = {
    vus: 20,
    duration: '5m',
};

export default () => {
    testInfo();
    testWriteTransactions();
};
