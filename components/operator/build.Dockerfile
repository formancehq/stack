FROM ghcr.io/formancehq/base:22.04
RUN apt-get install postgresql-client
COPY operator /usr/bin/operator
ENV OTEL_SERVICE_NAME operator
ENTRYPOINT ["/usr/bin/operator"]
