FROM ghcr.io/formancehq/base:scratch
COPY search /usr/bin/search
ENV OTEL_SERVICE_NAME search
ENTRYPOINT ["/usr/bin/search"]
CMD ["serve"]
