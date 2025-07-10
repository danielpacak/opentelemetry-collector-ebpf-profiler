FROM alpine:3.19 AS certs
RUN apk --update add ca-certificates

FROM golang:1.23.6 AS build-stage
WORKDIR /build

COPY ./manifest.yaml manifest.yaml

RUN --mount=type=cache,target=/root/.cache/go-build GO111MODULE=on go install go.opentelemetry.io/collector/cmd/builder@v0.129.0
RUN --mount=type=cache,target=/root/.cache/go-build builder --config manifest.yaml

FROM gcr.io/distroless/base:latest

COPY ./collector-config.yaml /otelcol/collector-config.yaml
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --chmod=755 --from=build-stage /build/_build /otelcol

ENTRYPOINT ["/otelcol/otelcol-ebpf-profiler"]
CMD ["--config", "/otelcol/collector-config.yaml"]

EXPOSE 4317 4318 12001
