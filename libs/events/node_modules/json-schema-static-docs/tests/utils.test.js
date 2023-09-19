const { determineSchemaRelativePath } = require("../lib/utils.js");

const path = require("path");
const temp = path.join("./test/", "*/**.json");
console.log("temp", temp);

describe("Utils", () => {
  describe("determineSchemaRelativePath", () => {
    test("handles absolute schemaInputPath", () => {
      const outputFilename = determineSchemaRelativePath(
        "/Users/test/schema/foo/bar.json",
        "/Users/test/schema"
      );
      expect(outputFilename).toBe("foo/bar.json");
    });

    test("handles absolute schemaInputPath with trailing slash", () => {
      const outputFilename = determineSchemaRelativePath(
        "/Users/test/schema/foo/bar.json",
        "/Users/test/schema/"
      );
      expect(outputFilename).toBe("foo/bar.json");
    });

    test("handles relative schemaInputPath", () => {
      const outputFilename = determineSchemaRelativePath(
        "tests/examples/schema/foo/bar.json",
        "tests/examples/schema"
      );
      expect(outputFilename).toBe("foo/bar.json");
    });

    test("handles relative schemaInputPath with trailing slash", () => {
      const outputFilename = determineSchemaRelativePath(
        "tests/examples/schema/foo/bar.json",
        "tests/examples/schema/"
      );
      expect(outputFilename).toBe("foo/bar.json");
    });
  });
});
