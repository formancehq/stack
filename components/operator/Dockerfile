# Build the manager binary
FROM golang:1.19-alpine as builder
WORKDIR /src
ENV CGO_ENABLED=0
COPY libs/go-libs libs/go-libs
COPY components/payments components/payments
COPY components/search components/search
COPY components/operator components/operator
WORKDIR /src/components/operator
RUN apk add --no-cache postgresql-client
RUN --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    --mount=type=cache,id=gomodcache,target=/go/pkg/mod \
    go build -v -a -o manager main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot as release
LABEL org.opencontainers.image.source=https://github.com/formancehq/operator
WORKDIR /
COPY --from=builder /src/components/operator/manager /usr/bin/operator
USER 65532:65532
ENTRYPOINT ["/usr/bin/operator"]
