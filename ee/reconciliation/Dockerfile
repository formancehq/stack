FROM golang:1.20-bullseye AS builder

RUN apt-get update && \
    apt-get install -y gcc-aarch64-linux-gnu gcc-x86-64-linux-gnu && \
    ln -s /usr/bin/aarch64-linux-gnu-gcc /usr/bin/arm64-linux-gnu-gcc  && \
    ln -s /usr/bin/x86_64-linux-gnu-gcc /usr/bin/amd64-linux-gnu-gcc

ARG TARGETARCH
ARG APP_SHA
ARG VERSION

WORKDIR /src

# get deps first so it's cached
COPY . .
WORKDIR /src/ee/reconciliation
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go mod vendor

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH \
    CC=$TARGETARCH-linux-gnu-gcc \
    go build -o bin/reconciliation \
    -ldflags="-X github.com/formancehq/reconciliation/cmd.Version=${VERSION} \
    -X github.com/formancehq/reconciliation/cmd.BuildDate=$(date +%s) \
    -X github.com/formancehq/reconciliation/cmd.Commit=${APP_SHA}" ./

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/ee/reconciliation/bin/reconciliation /usr/local/bin/reconciliation

EXPOSE 8080

ENTRYPOINT ["reconciliation"]
ENV OTEL_SERVICE_NAME=reconciliation
CMD ["serve"]
