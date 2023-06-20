FROM ubuntu:22.04
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY wallets /usr/bin/wallets
ENV OTEL_SERVICE_NAME wallets
ENTRYPOINT ["/usr/bin/wallets"]
CMD ["server"]
