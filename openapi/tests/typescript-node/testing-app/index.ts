import {ServerApi, createConfiguration} from '@formancehq/formance';

const main = () => {
    console.info("Starting test");
    const configuration = createConfiguration();
    const serverApi = new ServerApi(configuration);
    console.info('TODO: Need other checks. Actually just checking we can install the SDK and use it.');
}
main()
