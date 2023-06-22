# syntax=docker/dockerfile:1
FROM ghcr.io/formancehq/gateway:v0.1.7 as gateway
FROM ghcr.io/formancehq/ledger:v1.10.3 as ledger
FROM ghcr.io/formancehq/payments:v0.7.1 as payments
FROM ghcr.io/formancehq/auth:v0.4.3 as auth
FROM ghcr.io/formancehq/search:v0.7.0 as search
FROM ghcr.io/formancehq/wallets:v0.4.3 as wallets
FROM ghcr.io/formancehq/webhooks:v0.6.6 as webhooks
FROM ghcr.io/formancehq/control:v1.7.0 as control
FROM ghcr.io/formancehq/auth-dex:latest as dex
FROM jeffail/benthos:4.12.1 as benthos

FROM node:18
WORKDIR /tmp
RUN apt update && apt install -y wget &&\
  export ARCH=$(dpkg --print-architecture) && \
  wget https://github.com/F1bonacc1/process-compose/releases/download/v0.51.4/process-compose_Linux_$ARCH.tar.gz && \
  tar -xvf process-compose_Linux_$ARCH.tar.gz -C /tmp &&\
  cp -R process-compose /usr/bin/process-compose
COPY --from=gateway /usr/bin/caddy /usr/bin/gateway
COPY --from=dex /usr/local/bin/dex /usr/bin/dex
COPY --from=auth /usr/bin/auth /usr/bin/auth
COPY --from=ledger /usr/local/bin/numary /usr/bin/ledger
COPY --from=payments /usr/bin/payments /usr/bin/payments
COPY --from=search /usr/bin/search /usr/bin/search
COPY --from=wallets /usr/bin/wallets /usr/bin/wallets
COPY --from=webhooks /usr/bin/webhooks /usr/bin/webhooks
COPY --from=benthos /benthos /usr/bin/benthos
COPY --from=control /app /app
COPY ./config /etc/formance
COPY ./components/search/benthos /benthos
WORKDIR /app
COPY process-compose.yaml /running/process-compose.yaml
CMD ["process-compose", "up", "--config", "/running/process-compose.yaml", "--tui=false", "-p", "9090"]
