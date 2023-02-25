#!/bin/sh

set -x;
set -e;

cd testing-app;
npm install;
npm run build;
npm run test;
