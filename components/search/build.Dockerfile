FROM ubuntu:jammy
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY search /usr/bin/search
ENV OTEL_SERVICE_NAME search
ENTRYPOINT ["/usr/bin/search"]
CMD ["serve"]
