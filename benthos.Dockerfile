FROM jeffail/benthos:3.65
WORKDIR /config
COPY benthos /config
CMD ["-c", "/config/config.yml"]
