FROM ubuntu:22.04
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY payments /usr/bin/payments
ENV OTEL_SERVICE_NAME payments
ENTRYPOINT ["/usr/bin/payments"]
CMD ["server"]
