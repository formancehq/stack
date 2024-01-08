FROM scratch
COPY wallets /usr/bin/wallets
ENV OTEL_SERVICE_NAME wallets
ENTRYPOINT ["/usr/bin/wallets"]
CMD ["server"]
