const fastGlob = require("fast-glob");
const fs = require("fs");
const { promisify } = require("util");
const YAML = require("yaml");
const readFileAsync = promisify(fs.readFile);

const loadFiles = async (files) => {
  return await Promise.all(
    files.map(async (file) => {
      let dataObject;

      let fileContent = await readFileAsync(file).catch((e) => {
        console.error("Error reading file: ", file);
        console.error(e);
      });

      fileContent = fileContent.toString();
      const extension = file.split(".").pop();

      try {
        if (extension === "yml" || extension == "yaml") {
          dataObject = YAML.parse(fileContent);
        } else {
          dataObject = JSON.parse(fileContent);
        }
      } catch (e) {
        console.error("Error parsing file: ", file);
        console.error(e);
      }

      return {
        filename: file,
        data: dataObject,
      };
    })
  );
};

var Loader = function () {};

Loader.loadFiles = async function (glob) {
  const files = fastGlob.sync(glob);
  const results = await loadFiles(files);
  return results.filter((result) => result.data !== undefined);
};

module.exports = Loader;
