FROM jeffail/benthos:3.65
ENV BENTHOS_MODE=kafka
WORKDIR /config
COPY benthos /config
CMD ["-c", "/config/config_${BENTHOS_MODE}.yml"]