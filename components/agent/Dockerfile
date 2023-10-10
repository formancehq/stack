# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS builder
WORKDIR /src
RUN apk update && apk add git
COPY libs/go-libs libs/go-libs
COPY components/operator components/operator
COPY components/agent components/agent
WORKDIR components/agent

ARG APP_SHA
ARG VERSION
ARG CGO_ENABLED=0
RUN --mount=type=cache,id=gobuild,target=/root/.cache/go-build --mount=type=cache,id=gomodcache,target=/go/pkg/mod \
    CGO_ENABLED=$CGO_ENABLED GOOS=linux \
    go build -o agent -v \
    -ldflags="-X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Version=${VERSION} \
        -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.BuildDate=$(date +%s) \
        -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Commit=${APP_SHA}" ./

FROM alpine:3.18
RUN apk update && apk add ca-certificates curl
COPY --from=builder /src/components/agent/agent /usr/local/bin/agent
ENTRYPOINT ["/usr/local/bin/agent"]
ENV OTEL_SERVICE_NAME /usr/local/bin/agent
