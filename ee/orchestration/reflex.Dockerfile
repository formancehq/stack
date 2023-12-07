FROM golang:1.19-alpine
RUN go install github.com/cespare/reflex@latest
