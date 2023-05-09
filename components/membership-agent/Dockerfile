# syntax=docker/dockerfile:1
FROM golang:1.19-alpine AS builder
WORKDIR /src
RUN apk update && apk add git
COPY components/membership-agent components/membership-agent
WORKDIR components/membership-agent

ARG APP_SHA
ARG VERSION
ARG CGO_ENABLED=0
RUN --mount=type=cache,id=gobuild,target=/root/.cache/go-build --mount=type=cache,id=gomodcache,target=/go/pkg/mod \
    CGO_ENABLED=$CGO_ENABLED GOOS=linux \
    go build -o membership-agent -v \
    -ldflags="-X github.com/formancehq/membership-agent/cmd.Version=${VERSION} \
    -X github.com/formancehq/membership-agent/cmd.BuildDate=$(date +%s) \
    -X github.com/formancehq/membership-agent/cmd.Commit=${APP_SHA}" ./

FROM alpine:3.16
RUN apk update && apk add ca-certificates curl
COPY --from=builder /src/components/membership-agent/membership-agent /usr/local/bin/membership-agent
ENTRYPOINT ["/usr/local/bin/membership-agent"]
ENV OTEL_SERVICE_NAME /usr/local/bin/membership-agent
