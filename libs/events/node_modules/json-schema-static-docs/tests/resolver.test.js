const Resolver = require("../lib/resolver.js");

test("resolves single schema", async () => {
  const result = await Resolver.resolveSchema(
    "./tests/examples/schema/person.json"
  );
  expect(result.properties.name.$ref).toBeUndefined();
  expect(result.properties.name.$id).toBe(
    "http://example.com/schemas/name.json"
  );
});

test("resolves multiple schemas", async () => {
  const results = await Resolver.resolveSchemas(
    "./tests/examples/schema/**.json"
  );
  expect(results).toHaveLength(2);
  expect(results[0].data.$id).toBe("http://example.com/schemas/name.json");
  expect(results[1].data.$id).toBe("http://example.com/schemas/person.json");
  expect(results[1].data.properties.name.$id).toBe(
    "http://example.com/schemas/name.json"
  );
  expect(results[1].data.properties.name.$ref).toBeUndefined();
});
