FROM node:19-alpine3.16
WORKDIR /src
COPY sdks/typescript-node sdks/typescript-node
COPY openapi/tests/typescript-node openapi/tests/typescript-node

WORKDIR openapi/tests/typescript-node
ENTRYPOINT ["/src/openapi/tests/typescript-node/entrypoint.sh"]
