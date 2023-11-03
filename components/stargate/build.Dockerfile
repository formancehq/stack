FROM ghcr.io/formancehq/base:22.04
COPY stargate /usr/bin/stargate
ENV OTEL_SERVICE_NAME stargate
ENTRYPOINT ["/usr/bin/stargate"]
CMD ["client"]