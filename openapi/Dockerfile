FROM node:19-alpine
ENV PORT 3000
WORKDIR /app
COPY package.json ./
COPY yarn.lock ./
COPY main.js /app/main.js
COPY sdk/generate.json /app/sdk/generate.json
RUN apk add --no-cache ca-certificates curl
RUN yarn install --production
EXPOSE 3000
CMD ["yarn", "run", "start"]