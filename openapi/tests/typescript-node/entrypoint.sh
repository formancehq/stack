#!/bin/sh

set -x;
set -e;

cd /src/sdks/typescript-node;
npm install;
npm run build;
cd /src/openapi/tests/typescript-node/testing-app;
npm install;
npm run build;
npm run test;
