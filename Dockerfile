FROM golang:1.27.1-alpine

RUN apk update && \
  apk --no-cache add binutils

