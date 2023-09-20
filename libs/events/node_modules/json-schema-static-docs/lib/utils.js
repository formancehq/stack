const determineSchemaRelativePath = (schemaFilename, schemaInputPath) => {
  let outputFilename = schemaFilename.substr(schemaInputPath.length);
  if (outputFilename.substr(0, 1) === "/") {
    outputFilename = outputFilename.substr(1);
  }
  return outputFilename;
};

module.exports = {
  determineSchemaRelativePath: determineSchemaRelativePath,
};
