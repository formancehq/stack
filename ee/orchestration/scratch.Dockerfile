FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY orchestration /usr/bin/orchestration
ENV OTEL_SERVICE_NAME orchestration
ENTRYPOINT ["/usr/bin/orchestration"]
CMD ["serve"]
