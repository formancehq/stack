FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY stargate /usr/bin/stargate
ENV OTEL_SERVICE_NAME stargate
ENTRYPOINT ["/usr/bin/stargate"]
CMD ["client"]
