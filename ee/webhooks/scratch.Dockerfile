FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY webhooks /usr/bin/webhooks
ENV OTEL_SERVICE_NAME webhooks
ENTRYPOINT ["/usr/bin/webhooks"]
CMD ["server"]
