FROM ghcr.io/formancehq/base:22.04
COPY agent /usr/bin/agent
ENV OTEL_SERVICE_NAME agent
ENTRYPOINT ["/usr/bin/agent"]
