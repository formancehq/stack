const RendererMarkdown = require("../lib/renderer-markdown.js");
const Handlebars = require("handlebars");

test("defines expected helpers", () => {
  expect(Handlebars.helpers["getHtmlAnchorForRef"]).toBeDefined();
});
