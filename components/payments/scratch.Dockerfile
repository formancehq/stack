FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY payments /usr/bin/payments
ENV OTEL_SERVICE_NAME payments
ENTRYPOINT ["/usr/bin/payments"]
CMD ["server"]
