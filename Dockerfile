FROM golang:1.25.1-alpine

RUN apk update && \
  apk --no-cache add binutils

