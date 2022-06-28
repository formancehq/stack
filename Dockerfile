ARG VERSION=latest

FROM golang:1.18-buster as src
COPY . /app
WORKDIR /app

FROM src as dev
RUN apt-get update && apt-get install -y ca-certificates git-core ssh
RUN go install github.com/cespare/reflex@latest
RUN git config --global url.ssh://git@github.com/numary.insteadOf https://github.com/numary

FROM src as compiler
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=${VERSION}" .

FROM alpine as app
RUN apk add --no-cache ca-certificates curl
COPY --from=compiler /app/webhooks-cloud /usr/local/bin/webhooks
EXPOSE 8080
CMD ["webhooks"]
