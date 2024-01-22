FROM ghcr.io/formancehq/base:scratch
COPY wallets /usr/bin/wallets
ENV OTEL_SERVICE_NAME wallets
ENTRYPOINT ["/usr/bin/wallets"]
CMD ["server"]
