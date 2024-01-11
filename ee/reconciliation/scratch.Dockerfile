FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY reconciliation /usr/bin/reconciliation
ENV OTEL_SERVICE_NAME reconciliation
ENTRYPOINT ["/usr/bin/reconciliation"]
CMD ["serve"]
