FROM golang:1.20-alpine as builder
ARG APP_SHA
ARG VERSION
WORKDIR /src
COPY components/gateway components/gateway
WORKDIR /src/ee/gateway
RUN --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    --mount=type=cache,id=gomodcache,target=/go/pkg/mod \
    go build -o caddy \
    -ldflags="-X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Version=${VERSION} \
    -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.BuildDate=$(date +%s) \
    -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Commit=${APP_SHA}" ./

FROM alpine:3.16
RUN apk update && apk add ca-certificates curl
COPY components/gateway/Caddyfile /etc/caddy/Caddyfile
COPY --from=builder src/components/gateway/caddy /usr/bin/caddy
ENV OTEL_SERVICE_NAME gateway
ENTRYPOINT ["/usr/bin/caddy"]
CMD ["run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
