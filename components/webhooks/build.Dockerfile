FROM ghcr.io/formancehq/base:22.04
COPY webhooks /usr/bin/webhooks
ENV OTEL_SERVICE_NAME webhooks
ENTRYPOINT ["/usr/bin/webhooks"]
CMD ["server"]
