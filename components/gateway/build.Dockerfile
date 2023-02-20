FROM ubuntu:jammy
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY gateway /usr/bin/caddy
COPY Caddyfile /etc/caddy/Caddyfile
ENV OTEL_SERVICE_NAME gateway
ENTRYPOINT ["/usr/bin/caddy"]
CMD ["run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
