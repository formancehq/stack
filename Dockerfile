FROM node:latest

WORKDIR /app
RUN npm install -g @moonrepo/cli
COPY ./.moon/docker/workspace .
RUN moon docker setup
RUN moon run auth:lint
