const fs = require("fs");
const path = require("path");
const rimraf = require("rimraf");
const _ = require("lodash");

const JsonSchamaStaticDocs = require("../lib/json-schema-static-docs.js");

let defaultTestOptions = {
  inputPath: "./tests/examples/schema/",
  outputPath: "./tests/docs/",
};
let testOptions = {};

beforeEach(() => {
  testOptions = _.cloneDeep(defaultTestOptions);
  if (fs.existsSync(testOptions.outputPath)) {
    rimraf.sync(testOptions.outputPath);
  }
  fs.mkdirSync(testOptions.outputPath);
});

afterEach(() => {});

test("resolves single schema", async () => {
  const jsonSchameStaticDocs = new JsonSchamaStaticDocs(testOptions);
  const result = await jsonSchameStaticDocs.generate();
  expect(result.length).toBe(3);
});

test("handles absolute paths", async () => {
  testOptions.inputPath = path.resolve(__dirname, "./examples/schema/");
  testOptions.outputPath = path.resolve(__dirname, "./docs/");
  const jsonSchameStaticDocs = new JsonSchamaStaticDocs(testOptions);
  let mergedSchemas = await jsonSchameStaticDocs.generate();
  expect(mergedSchemas[2].relativeFilename).toBe("sub/different-person.json");
  const exists = fs.existsSync("./tests/docs/sub/different-person.md");
  expect(exists).toBe(true);
});

test("supports custom templates", async () => {
  testOptions.templatePath = path.join(__dirname, "examples/templates/");
  const jsonSchameStaticDocs = new JsonSchamaStaticDocs(testOptions);
  await jsonSchameStaticDocs.generate();
  let result = fs.readFileSync(path.join(testOptions.outputPath, "name.md"));
  expect(result.toString()).toBe("foo");
});

test("allows templates to be skipped", async () => {
  testOptions.skipTemplates = true;
  const jsonSchameStaticDocs = new JsonSchamaStaticDocs(testOptions);
  await jsonSchameStaticDocs.generate();
  let result = fs.readFileSync(
    path.join(testOptions.outputPath, "person.json")
  );
  let schema = JSON.parse(result.toString());
  expect(schema.title).toBe("Person");
  expect(schema.properties.name.properties.firstNames.type).toBe("string");
});
