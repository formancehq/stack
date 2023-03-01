FROM python:3.12.0a5-alpine3.17
WORKDIR /src

RUN apk update && apk add git

COPY sdks/python sdks/python
COPY openapi/tests/python openapi/tests/python

WORKDIR openapi/tests/python
ENTRYPOINT ["/src/openapi/tests/python/entrypoint.sh"]
