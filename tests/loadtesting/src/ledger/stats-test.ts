import {withClient} from "../../libs/client";

export default () => {
    withClient('1.3', client => {
        client.getStats(__ENV.LEDGER);
    });
};
