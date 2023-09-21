const {
  getHtmlAnchorForRef,
  getLabelForProperty,
} = require("../../lib/renderer/helpers.js");

describe("renderer helpers", () => {
  describe("getHtmlAnchorForRef", () => {
    describe("handles $ref only to relative file", () => {
      const result = getHtmlAnchorForRef("../foo.json");
      expect(result).toBe('<a href="../foo.html">../foo.html</a>');
    });
    describe("handles $ref only to relative file with an anchor", () => {
      const result = getHtmlAnchorForRef("../foo.json#/definitions/bar");
      expect(result).toBe('<a href="../foo.html#bar">../foo.html#bar</a>');
    });
    describe("handles $ref only to relative file with an anchor converting to lower case", () => {
      const result = getHtmlAnchorForRef("../foo.json#/definitions/barHumbug");
      expect(result).toBe(
        '<a href="../foo.html#barhumbug">../foo.html#barhumbug</a>'
      );
    });
    describe("handles $ref only to relative file at the same level", () => {
      const result = getHtmlAnchorForRef("../foo.json");
      expect(result).toBe('<a href="../foo.html">../foo.html</a>');
    });
    describe("handles $ref only to relative file at the same level with an anchor", () => {
      const result = getHtmlAnchorForRef("../foo.json#/definitions/bar");
      expect(result).toBe('<a href="../foo.html#bar">../foo.html#bar</a>');
    });
  });
  describe("propertyToTypeLabel", () => {
    describe("simple property", () => {
      test("returns expected result for a string property", async () => {
        const property = {
          type: "string",
        };
        const result = getLabelForProperty(property);
        expect(result).toBe("String");
      });
      test("returns expected result for a number property", async () => {
        const property = {
          type: "number",
        };
        const result = getLabelForProperty(property);
        expect(result).toBe("Number");
      });
    });
    describe("simple property with $ref", () => {
      test("returns expected result for a string with a local $ref", async () => {
        const property = {
          type: "string",
          $ref: "#/definitions/foo",
        };
        const result = getLabelForProperty(property);
        expect(result).toBe("String");
      });

      test("returns expected result for a string with a relative $ref", async () => {
        const property = {
          type: "string",
          $ref: "../foo.json#/definitions/bar",
        };
        const result = getLabelForProperty(property);
        expect(result).toBe("String");
      });

      test("returns expected result for a string with a http $ref", async () => {
        const property = {
          type: "string",
          $ref: "http://example.com/foo.json#/definitions/bar",
        };
        const result = getLabelForProperty(property);
        expect(result).toBe("String");
      });
    });

    describe("object property", () => {
      test("returns expected result for an object property", async () => {
        const property = {
          type: "object",
        };
        const result = getLabelForProperty(property);
        expect(result).toBe("Object");
      });

      test("returns expected result for a object with a local $ref", async () => {
        const property = {
          type: "object",
          $ref: "#/definitions/foo",
        };
        const result = getLabelForProperty(property);
        expect(result).toBe("Object");
      });

      test("returns expected result for a object with a local $ref and a title", async () => {
        const property = {
          type: "object",
          $ref: "#/definitions/foo",
          title: "Foo",
        };
        const result = getLabelForProperty(property);
        expect(result).toBe("Object");
      });

      test("returns expected result for a object with a relative $ref", async () => {
        const property = {
          type: "object",
          $ref: "../foo.json#/definitions/bar",
        };
        const result = getLabelForProperty(property);
        expect(result).toBe(
          'Object (of type <a href="../foo.html#bar">../foo.html#bar</a>)'
        );
      });
      test("returns expected result for a object with a relative $ref at the same directory level", async () => {
        const property = {
          type: "object",
          $ref: "foo.json",
        };
        const result = getLabelForProperty(property);
        expect(result).toBe('Object (of type <a href="foo.html">foo.html</a>)');
      });
      test("returns expected result for a object with a relative $ref at the same directory level with an anchor", async () => {
        const property = {
          type: "object",
          $ref: "foo.json#/definitions/bar",
        };
        const result = getLabelForProperty(property);
        expect(result).toBe(
          'Object (of type <a href="foo.html#bar">foo.html#bar</a>)'
        );
      });
      test("returns expected result for a object with a relative $ref and a title", async () => {
        const property = {
          type: "object",
          $ref: "../foo.json#/definitions/bar",
          title: "Foo",
        };
        const result = getLabelForProperty(property);
        expect(result).toBe(
          'Object (of type <a href="../foo.html#bar">Foo</a>)'
        );
      });

      test("returns expected result for a object with a http $ref", async () => {
        const property = {
          type: "object",
          $ref: "http://example.com/foo.json#/definitions/bar",
        };
        const result = getLabelForProperty(property);
        expect(result).toBe("Object");
      });
      test("returns expected result for a object with a http $ref and a title", async () => {
        const property = {
          type: "object",
          $ref: "http://example.com/foo.json#/definitions/bar",
          title: "Foo",
        };
        const result = getLabelForProperty(property);
        expect(result).toBe("Object");
      });
    });
  });
});
