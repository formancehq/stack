---
title: JSON Schema Static Docs
---

# Json Schema Static Docs

[![npm version](https://badge.fury.io/js/json-schema-static-docs.svg)](https://badge.fury.io/js/json-schema-static-docs) [![CircleCI](https://circleci.com/gh/tomcollins/json-schema-static-docs/tree/master.svg?style=svg)](https://circleci.com/gh/tomcollins/json-schema-static-docs/tree/master)

This library generates human friendly static markdown documents based on a set of JSON schema documents.

## Features

- quickly generate static documentation with optional frontmatter
- generates an index of documents
- able to display nested properties and objects
- describes array item schema
- displays examples for schema and properties
- generates links between documents when schema include relative $ref values
- support for descriptions when using string enum values via a custom keyword

See [this post](https://careers.dft.gov.uk/dvla-software-developers-behind-the-screens/) describing how this library is used by the [DVLA](https://github.com/dvla/).

## Examples

{% include examples-table.md %}

See the [examples page](/json-schema-static-docs/examples/) for more info.

## Support for JSON schema specification versions

Currently supports schema specified using the following [specification versions](https://json-schema.org/specification-links.html):
draft-06, draft-07 and draft-2019-09.

Examples for each specification version can be found below.

You can view a [detailed description of supported keywords](/json-schema-static-docs/support/).

## Installation

```bash
npm install json-schema-static-docs
```

## Usage

```javascript
const JsonSchemaStaticDocs = require("json-schema-static-docs");

(async () => {
  let jsonSchemaStaticDocs = new JsonSchemaStaticDocs({
    inputPath: "schema",
    outputPath: "docs",
    ajvOptions: {
      allowUnionTypes: true,
    },
  });
  await jsonSchemaStaticDocs.generate();
  console.log("Documents generated.");
})();
```

## Options

| Parameter      | Description                                   | Default            |
| -------------- | --------------------------------------------- | ------------------ |
| inputPath      | Directory containing your schema              | "schema"           |
| inputFileGlob  | Glob used to look for schema files            | "\*_/_.{yml,json}" |
| outputPath     | Where to write documentation files            | "docs"             |
| createIndex    | Create an index of documents                  | true               |
| indexPath      | Index file path (relative to outputPath)      | "index.md"         |
| indexTitle     | Title of the generated index page             | "Index"            |
| templatePath   | Where to find templates                       | "templates"        |
| ajvOptions     | Options to pass to [AJV](https://ajv.js.org/) | {}                 |
| enableMetaEnum | Allow documentation of enum values            | false              |
| addFrontMatter | Add front matter to generated documentation   | false              |
| displaySchema  | Display schema JSON in output                 | true               |

## Index Creation

By default a root level index will be created in the specified `outputPath`.

You can see an example of the [here](examples/examples-index.html);

### Customising the index

You can specify a path and title for the index.

```javascript
let jsonSchemaStaticDocs = new JsonSchemaStaticDocs({
  inputPath: "schema",
  outputPath: "docs",
  indexPath: "schema-index.md",
  indexTitle: "List of schema with custom title",
});
```

## Markdown Front Matter

If you want to include markdown front matter (for Jekyll, Hugo etc) use the `addFrontMatter` options.

```javascript
let jsonSchemaStaticDocs = new JsonSchemaStaticDocs({
  inputPath: "schema",
  outputPath: "docs",
  addFrontMatter: true,
});
await jsonSchemaStaticDocs.generate();
```

This will prepend generated markdown documents with the schema title or it.

```yml
---
title: The schema title or $id
---
# documentation starts here
```

## Describing Enums

Json-schema allows a set of enumeration values to be defined for a string property but does not allow descriptions to be defined for each value. Descriptions within documentation can be very useful.

This library supports the `meta:enum` convention (inspired by [adobe/jsonschema2md](https://github.com/adobe/jsonschema2md) to allow descriptions to be defined for enums.

You will need to enable this feature using the `enableMetaEnum` option:

```javascript
let jsonSchemaStaticDocs = new JsonSchemaStaticDocs({
  inputPath: "schema",
  outputPath: "docs",
  enableMetaEnum: true,
});
await jsonSchemaStaticDocs.generate();
```

_This allows the `meta:enum` keyword to be used when applying strict validation._

And then define the `meta:enum` descriptions adjacent to your `enum` e.g.

```yml
property1:
  title: "Property 1"
  type: "string"
  enum: ["foo", 42]
  meta:enum:
    foo: Description for foo
    42: Description for 42
```

## Custom Templates

Templates are authored in [handlebars.js](https://handlebarsjs.com).

The default template is [templates/markdown/schema.hbs](https://github.com/tomcollins/json-schema-static-docs/blob/master/templates/markdown/schema.hbs).

You can provide your own custom templates using the `templatePath` option.

In the example below you will be expected to provide `./your-templates/schema.hbs'.

```javascript
const JsonSchemaStaticDocs = require("json-schema-static-docs");

(async () => {
  let jsonSchemaStaticDocs = new JsonSchemaStaticDocs({
    inputPath: "schema",
    outputPath: "docs",
    templatePath: "your-templates/",
  });
  await jsonSchemaStaticDocs.generate();
  console.log("Documents generated.");
})();
```

_There are currently limitations when using custom templates. Some elements are rendered through handlebars helpers outside the templates._
