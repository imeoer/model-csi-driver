# syntax=docker/dockerfile:1.4

FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  make release

FROM ubuntu:24.04
COPY --from=builder /app/model-csi-driver /usr/bin/model-csi-driver
ENTRYPOINT ["/usr/bin/model-csi-driver"]
