FROM ubuntu:22.04
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY auth /usr/bin/auth
ENV OTEL_SERVICE_NAME auth
ENTRYPOINT ["/usr/bin/auth"]
CMD ["serve"]
