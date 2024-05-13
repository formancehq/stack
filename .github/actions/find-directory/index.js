const fs = require('fs');
const path = require('path');

function findDockerFile(dir) {
    let results = [];
    const list = fs.readdirSync(dir);
    list.forEach(file => {
        file = path.resolve(dir, file);
        const stat = fs.statSync(file);
        if (stat && stat.isDirectory()) {
            /* if it is a directory, recurse */
            results = results.concat(findDockerFile(file));
        } else {
            if (path.basename(file) === ".goreleaser.yml") {
                results.push(path.relative(".", path.dirname(file)));
            }
        }
    });
    return results;
}

const dataComponents = findDockerFile("./components");
const dataEe = findDockerFile("./ee");
const data = dataComponents.concat(dataEe);
console.log(JSON.stringify(data, null, 0));