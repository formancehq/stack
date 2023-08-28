FROM golang:1.19 AS builder
ARG APP_SHA
ARG VERSION
WORKDIR /src
COPY . .
WORKDIR /src/components/fctl
RUN go mod download
RUN GOOS=linux go build -o fctl \
    -ldflags="-X $(cat go.mod |head -1|cut -d \  -f2)/pkg.Version=${VERSION} \
    -X $(cat go.mod |head -1|cut -d \  -f2)/pkg.BuildDate=$(date +%s) \
    -X $(cat go.mod |head -1|cut -d \  -f2)/pkg.Commit=${APP_SHA}" ./

FROM ubuntu:22.04
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY --from=builder /src/components/fctl/fctl /fctl
EXPOSE 3068
ENV OTEL_SERVICE_NAME fctl
ENTRYPOINT ["/fctl"]
CMD ["--help"]
