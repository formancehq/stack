FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY search /usr/bin/search
ENV OTEL_SERVICE_NAME search
ENTRYPOINT ["/usr/bin/search"]
CMD ["serve"]
