FROM ghcr.io/formancehq/base:scratch
COPY operator /usr/bin/operator
ENV OTEL_SERVICE_NAME operator
ENTRYPOINT ["/usr/bin/operator"]
