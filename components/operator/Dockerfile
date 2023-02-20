# Build the manager binary
FROM golang:1.18-alpine as builder
WORKDIR /workspace
ENV CGO_ENABLED=0
ENV GOOS=linux
COPY . .
WORKDIR /workspace/components/operator
RUN go build -v -a -o manager main.go

FROM golang:1.18-alpine as reloader
RUN go install github.com/cosmtrek/air@latest

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot as release
LABEL org.opencontainers.image.source=https://github.com/formancehq/operator
WORKDIR /
COPY --from=builder /workspace/components/operator/manager /usr/bin/operator
USER 65532:65532
ENTRYPOINT ["/usr/bin/operator"]
