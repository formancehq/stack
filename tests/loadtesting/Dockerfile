FROM node:19-alpine3.16
# Install ca-certificates
RUN apk update && apk add ca-certificates curl
# Install k6
ADD https://github.com/formancehq/xk6-extension/releases/download/v0.1.1/k6 /usr/bin/k6
RUN chmod +x /usr/bin/k6
# Build deps
WORKDIR /k6
COPY . .
RUN yarn install
RUN yarn run build
ENTRYPOINT ["/usr/bin/k6"]
