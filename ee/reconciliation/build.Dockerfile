FROM ghcr.io/formancehq/base:22.04
COPY reconciliation /usr/bin/reconciliation
ENV OTEL_SERVICE_NAME reconciliation
ENTRYPOINT ["/usr/bin/reconciliation"]
CMD ["serve"]
