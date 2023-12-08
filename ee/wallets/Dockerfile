FROM golang:1.20 AS builder
ARG APP_SHA
ARG VERSION
WORKDIR /src
COPY libs/clients/go /src/libs/clients/go
COPY libs /src/libs
COPY components/wallets /src/ee/wallets
WORKDIR /src/ee/wallets
RUN --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    --mount=type=cache,id=gomodcache,target=/go/pkg/mod \
    GOOS=linux go build -o wallets \
    -ldflags="-X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Version=${VERSION} \
    -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.BuildDate=$(date +%s) \
    -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Commit=${APP_SHA}" ./

FROM ubuntu:22.04
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY --from=builder /src/ee/wallets/wallets /wallets
EXPOSE 3068
ENV OTEL_SERVICE_NAME wallets
ENTRYPOINT ["/wallets"]
CMD ["serve"]
