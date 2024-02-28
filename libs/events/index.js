const fs = require("fs/promises");
const yaml = require('yaml');

(async () => {

    const rawBase = await fs.readFile("./base.yaml", { encoding: 'utf8' });
    const base = yaml.parse(rawBase);

    for(const service of await fs.readdir("services")) {
        for(const version of await fs.readdir('services/' + service)) {
            for(const event of await fs.readdir('services/' + service + '/' + version)) {
                const rawEventData = await fs.readFile('services/' + service + '/' + version + '/' + event, { encoding: 'utf8' });
                base.properties.payload = yaml.parse(rawEventData);
                const directory = 'generated/' + service + '/' + version + '/';
                await fs.mkdir(directory, { recursive: true });
                await fs.writeFile(directory + event.replace('.yaml', '.json'), JSON.stringify(base, null, 2));
            }
        }
    }
})();