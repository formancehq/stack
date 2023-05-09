FROM ubuntu:jammy
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY membership-agent /usr/bin/membership-agent
ENV OTEL_SERVICE_NAME membership-agent
ENTRYPOINT ["/usr/bin/membership-agent"]
CMD ["server"]
