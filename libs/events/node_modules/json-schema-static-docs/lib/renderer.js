let RendererMarkdown = require("./renderer-markdown");

function Renderer(templatePath, options) {
  options = options || {};
  this.renderer = new RendererMarkdown(templatePath, options);
}

Renderer.prototype.setup = async function () {
  await this.renderer.setup();
};

Renderer.prototype.renderSchema = function (data) {
  return this.renderer.renderSchema(data);
};

module.exports = Renderer;
