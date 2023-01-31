FROM ubuntu:jammy
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY webhooks /usr/bin/webhooks
ENV OTEL_SERVICE_NAME webhooks
ENTRYPOINT ["/usr/bin/webhooks"]
CMD ["server"]
