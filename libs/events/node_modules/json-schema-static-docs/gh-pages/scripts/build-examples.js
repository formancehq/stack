const JsonSchemaStaticDocs = require("../../lib/json-schema-static-docs");
const path = require("path");
const fastGlob = require("fast-glob");

const ymlPath = path.resolve(__dirname, "../yml/");
const examplesPath = path.resolve(__dirname, "../examples/");

(async () => {
  let jsonSchemaStaticDocs = new JsonSchemaStaticDocs({
    inputPath: ymlPath,
    outputPath: examplesPath,
    indexPath: "examples-index.md",
    indexTitle: "List of Examples",
    ajvOptions: {
      allowUnionTypes: true,
    },
    enableMetaEnum: true,
    addFrontMatter: true,
  });
  await jsonSchemaStaticDocs.generate();
  console.log("Documents generated.");
})();
