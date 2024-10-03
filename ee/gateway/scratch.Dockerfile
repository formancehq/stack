FROM ghcr.io/formancehq/base:scratch
ADD https://raw.githubusercontent.com/formancehq/gateway/main/Caddyfile /etc/caddy/Caddyfile
COPY gateway /usr/bin/caddy
ENV OTEL_SERVICE_NAME gateway
ENTRYPOINT ["/usr/bin/caddy"]
CMD ["run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]