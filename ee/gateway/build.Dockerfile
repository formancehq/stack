FROM ghcr.io/formancehq/base:22.04
ADD https://raw.githubusercontent.com/formancehq/stack/main/ee/gateway/Caddyfile /etc/caddy/Caddyfile
COPY gateway /usr/bin/caddy
ENV OTEL_SERVICE_NAME gateway
ENTRYPOINT ["/usr/bin/caddy"]
CMD ["run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
