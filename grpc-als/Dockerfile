# FROM golang:1.15.10-stretch AS builder
FROM golang:1.15-alpine3.12 AS base

ENV GO111MODULE=on \
    CGO_ENABLED=1

WORKDIR /gRPC-als
COPY . .
RUN go build ./...
RUN ls
ENTRYPOINT [ "./gRPC-als" ]