FROM ubuntu:22.04
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
ADD https://raw.githubusercontent.com/formancehq/stack/main/components/gateway/Caddyfile /etc/caddy/Caddyfile
COPY gateway /usr/bin/caddy
ENV OTEL_SERVICE_NAME gateway
ENTRYPOINT ["/usr/bin/caddy"]
CMD ["run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
