const fs = require("fs/promises");
const path = require("path");
const yaml = require('yaml');

(async () => {

    const rawBase = await fs.readFile(path.join(__dirname, "base.yaml"), { encoding: 'utf8' });
    const aggregated = {};

    for(const service of await fs.readdir(path.join(__dirname, "services"))) {
        aggregated[service] = {};
        for(const version of await fs.readdir(path.join(__dirname, 'services', service))) {
            aggregated[service][version] = {};
            for(const event of await fs.readdir(path.join(__dirname, 'services', service, version))) {
                const rawEventData = await fs.readFile(path.join(__dirname, 'services', service, version, event), { encoding: 'utf8' });
                const base = yaml.parse(rawBase);
                base.properties.payload = yaml.parse(rawEventData);
                const directory = path.join(__dirname, 'generated', service, version);
                await fs.mkdir(directory, { recursive: true });
                await fs.writeFile(path.join(directory, event.replace('.yaml', '.json')), JSON.stringify(base, null, 2));

                aggregated[service][version][event.replace('.yaml', '')] = base;
            }
        }
    }

    console.log(aggregated);
    const aggregatedJSON = JSON.stringify(aggregated, null, 2);
    await fs.writeFile(path.join(__dirname, 'generated', 'all.json'), aggregatedJSON);

    // Keep the legacy location in sync while existing consumers (notably the
    // Platform UI) still fetch the event catalogue from libs/events.
    const legacyGeneratedDirectory = path.join(__dirname, '..', 'libs', 'events', 'generated');
    await fs.mkdir(legacyGeneratedDirectory, { recursive: true });
    await fs.writeFile(path.join(legacyGeneratedDirectory, 'all.json'), aggregatedJSON);
})();
