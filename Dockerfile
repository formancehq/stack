# syntax=docker/dockerfile:1
FROM ghcr.io/formancehq/gateway:v0.1.7 as gateway
FROM ghcr.io/formancehq/ledger:v1.10.3 as ledger
FROM ghcr.io/formancehq/payments:v0.9.7 as payments
FROM ghcr.io/formancehq/orchestration:v0.1.5 as orchestration
FROM ghcr.io/formancehq/auth:v0.4.3 as auth
FROM ghcr.io/formancehq/search:v0.8.0 as search
FROM ghcr.io/formancehq/wallets:v0.4.3 as wallets
FROM ghcr.io/formancehq/webhooks:v0.6.6 as webhooks
FROM ghcr.io/formancehq/control:v1.7.0 as control
FROM ghcr.io/formancehq/auth-dex:latest as dex
FROM jeffail/benthos:4.23.0 as benthos

FROM golang:1.20 as builder
WORKDIR /tmp
RUN apt update && apt install -y wget
RUN wget https://github.com/F1bonacc1/process-compose/archive/refs/tags/v0.51.4.tar.gz \
    && tar -xvf v0.51.4.tar.gz
WORKDIR /tmp/process-compose-0.51.4
RUN go build -o /usr/bin/process-compose ./src

FROM node:18
WORKDIR /tmp
COPY --from=builder /usr/bin/process-compose /usr/bin/process-compose
COPY --from=gateway /usr/bin/caddy /usr/bin/gateway
COPY --from=orchestration /usr/bin/orchestration /usr/bin/orchestration
COPY --from=dex /usr/local/bin/dex /usr/bin/dex
COPY --from=auth /usr/bin/auth /usr/bin/auth
COPY --from=ledger /usr/local/bin/numary /usr/bin/ledger
COPY --from=payments /usr/bin/payments /usr/bin/payments
COPY --from=search /usr/bin/search /usr/bin/search
COPY --from=wallets /usr/bin/wallets /usr/bin/wallets
COPY --from=webhooks /usr/bin/webhooks /usr/bin/webhooks
COPY --from=benthos /benthos /usr/bin/benthos
COPY --from=control /app /app
COPY .local /etc/formance
COPY ./components/search/benthos /benthos
WORKDIR /app
COPY .local/process-compose.yaml /running/process-compose.yaml
CMD ["process-compose", "up", "--config", "/running/process-compose.yaml", "--tui=false", "-p", "9090"]
