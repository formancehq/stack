const _ = require("lodash");
const pointer = require("jsonpointer");

const getFilename = (path) => {
  let parts = path.split("/");
  return parts[parts.length - 1];
};

const resolveRef = (
  unresolvedSchema,
  mergedSchema,
  unresolvedProperty,
  mergedProperty
) => {
  let ref;
  if (unresolvedProperty.type === "array") {
    ref = unresolvedProperty.items?.$ref;
  } else {
    ref = unresolvedProperty.$ref;
  }

  if (!ref) {
    return;
  } else if (ref.substr(0, 1) !== "#") {
    mergedProperty.$ref = ref;
    return;
  }

  let refPointer = ref.substr(1);
  let resolvedRef;

  try {
    resolvedRef = pointer.get(unresolvedSchema, refPointer);
  } catch (e) {
    console.error(
      "Error resolving JSON pointer",
      refPointer,
      "in property",
      property.title,
      "of schema",
      unresolvedSchema.$id
    );
    console.error(e);
  }

  // @todo technically this could also require resolution

  if (resolvedRef) {
    if (resolvedRef.type === "array") {
      mergedProperty.items.$ref = resolvedRef.items.$ref || ref;
    } else {
      mergedProperty.$ref = resolvedRef.$ref || ref;
    }

    // @todo handle allOf, anyOf
    if (resolvedRef.oneOf) {
      resolvedRef.oneOf.forEach((oneItem, index) => {
        if (oneItem.$ref) {
          let mergedResolvedRef = pointer.get(mergedSchema, refPointer);
          if (mergedResolvedRef.oneOf && mergedResolvedRef.oneOf[index]) {
            mergedResolvedRef.oneOf[index].$ref = oneItem.$ref;
          }
        }
      });
    }
  }

  return;
};

const mergeProperty = (unresolvedSchema, mergedSchema, key) => {
  let unresolvedProperty;
  let mergedProperty;

  if (mergedSchema.properties) {
    mergedProperty = mergedSchema.properties[key];
  } else {
    return;
  }
  if (unresolvedSchema && unresolvedSchema.properties) {
    unresolvedProperty = unresolvedSchema.properties[key];
  }

  if (mergedProperty === undefined) {
    return;
  }

  mergedProperty.isRequired =
    Array.isArray(mergedSchema.required) && mergedSchema.required.includes(key);

  if (unresolvedProperty) {
    resolveRef(
      unresolvedSchema,
      mergedSchema,
      unresolvedProperty,
      mergedProperty
    );
  }

  if (mergedProperty.properties) {
    for (const key in mergedProperty.properties) {
      mergeProperty(unresolvedProperty, mergedProperty, key);
    }
  }
};

var Merger = function () {};

// resolvedSchemas do not have any $ref data
// add the $ref from the related unresolvedSchemas into each resolvedSchema
Merger.mergeSchemas = function (unresolvedSchemas, resolvedSchemas) {
  return unresolvedSchemas.map((unresolvedSchema) => {
    let resolvedSchema = resolvedSchemas.find(
      (resolvedSchema) => resolvedSchema.filename === unresolvedSchema.filename
    );

    if (!resolvedSchema) {
      throw "Unable to resolve schema " + unresolvedSchema.filename;
    }

    let mergedSchema = _.cloneDeep(resolvedSchema);
    mergedSchema.cleanSchema = _.cloneDeep(resolvedSchema.data);

    if (mergedSchema.data.properties) {
      for (const key in mergedSchema.data.properties) {
        mergeProperty(unresolvedSchema.data, mergedSchema.data, key);
      }
    }

    // @todo handle top-level oneOf, allOf, anyOf
    if (mergedSchema.data.oneOf) {
      mergedSchema.data.oneOf.forEach((one, index) => {
        // @todo this is going to throw exceptions
        mergedSchema.data.oneOf[index].$ref =
          unresolvedSchema.data.oneOf[index].$ref;
      });
    }

    mergedSchema.schema = mergedSchema.data;
    delete mergedSchema.data;

    mergedSchema.unresolvedSchema = unresolvedSchema.data;

    return mergedSchema;
  });
};

module.exports = Merger;
