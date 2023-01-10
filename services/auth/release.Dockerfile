FROM ubuntu:jammy
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY auth /usr/bin/auth
EXPOSE 8080
ENV OTEL_SERVICE_NAME auth
ENTRYPOINT ["/auth"]
CMD ["--help"]
