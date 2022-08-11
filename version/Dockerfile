FROM golang:1.18-alpine as build
WORKDIR /app
COPY * ./
RUN go mod download
RUN go build -o version

FROM alpine:3.16
WORKDIR /
COPY --from=build /app/version /usr/bin/version
EXPOSE 3333
ENTRYPOINT ["version"]
