FROM ubuntu:22.04
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
