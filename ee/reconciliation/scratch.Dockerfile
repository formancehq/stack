FROM scratch
COPY reconciliation /usr/bin/reconciliation
ENV OTEL_SERVICE_NAME reconciliation
ENTRYPOINT ["/usr/bin/reconciliation"]
CMD ["serve"]
