ARG VERSION=latest

FROM golang:1.18-buster as src
COPY . /app
WORKDIR /app

FROM src as dev
RUN apt-get update && apt-get install -y ca-certificates git-core ssh

FROM src as compiler
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=${VERSION}" .

FROM alpine as app
RUN apk add --no-cache ca-certificates curl
COPY --from=compiler /app/search /usr/local/bin/search
EXPOSE 8080
CMD ["search"]
