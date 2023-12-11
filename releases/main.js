import { exec } from "child_process";
import {Octokit} from "@octokit/core";
import {paginateRest} from "@octokit/plugin-paginate-rest";
import semver from "semver";
import YAML from 'yaml'
import {Command} from "commander";
import fs from "fs";

const getVersions = async (string) => {
    const data = [];

    const MyOctokit = Octokit.plugin(paginateRest);
    const octokit = new MyOctokit({
        auth: process.env.GITHUB_TOKEN,
    });

    const stackResponse = await octokit.paginate("GET /repos/{owner}/{repo}/releases", {
        owner: "formancehq",
        repo: "stack",
        per_page: 100,
    });

    stackResponse.forEach((release) => {
        if (release.tag_name.includes('components/'))
        {
            const name = release.tag_name.match(/([a-z].*)\/([a-z].*)\/([a-z].*)/)[2];
            if (name !== "ledger" || name !== "fctl" || name !== "operator" || name !== "agent")
            {
                data.push(
                    release.tag_name.replace('components/', '')
                );
            }
        }
        if (release.tag_name.includes('ee/'))
        {
            const name = release.tag_name.match(/([a-z].*)\/([a-z].*)\/([a-z].*)/)[2];
            if (name !== "ledger" || name !== "fctl" || name !== "operator" || name !== "agent")
            {
                data.push(
                    release.tag_name.replace('ee/', '')
                );
            }
        }
    });

    const ledgerResponse = await octokit.paginate("GET /repos/{owner}/{repo}/releases", {
        owner: "formancehq",
        repo: "ledger",
        per_page: 100,
    });

    ledgerResponse.forEach((release) => {
        if (semver.satisfies(release.tag_name, '1.x || <= 2.0.0'))
        {
            data.push("ledger/"+release.tag_name);
        }
    });

    return data;
}

const generateVersions = async () => {
    const releases = await getVersions();
    const data = {};
    releases.forEach((release) => {
        const regex = release.match(/(.*)\/(.*)/)
        const service = regex[1];
        const version = regex[2];
        if (service === "fctl" || service === "operator" || service === "agent")
        {
            return;
        }
        if (semver.prerelease(version) !== null)
        {
            return;
        }

        if (data[service] === undefined || semver.lt(data[service], version)){
            data[service] = version;
        }
    });
    return data;
}

const versionMajor = (version) => {
    // Extract version for folder name
    return `v${semver.major(version)}.${semver.minor(version)}`;
}

const program = new Command();
program.command('create')
    .description('Create a new version file')
    .argument('<version>', 'Version to create')
    .action(async (str, options) => {
        const version = str;

        // Check if version is semver
        if (!semver.valid(version)) {
            console.error("Version is not semver");
            process.exit(1);
        }

        const components = await generateVersions();
        const jsonObject = {
            "version": version,
            "components": components
        }
        const doc = new YAML.Document();
        doc.contents = jsonObject;

        // Extract version for folder name
        const versionFolder = versionMajor(version);

        // Write to file in versions folder
        const dir = `./versions/${versionFolder}`;
        if (!fs.existsSync(dir)){
            fs.mkdirSync(dir);
        }
        fs.writeFileSync(`${dir}/main.yaml`, doc.toString())
        console.log(`Created version ${version}`)
    });

program.command('generate')
    .description('Generate OPENAPI')
    .argument('<version>', 'Version to create')
    .action(async (str, options) => {
        const version = str;

        // Create directory if not exist build
        if (!fs.existsSync('./build')){
            fs.mkdirSync('./build');
        }

        // Extract version for folder name
        const versionFolder = versionMajor(version);
        const dirVersion = `./versions/${versionFolder}`;

        const versionFile = fs.readFileSync(`${dirVersion}/main.yaml`, 'utf8')
        const CONTENT = YAML.parse(versionFile)

        // Generate openapi-merge.json
        const openapiFile = fs.readFileSync(`./templates/openapi/openapi-merge.json`, 'utf8')
        const openapiConfig = openapiFile
            .replace('AUTH_VERSION', CONTENT.components.auth)
            .replace('LEDGER_VERSION', CONTENT.components.ledger)
            .replace('PAYMENTS_VERSION', CONTENT.components.payments)
            .replace('SEARCH_VERSION', CONTENT.components.search)
            .replace('ORCHESTRATION_VERSION', CONTENT.components.orchestration)
            .replace('WALLETS_VERSION', CONTENT.components.wallets)
            .replace('WEBHOOKS_VERSION', CONTENT.components.webhooks)
            .replace('STARGATE_VERSION', CONTENT.components.stargate)
            .replace('GATEWAY_VERSION', CONTENT.components.gateway)
            .replace('OUTPUT_FILE', `./../${dirVersion}/openapi.json`)
        fs.writeFileSync(`./build/openapi-merge.json`, openapiConfig)

        // Generate base.yaml
        const baseFile = fs.readFileSync(`./templates/openapi/base.yaml`, 'utf8')
        const baseConfig = baseFile
            .replace('SDK_VERSION', CONTENT.version)
        fs.writeFileSync(`./build/base.yaml`, baseConfig)

        // Run command for generate openapi.json
        // openapi-merge-cli --config ./build/openapi-merge.json
        exec(`npm run build`, (error, stdout, stderr) => {
            if (error) {
                console.log(`error: ${error.message}`);
                process.exit(1);
            }
        });

        console.log(`Generated version ${version}`)
    });

program.command('operator')
    .description('Create a new version file for Operator')
    .argument('<version>', 'Version to create')
    .action(async (str, options) => {
        const version = str;

        // Extract version for folder name
        const versionFile = fs.readFileSync(`./versions/${version}/main.yaml`, 'utf8')
        const CONTENT = YAML.parse(versionFile)

        const jsonObject = {
            "apiVersion": "stack.formance.com/v1beta3",
            "kind": "Versions",
            "metadata": {
                "name": CONTENT.version
            },
            "spec": CONTENT.components
        }
        const doc = new YAML.Document();
        doc.contents = jsonObject;

        // Write to file in versions folder
        const dir = `./versions/${version}`;
        if (!fs.existsSync(dir)){
            fs.mkdirSync(dir);
        }
        fs.writeFileSync(`${dir}/operator.yaml`, doc.toString())
        console.log(`Created version ${version}`)

    });

program.parse();
