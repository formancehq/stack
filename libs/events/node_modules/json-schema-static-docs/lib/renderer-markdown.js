const Handlebars = require("handlebars");
const fs = require("fs");
const path = require("path");
const util = require("util");
const _ = require("lodash");
const {
  getHtmlAnchorForRef,
  getLabelForProperty,
  upperCaseFirstCharacter,
} = require("./renderer/helpers");

const readFile = util.promisify(fs.readFile);

/* ------ */

Handlebars.registerHelper("isExternalRef", function (value) {
  return (
    typeof value === "string" &&
    (value.match(/^(.+)#/) || value.match(/^(.+)\.(json|yaml|yml)$/))
  );
});
Handlebars.registerHelper("getHtmlAnchorForRef", function (ref, label) {
  return getHtmlAnchorForRef(ref, label);
});

/* ------ */

Handlebars.registerHelper("json", function (context) {
  return JSON.stringify(context, undefined, 4);
});

Handlebars.registerHelper("enumValues", function (enums) {
  let html = "<ul>";
  if (Array.isArray(enums)) {
    enums.forEach((e) => {
      html += `<li>${String(e)}</li>`;
    });
  }
  html += "</ul>";
  return new Handlebars.SafeString(html);
});

Handlebars.registerHelper("isDefined", function (value) {
  return value !== undefined;
});

Handlebars.registerHelper("isArrayOfStringLikeValues", function (value) {
  return Array.isArray(value) && !value.find((v) => typeof v === "object");
});

const propertyRow = (name, property, parentKey) => {
  let href = "";
  if (typeof parentKey === "string") {
    href = parentKey.toLowerCase();
  }
  href += name.toLowerCase();

  let typeLabel;
  if (property.const) {
    typeLabel = `${upperCaseFirstCharacter(property.type)}=${property.const}`;
  } else {
    typeLabel = getLabelForProperty(property);
  }
  let html =
    `<tr><td colspan="2"><a href="#${href}">${name}</a></td>` +
    `<td>${typeLabel}</td></tr>`;

  return html;
};

const oneOfRow = (name, property) => {
  const length = Object.keys(property.oneOf).length;
  let html =
    `<tr><th rowspan="${length}">${name}</th>` +
    `<td rowspan="${length}">One of:</td>`;
  let isFirstRow = true;

  for (const key in property.oneOf) {
    let oneOfItem = property.oneOf[key];

    if (!isFirstRow) {
      html += "<tr>";
    }

    // oneOf allows common properties to sit above the oneOf items so need to support
    // both approaches
    html += `<td>${getLabelForProperty(
      oneOfItem.type ? oneOfItem : property
    )}</td></tr>`;

    isFirstRow = false;
  }

  return html;
};

const oneOfArrayRow = (oneOf) => {
  const length = oneOf.length;
  let html = "";
  let isFirstRow = true;

  oneOf.forEach((one) => {
    if (isFirstRow) {
      html += `<tr><th colspan="2" rowspan="${length}">One of:</th>`;
    } else {
      html += "<tr>";
    }

    html += `<th>${getLabelForProperty(one)}</th></tr>`;

    isFirstRow = false;
  });

  return html;
};

const propertiesTable = (schema, parentKey) => {
  let properties = schema.properties;

  let html =
    "<table>" +
    '<thead><tr><th colspan="2">Name</th><th>Type</th></tr></thead>' +
    "<tbody>";

  for (const key in properties) {
    let property = properties[key];

    if (
      property &&
      property.oneOf != undefined &&
      (!property.$ref || property.$ref.match(/^#/))
    ) {
      html += oneOfRow(key, property);
    } else {
      html += propertyRow(key, property, parentKey);
    }
  }

  let oneOf = schema.oneOf;
  if (Array.isArray(oneOf) && oneOf.length > 0) {
    html += oneOfArrayRow(oneOf);
  }

  html += "</tbody></table>";

  return html;
};

Handlebars.registerHelper("concat", function () {
  const subArguments = [...arguments].slice(0, -1);
  return subArguments.join("");
});

Handlebars.registerHelper("match", function (string, regex) {
  let result;
  if (typeof string === "string") {
    result = string?.match(regex);
  }
  return result;
});

Handlebars.registerHelper("propertiesTable", function (schema, parentKey) {
  return new Handlebars.SafeString(propertiesTable(schema, parentKey));
});

Handlebars.registerHelper("propertyTypeRow", function (property) {
  let html;
  if (property.oneOf != undefined) {
    html = `<tr>${oneOfRow("Type", property)}</tr>`;
  } else {
    const typeLabel = getLabelForProperty(property);
    if (typeLabel) {
      html = `<tr><th>Type</th><td colspan="2">${typeLabel}</td></tr>`;
    }
  }
  return html ? new Handlebars.SafeString(html) : "";
});

var RendererMarkdown = function (templatePath, options) {
  options = options || {};
  this.templatePath = templatePath;
  this.addFrontMatter =
    options.addFrontMatter !== undefined ? options.addFrontMatter : false;
  this.displaySchema =
    options.displaySchema !== undefined ? options.displaySchema : false;
};

RendererMarkdown.prototype.setup = async function () {
  const schemaTemplatePath = path.join(this.templatePath, "schema.hbs");
  let templateSource = await readFile(schemaTemplatePath);
  this.templateSchema = Handlebars.compile(templateSource.toString());
};

RendererMarkdown.prototype.renderSchema = function (data) {
  data.displaySchema = this.displaySchema;
  let result = this.templateSchema(data);

  // this fixes tables that have been broken by unintended double line breaks
  // introduced by optional rows.
  // could possibly be avoided by use of ~ in templates
  result = result.replace(/\|(\n)+\|/g, "|\n|");

  if (this.addFrontMatter) {
    let frontMatter = `title: ${data.schema.title || data.schema.id}\n`;
    if (data.schema.description) {
      const description = data.schema.description.replace(/\n/g, "");
      frontMatter += `description: ${description}\n`;
    }

    result = `---\n${frontMatter}\n---\n${result}`;
  }

  return result;
};

module.exports = RendererMarkdown;
