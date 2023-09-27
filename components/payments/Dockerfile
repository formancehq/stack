FROM golang:1.19.3-bullseye AS builder

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
WORKDIR /src/components/payments
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go mod vendor

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH \
    CC=$TARGETARCH-linux-gnu-gcc \
    go build -o bin/payments \
    -ldflags="-X github.com/formancehq/payments/cmd.Version=${VERSION} \
    -X github.com/formancehq/payments/cmd.BuildDate=$(date +%s) \
    -X github.com/formancehq/payments/cmd.Commit=${APP_SHA}" ./

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/components/payments/bin/payments /usr/local/bin/payments

EXPOSE 8080

ENTRYPOINT ["payments"]
ENV OTEL_SERVICE_NAME=payments
CMD ["serve"]
