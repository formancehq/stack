FROM golang:1.20-alpine AS builder
ARG APP_SHA
ARG VERSION
WORKDIR /src
COPY components/search /src/ee/search
COPY libs/go-libs /src/libs/go-libs
WORKDIR /src/ee/search
RUN GOOS=linux go build -o search \
    -ldflags="-X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Version=${VERSION} \
    -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.BuildDate=$(date +%s) \
    -X $(cat go.mod |head -1|cut -d \  -f2)/cmd.Commit=${APP_SHA}" ./

FROM alpine:3.16
RUN apk update && apk add ca-certificates curl
COPY --from=builder /src/ee/search/search /search
EXPOSE 3068
ENV OTEL_SERVICE_NAME search
ENTRYPOINT ["/search"]
CMD ["serve"]
