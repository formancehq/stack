FROM ghcr.io/formancehq/base:scratch
COPY agent /usr/bin/agent
ENV OTEL_SERVICE_NAME agent
ENTRYPOINT ["/usr/bin/agent"]
