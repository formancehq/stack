FROM maven:3.8-jdk-8-slim
WORKDIR /src

COPY sdks/java sdks/java
COPY openapi/tests/java openapi/tests/java

WORKDIR openapi/tests/java
ENTRYPOINT ["/src/openapi/tests/java/entrypoint.sh"]
