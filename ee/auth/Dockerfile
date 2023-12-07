FROM golang:1.20-alpine AS builder
WORKDIR /src
COPY components/auth components/auth
COPY libs libs
WORKDIR /src/ee/auth
ARG CGO_ENABLED=0
ARG APP_SHA
ARG VERSION
RUN --mount=type=cache,id=gobuild,target=/root/.cache/go-build --mount=type=cache,id=gomodcache,target=/go/pkg/mod go build -o auth \
    -ldflags="-X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Version=${VERSION} \
    -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.BuildDate=$(date +%s) \
    -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Commit=${APP_SHA}" ./

FROM alpine:3.16
RUN apk update && apk add ca-certificates curl
COPY --from=builder /src/ee/auth/auth /main
ENV OTEL_SERVICE_NAME auth
ENTRYPOINT ["/main"]
CMD ["--help"]
