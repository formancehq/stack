FROM golang:1.18 as gateway-builder
ARG APP_SHA
ARG VERSION
WORKDIR /src
COPY . .
WORKDIR /src/components/gateway
RUN go mod download
RUN GOOS=linux go build -o gateway \
    -ldflags="-X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Version=${VERSION} \
    -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.BuildDate=$(date +%s) \
    -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Commit=${APP_SHA}" ./

FROM golang:1.18 as builder
WORKDIR /src
COPY --from=gateway-builder /src/components/gateway/gateway /gateway
COPY --from=gateway-builder /src/config/gateway/Caddyfile .
COPY --from=gateway-builder /src/config/gateway/builder-config.json .
RUN /gateway build \
    --caddy-builder-config-path builder-config.json \
    --caddy-binary-output-path /usr/bin/caddy

FROM ubuntu:jammy
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY --from=builder /usr/bin/caddy /usr/bin/caddy
COPY --from=builder /src/Caddyfile /etc/caddy/Caddyfile
ENV OTEL_SERVICE_NAME gateway
ENTRYPOINT ["/usr/bin/caddy"]
CMD ["run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]