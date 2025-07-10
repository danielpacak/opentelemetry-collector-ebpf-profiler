# OpenTelemetry Collector eBPF Profiling Distribution

This OpenTelemetry Collector distribution is made specifically to be used as a node agent to gather
profiles on all processes running on the system. It contains the [eBPF profiler receiver] as well as
a subset of components from OpenTelemetry Collector Core and OpenTelemetry Collector Contrib.

``` yaml
# collector-config.yaml
receivers:
  profiling:
    Tracers: "php,python"
    SamplesPerSecond: 20
exporters:
  debug:
    verbosity: normal
service:
  pipelines:
    profiles:
      receivers:
        - profiling
      exporters:
        - debug
```

```
docker run --name collector-ebpf-profiler --privileged --pid=host -it --rm \
  -v /sys/kernel/debug:/sys/kernel/debug \
  -v /sys/fs/cgroup:/sys/fs/cgroup \
  -v /proc:/proc \
  -v $PWD/collector-config.yaml:/etc/config.yaml \
  -p 4317:4317 -p 4318:4318 \
  danielpacak/opentelemetry-collector-ebpf-profiler:latest \
  --config=/etc/config.yaml \
  --feature-gates=service.profilesSupport
```

``` mermaid
flowchart LR
  subgraph cluster["Kubernetes Cluster"]
  direction LR
  subgraph nodeA["Kubernetes Node A"]
    collector-ebpf-profiler-a["Collector eBPF Profiler"]
  end
  subgraph nodeB["Kubernetes Node B"]
    collector-ebpf-profiler-b["Collector eBPF Profiler"]
  end
  subgraph nodeC["Kubernetes Node C"]
    collector-ebpf-profiler-c["Collector eBPF Profiler"]
  end
  end
  pyroscope-development["Pyroscope"]
  pyroscope-backend["Pyroscope"]
  otel-collector["OTel Collector"]
  collector-ebpf-profiler-a --> pyroscope-development
  collector-ebpf-profiler-a --> otel-collector
  collector-ebpf-profiler-b --> otel-collector
  collector-ebpf-profiler-c --> otel-collector
  otel-collector --> pyroscope-backend
```

---

manifest.yaml - builder manifest

Install the builder:

```
curl --proto '=https' --tlsv1.2 -fL -o ocb \
https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/cmd%2Fbuilder%2Fv0.129.0/ocb_0.129.0_linux_amd64
chmod +x ocb
```

Generate the code and build your collector's distribution:

```
./ocb --config manifest.yaml
```

Containerize your Collector’s distribution:

```
# Enable Docker multi-arch builds
docker run --rm --privileged tonistiigi/binfmt --install all
docker buildx create --name mybuilder --use
```

```
# Build the Docker image as Linux AMD and ARM,
# and loads the build result to "docker images"
docker buildx build --load \
  -t danielpacak/opentelemetry-collector-ebpf-profiler:latest \
  --platform=linux/amd64,linux/arm64 .
```

```
# Test the newly-built image
docker run --name collector-ebpf-profiler \
  --privileged \
  --pid=host \
  -it \
  --rm \
  -v /sys/kernel/debug:/sys/kernel/debug \
  -v /sys/fs/cgroup:/sys/fs/cgroup \
  -v /proc:/proc \
  -v $PWD/collector-config.yaml:/etc/config.yaml \
  --publish=4317:4317 \
  --publish=4318:4318 \
  danielpacak/opentelemetry-collector-ebpf-profiler:latest \
  --config=/etc/config.yaml \
  --feature-gates=service.profilesSupport
```

## Resources

1. https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-ebpf-profiler
2. https://github.com/open-telemetry/opentelemetry-ebpf-profiler/issues/521
3. https://opentelemetry.io/docs/collector/custom-collector/
4. https://blog.jaimyn.dev/how-to-build-multi-architecture-docker-images-on-an-m1-mac/
5. https://github.com/grafana/pyroscope/blob/main/examples/grafana-alloy-auto-instrumentation/ebpf-otel/docker/docker-compose.yml

[eBPF profiler receiver]: https://github.com/open-telemetry/opentelemetry-ebpf-profiler
