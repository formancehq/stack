FROM ghcr.io/formancehq/base:22.04
COPY operator /usr/bin/operator
ENV OTEL_SERVICE_NAME operator
ENTRYPOINT ["/usr/bin/operator"]
