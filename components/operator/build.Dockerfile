FROM ubuntu:jammy
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY operator /usr/bin/operator
ENV OTEL_SERVICE_NAME operator
ENTRYPOINT ["/usr/bin/operator"]
