const fs = require("fs/promises");
const path = require("path");
const yaml = require('yaml');

const VERSION_BASE_FILE = 'base.yaml';

async function sortedDirectories(directory) {
    return (await fs.readdir(directory, { withFileTypes: true }))
        .filter((entry) => entry.isDirectory())
        .map((entry) => entry.name)
        .sort();
}

async function sortedEventFiles(directory) {
    return (await fs.readdir(directory, { withFileTypes: true }))
        .filter((entry) => entry.isFile() && entry.name.endsWith('.yaml') && entry.name !== VERSION_BASE_FILE)
        .map((entry) => entry.name)
        .sort();
}

async function readVersionBase(versionDirectory) {
    try {
        return await fs.readFile(path.join(versionDirectory, VERSION_BASE_FILE), { encoding: 'utf8' });
    } catch (error) {
        if (error.code === 'ENOENT') {
            return null;
        }
        throw error;
    }
}

function composeSchema(rawDefaultBase, rawVersionBase, rawEventData) {
    const event = yaml.parse(rawEventData);
    if (rawVersionBase !== null) {
        const base = yaml.parse(rawVersionBase);
        base.allOf = [...(base.allOf || []), event];
        return base;
    }

    const base = yaml.parse(rawDefaultBase);
    base.properties.payload = event;
    return base;
}

(async () => {

    const rawBase = await fs.readFile(path.join(__dirname, "base.yaml"), { encoding: 'utf8' });
    const aggregated = {};

    for(const service of await sortedDirectories(path.join(__dirname, "services"))) {
        aggregated[service] = {};
        for(const version of await sortedDirectories(path.join(__dirname, 'services', service))) {
            aggregated[service][version] = {};
            const versionDirectory = path.join(__dirname, 'services', service, version);
            const rawVersionBase = await readVersionBase(versionDirectory);
            for(const event of await sortedEventFiles(versionDirectory)) {
                const rawEventData = await fs.readFile(path.join(versionDirectory, event), { encoding: 'utf8' });
                const schema = composeSchema(rawBase, rawVersionBase, rawEventData);
                const directory = path.join(__dirname, 'generated', service, version);
                await fs.mkdir(directory, { recursive: true });
                await fs.writeFile(path.join(directory, event.replace('.yaml', '.json')), JSON.stringify(schema, null, 2));

                aggregated[service][version][event.replace('.yaml', '')] = schema;
            }
        }
    }

    const aggregatedJSON = JSON.stringify(aggregated, null, 2);
    await fs.writeFile(path.join(__dirname, 'generated', 'all.json'), aggregatedJSON);

    // Keep the legacy location in sync while existing consumers (notably the
    // Platform UI) still fetch the event catalogue from libs/events.
    const legacyGeneratedDirectory = path.join(__dirname, '..', 'libs', 'events', 'generated');
    await fs.mkdir(legacyGeneratedDirectory, { recursive: true });
    await fs.writeFile(path.join(legacyGeneratedDirectory, 'all.json'), aggregatedJSON);
})();
