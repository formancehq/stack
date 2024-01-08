FROM scratch
COPY fctl /usr/bin/fctl
ENV OTEL_SERVICE_NAME fctl
ENTRYPOINT ["/usr/bin/fctl"]
