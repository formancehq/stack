const fs = require('fs');
const makeDir = require('make-dir');

var Writer = function(){}

Writer.writeFile = async function(filename, data) {
  let parts = filename.split('/');
  let dirName = parts.splice(0, parts.length-1).join('/');
  if (!fs.existsSync(dirName)) {
    await makeDir(dirName);
  }
  fs.writeFileSync(filename, data);
};

module.exports = Writer;