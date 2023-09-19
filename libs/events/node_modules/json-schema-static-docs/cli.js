#! /usr/bin/env node

const fs = require('fs');
const path = require('path');
const JsonSchemaStaticDocs = require('./lib/json-schema-static-docs');

var argv = require('optimist')
  .usage('Convert json schema into markdown docs.')
  .demand('i')
  .demand('o')
  .alias('i', 'inputPath')
  .describe('i', 'path to input directory')
  .alias('o', 'outputPath')
  .describe('o', 'path to output directory')
  .check(function(args) {
    if (!fs.existsSync(args.inputPath)) {
      throw 'Input path "' + args.inputPath + '" does not exist.';
    }
  })
  .argv;

( async () => {
  let jsonSchemaStaticDocs = new JsonSchemaStaticDocs({
    inputPath: argv.i,
    outputPath: argv.o
  });
  await jsonSchemaStaticDocs.generate();
  console.log('Documents generated into ' + argv.o);
})();

