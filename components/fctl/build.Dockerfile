FROM ghcr.io/formancehq/base:22.04
COPY fctl /usr/bin/fctl
ENV OTEL_SERVICE_NAME fctl
ENTRYPOINT ["/usr/bin/fctl"]
