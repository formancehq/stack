FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY fctl /usr/bin/fctl
ENV OTEL_SERVICE_NAME fctl
ENTRYPOINT ["/usr/bin/fctl"]
