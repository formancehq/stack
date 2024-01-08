FROM scratch
COPY auth /usr/bin/auth
ENV OTEL_SERVICE_NAME auth
ENTRYPOINT ["/usr/bin/auth"]
CMD ["serve"]
