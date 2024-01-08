FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY operator /usr/bin/operator
ENV OTEL_SERVICE_NAME operator
ENTRYPOINT ["/usr/bin/operator"]
