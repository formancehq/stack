FROM ghcr.io/formancehq/base:scratch
COPY orchestration /usr/bin/orchestration
ENV OTEL_SERVICE_NAME orchestration
ENTRYPOINT ["/usr/bin/orchestration"]
CMD ["serve"]
