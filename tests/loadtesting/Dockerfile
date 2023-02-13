FROM node:16-buster
# Install ca-certificates
RUN apt update && apt install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
# Install k6
ADD https://github.com/formancehq/xk6-extension/releases/download/v0.1.1/k6 /usr/bin/k6
RUN chmod +x /usr/bin/k6
# Build deps
WORKDIR /k6
COPY . .
RUN yarn install
RUN yarn run build
ENTRYPOINT ["/usr/bin/k6"]
