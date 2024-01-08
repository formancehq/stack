FROM scratch
COPY ledger /usr/bin/ledger
ENV OTEL_SERVICE_NAME ledger
ENTRYPOINT ["/usr/bin/ledger"]
CMD ["serve"]
