FROM golang:1.19 as builder
ARG APP_SHA
ARG VERSION
WORKDIR /src
COPY . .
WORKDIR /src/components/gateway
RUN go mod download
RUN GOOS=linux go build -o caddy \
    -ldflags="-X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Version=${VERSION} \
    -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.BuildDate=$(date +%s) \
    -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Commit=${APP_SHA}" ./

FROM ubuntu:jammy
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY --from=builder src/components/gateway/caddy /usr/bin/caddy
COPY --from=builder /src/config/gateway/Caddyfile /etc/caddy/Caddyfile
ENV OTEL_SERVICE_NAME gateway
ENTRYPOINT ["/usr/bin/caddy"]
CMD ["run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
