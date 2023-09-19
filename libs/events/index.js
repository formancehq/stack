const fs = require("fs/promises");
const yaml = require('yaml')
const JsonSchemaStaticDocs = require("json-schema-static-docs");

(async () => {

    const rawBase = await fs.readFile("./base.yaml", { encoding: 'utf8' });
    const base = yaml.parse(rawBase);

    for(const service of await fs.readdir("services")) {
        for(const version of await fs.readdir('services/' + service)) {
            for(const event of await fs.readdir('services/' + service + '/' + version)) {
                const rawEventData = await fs.readFile('services/' + service + '/' + version + '/' + event, { encoding: 'utf8' });
                base.properties.payload = yaml.parse(rawEventData);
                fs.writeFile('generated/' + service + '-' + version + '-' + event + '.json', JSON.stringify(base));
            }
        }
    }
})();