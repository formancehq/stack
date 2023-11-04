FROM golang:1.21-alpine
RUN go install github.com/cespare/reflex@latest
RUN apk update && apk add curl
