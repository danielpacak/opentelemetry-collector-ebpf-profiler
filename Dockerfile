# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23.6

FROM --platform=$BUILDPLATFORM alpine:3.19 AS certs
RUN apk --update add ca-certificates

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build-stage
WORKDIR /build

COPY ./manifest.yaml manifest.yaml
COPY ./exporter exporter

RUN --mount=type=cache,target=/root/.cache/go-build GOARCH=$TARGETARCH go install go.opentelemetry.io/collector/cmd/builder@v0.131.0
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOARCH=$TARGETARCH builder --config manifest.yaml

FROM --platform=$BUILDPLATFORM gcr.io/distroless/base:latest

COPY ./collector-config.yaml /otelcol/collector-config.yaml
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --chmod=755 --from=build-stage /build/_build /otelcol

ENTRYPOINT ["/otelcol/otelcol-ebpf-profiler"]
CMD ["--config", "/otelcol/collector-config.yaml"]

EXPOSE 4317 4318 12001
