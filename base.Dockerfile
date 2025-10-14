FROM ubuntu:24.04 AS base
RUN apt update && DEBIAN_FRONTEND=noninteractive apt install -y ca-certificates curl tzdata && rm -rf /var/lib/apt/lists/*

FROM alpine:3.22 AS certs
RUN apk --update add ca-certificates

FROM scratch AS scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
