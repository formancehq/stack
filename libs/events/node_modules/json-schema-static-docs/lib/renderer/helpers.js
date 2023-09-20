const upperCaseFirstCharacter = (value) => {
  return typeof value === "string"
    ? value.substr(0, 1).toUpperCase() + value.substr(1)
    : value;
};

const convertRelativeUrlToHtmlDocsUrl = (href) => {
  return typeof href === "string"
    ? href
        // document pages use a .html extension
        // given $ref=schema.json#/definitions/fooBar
        // we want to link to schema.html#fooBar
        .replace(/\.json/, ".html")
        .replace(/\.yml/, ".html")
        // convert any defs into suitable anchors
        .replace(/#\/(definitions|\$defs)\//, "#")
        // anchor IDs will be lower-case
        .replace(/#(.)+$/, (match, offset, string) => {
          return match.toLowerCase();
        })
    : href;
};

const isHttpRef = (ref) => {
  return typeof ref === "string" && ref.match(/^http(s)?:/);
};
const isRelativeRef = (ref) => {
  return typeof ref === "string" && !isHttpRef(ref) && !ref.match(/^#/);
};

const formatLabel = (label) => {
  if (Array.isArray(label)) {
    label = label.join(", ");
    label = `[${label}]`;
  } else {
    label = label || "";
    label = upperCaseFirstCharacter(label);
  }
  return label;
};

const getHtmlAnchorForRef = (ref, label) => {
  let refParts = ref.split("/");
  let filename = refParts[refParts.length - 1];

  let htmlAnchor;

  if (isHttpRef(ref)) {
    htmlAnchor = `<a href="${ref}">${filename}</a>`;
  } else if (isRelativeRef(ref)) {
    // this is a little ugly,
    // most of the time the $ref is passed in as the label by the template
    // may be a nicer way to handle this as an optional argument in a handlebars helper
    label =
      !label || label === ref
        ? convertRelativeUrlToHtmlDocsUrl(ref)
        : formatLabel(label);
    htmlAnchor = `<a href="${convertRelativeUrlToHtmlDocsUrl(
      ref
    )}">${label}</a>`;
  } else {
    label = formatLabel(label);
    htmlAnchor = `<a href="${ref}">${label}</a>`;
  }

  return htmlAnchor;
};

const getLabelForProperty = (property) => {
  let result;
  let ref = property.items?.$ref || property.$ref;
  let _isRelativeRef = isRelativeRef(ref);

  if (_isRelativeRef && property.type === "array") {
    result = `Array [${getHtmlAnchorForRef(ref, property.items.title)}]`;
  } else if (_isRelativeRef && property.type === "object") {
    result = `Object (of type ${getHtmlAnchorForRef(ref, property.title)})`;
  } else {
    result = formatLabel(property.type);
  }

  return result;
};

module.exports = {
  getHtmlAnchorForRef,
  getLabelForProperty,
  upperCaseFirstCharacter,
};
