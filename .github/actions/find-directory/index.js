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
                const obj = { component: path.basename(path.dirname(file)), type: path.basename(path.dirname(path.dirname(file))) };
                results.push(obj);
            }
        }
    });
    return results;
}

const dataComponents = findDockerFile("./components");
const dataEe = findDockerFile("./ee");
const dataTools = findDockerFile("./tools");
const data = dataComponents.concat(dataEe).concat(dataTools);
console.log(JSON.stringify(data, null, 0));