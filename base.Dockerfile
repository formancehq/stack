FROM ubuntu:22.04 AS base
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch AS scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt