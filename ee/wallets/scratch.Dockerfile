FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY wallets /usr/bin/wallets
ENV OTEL_SERVICE_NAME wallets
ENTRYPOINT ["/usr/bin/wallets"]
CMD ["server"]
