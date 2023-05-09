FROM ubuntu:jammy
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY stargate /usr/bin/stargate
ENV OTEL_SERVICE_NAME stargate
ENTRYPOINT ["/usr/bin/stargate"]
CMD ["server"]
