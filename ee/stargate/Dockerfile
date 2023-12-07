FROM golang:1.20-alpine as builder
ARG APP_SHA
ARG VERSION
WORKDIR /src
COPY libs/go-libs /src/libs/go-libs
COPY components/stargate /src/ee/stargate
WORKDIR /src/ee/stargate
RUN --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    --mount=type=cache,id=gomodcache,target=/go/pkg/mod \
    go build -o stargate \
    -ldflags="-X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Version=${VERSION} \
    -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.BuildDate=$(date +%s) \
    -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Commit=${APP_SHA}" ./

FROM alpine:3.16
RUN apk update && apk add ca-certificates curl
COPY --from=builder src/components/stargate/stargate /usr/bin/stargate
ENV OTEL_SERVICE_NAME stargate
RUN chmod +x /usr/bin/stargate
ENTRYPOINT ["/usr/bin/stargate"]
CMD ["client"]
