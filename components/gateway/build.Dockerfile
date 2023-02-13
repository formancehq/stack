FROM ubuntu:jammy as gateway-builder
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY gateway /usr/bin/gateway
ENTRYPOINT ["/usr/bin/gateway"]

FROM golang:1.18 as builder
WORKDIR /src
COPY --from=gateway-builder /usr/bin/gateway /usr/bin/gateway
COPY builder-config.json .
RUN /usr/bin/gateway build \
    --caddy-builder-config-path builder-config.json \
    --caddy-binary-output-path /usr/bin/caddy

FROM ubuntu:jammy
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY --from=builder /usr/bin/caddy /usr/bin/caddy
COPY Caddyfile /etc/caddy/Caddyfile
ENV OTEL_SERVICE_NAME gateway
ENTRYPOINT ["/usr/bin/caddy"]
CMD ["run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]