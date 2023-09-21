const fastGlob = require("fast-glob");
const $RefParser = require("@apidevtools/json-schema-ref-parser");

var Resolver = function () {};

Resolver.resolveSchemas = async function (fileGlob, resolvers) {
  const files = fastGlob.sync(fileGlob);
  const results = await Promise.all(
    files.map(async (filename) => {
      let resolvedSchema = await Resolver.resolveSchema(
        filename,
        resolvers
      ).catch((e) => {
        console.error("Error resolving", filename);
        console.error(e.message);
      });
      return {
        filename: filename,
        data: resolvedSchema,
      };
    })
  );
  return results.filter((results) => results.data != undefined);
};

Resolver.resolveSchema = async function (schemaToResolve, resolvers) {
  return await $RefParser.dereference(schemaToResolve, { resolve: resolvers });
};

module.exports = Resolver;
