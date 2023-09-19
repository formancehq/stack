const _ = require("lodash");
const Validator = require("../lib/validator.js");

let schema = {
  $id: "1",
  title: "1 unresolved",
  type: "object",
  properties: {
    property1: {
      type: ["string", "number"],
    },
    additionalProperties: false,
  },
  additionalProperties: false,
};

let dataValid = {
  property1: "bar",
};

let dataInvalidAdditional = {
  foo: "bar",
};

let dataInvalidType = {
  property1: null,
};

const defaultOptions = { allowUnionTypes: true };

test("validates schemas and data", () => {
  const validator = new Validator([schema], defaultOptions);
  const result = validator.validateSchemaAndData(schema, dataValid);
  expect(result).toBe(true);
});

test("fails with additional properties", () => {
  expect.assertions(2);
  const validator = new Validator([schema], defaultOptions);
  let result;
  try {
    result = validator.validateSchemaAndData(schema, dataInvalidAdditional);
  } catch (e) {
    expect(e.length).toBe(1);
    expect(e[0].keyword).toBe("additionalProperties");
  }
});

test("fails with invalid type", () => {
  expect.assertions(3);
  const validator = new Validator([schema], defaultOptions);
  let result;
  try {
    result = validator.validateSchemaAndData(schema, dataInvalidType);
  } catch (e) {
    expect(e.length).toBe(1);
    expect(e[0].keyword).toBe("type");
    expect(e[0].schemaPath).toBe("#/properties/property1/type");
  }
});
