FROM golang:1.20-alpine
WORKDIR /src

COPY sdks/go sdks/go
COPY openapi/tests/go openapi/tests/go

WORKDIR openapi/tests/go
ENTRYPOINT ["/src/openapi/tests/go/entrypoint.sh"]
