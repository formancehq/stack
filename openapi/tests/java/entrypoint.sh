#!/bin/sh

set -x;
set -e;

cd /src/sdks/java
mvn clean install -DskipTests;
cd /src/openapi/tests/java
mvn test;
