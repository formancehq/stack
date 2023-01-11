FROM golang:1.18 AS builder
ARG APP_SHA
ARG VERSION
WORKDIR /src
COPY . .
WORKDIR /src/services/payments
RUN go mod download
RUN GOOS=linux go build -o payments \
    -ldflags="-X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Version=${VERSION} \
    -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.BuildDate=$(date +%s) \
    -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Commit=${APP_SHA}" ./

FROM ubuntu:jammy
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY --from=builder /src/services/payments/payments /payments
EXPOSE 3068
ENV OTEL_SERVICE_NAME payments
ENTRYPOINT ["/payments"]
CMD ["server"]
