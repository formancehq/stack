const Loader = require("../lib/loader.js");

test("loads single schema", async () => {
  const results = await Loader.loadFiles("./tests/examples/schema/person.json");
  expect(results).toHaveLength(1);
});

test("loads multiple schema", async () => {
  const results = await Loader.loadFiles("./tests/examples/schema/**.json");
  expect(results).toHaveLength(2);
});

test("gracefully handles malformed schema", async () => {
  const results = await Loader.loadFiles(
    "./tests/examples/schema-with-errors/malformed.json"
  );
  expect(results).toHaveLength(0);
});
