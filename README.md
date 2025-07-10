# OpenTelemetry Collector eBPF Profiling Distribution

This OpenTelemetry Collector distribution is made specifically to be used as a node agent to gather
profiles on all processes running on the system. It contains the [eBPF profiler receiver] as well as
a subset of components from OpenTelemetry Collector Core and OpenTelemetry Collector Contrib.

## Quick Start

1. Create a collector configuration file. A very basic configuration may look like this:

    ``` yaml
    # collector-config.yaml
    receivers:
      profiling:
        Tracers: "php,python"
        SamplesPerSecond: 20

    # processors:
    #   custom_profiles_processor:
    #     foo: "bar"

    exporters:
      debug:
        verbosity: normal

    service:
      pipelines:
        profiles:
          receivers:
            - profiling
    #       processors:
    #         - custom_profiles_processor
          exporters:
            - debug
    ```
2. Create and run collector in a new container from the image:

    ```
    docker run --name collector-ebpf-profiling-distro --privileged --pid=host -it --rm \
      -v /sys/kernel/debug:/sys/kernel/debug \
      -v /sys/fs/cgroup:/sys/fs/cgroup \
      -v /proc:/proc \
      -v $PWD/collector-config.yaml:/etc/config.yaml \
      -p 4317:4317 -p 4318:4318 \
      docker.io/danielpacak/opentelemetry-collector-ebpf-profiler:latest \
        --config=/etc/config.yaml \
        --feature-gates=service.profilesSupport
    ```

## Example Kubernetes Deployment

``` mermaid
flowchart LR
  subgraph cluster["Kubernetes Cluster"]
  direction LR
  subgraph nodeA["Kubernetes Node A"]
    collector-ebpf-profiler-a["OTel Collector eBPF Profiling Distro"]
  end
  subgraph nodeB["Kubernetes Node B"]
    collector-ebpf-profiler-b["OTel Collector eBPF Profiling Distro"]
  end
  subgraph nodeC["Kubernetes Node C"]
    collector-ebpf-profiler-c["OTel Collector eBPF Profiling Distro"]
    collector-kubernetes["OTel Collector Kubernetes Distro"]
  end
  end
  pyroscope-development["Pyroscope"]
  pyroscope-backend["Pyroscope"]
  otel-collector["OTel Collector Contrib Distro"]
  collector-ebpf-profiler-a e1@--> pyroscope-development
  collector-ebpf-profiler-a e2@--> otel-collector
  collector-ebpf-profiler-b e3@--> otel-collector
  collector-ebpf-profiler-c e4@--> otel-collector
  otel-collector e5@--> pyroscope-backend
  collector-kubernetes e6@--> otel-collector

  e1@{ animate: true }
  e2@{ animate: true }
  e3@{ animate: true }
  e4@{ animate: true }
  e5@{ animate: true }
  e6@{ animate: true }
```

---

## Building and Running Collector eBPF Profiling Distribution Locally


1. Install the builder. For linux/amd64 platform you can use the following command:

   ```
   curl --proto '=https' --tlsv1.2 -fL -o ocb \
   https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/cmd%2Fbuilder%2Fv0.129.0/ocb_0.129.0_linux_amd64
   chmod +x ocb
   ```

2. Generate the code and build Collector's distribution:

   ```
   ./ocb --config manifest.yaml
   ```

3. Containerize Collector’s distribution:
   1. Enable Docker multi-arch builds:
      ```
      docker run --rm --privileged tonistiigi/binfmt --install all
      docker buildx create --name mybuilder --use
      ```
   2. Build the Docker image as Linux AMD and ARM, and load the build result to "docker images":
      ```
      docker buildx build --load \
        -t docker.io/danielpacak/opentelemetry-collector-ebpf-profiler:latest \
        --platform=linux/amd64,linux/arm64 .
      ```
   3. Test the newly-built image:
      ```
      docker run --name collector-ebpf-profiling-distro --privileged --pid=host -it --rm \
        -v /sys/kernel/debug:/sys/kernel/debug \
        -v /sys/fs/cgroup:/sys/fs/cgroup \
        -v /proc:/proc \
        -v $PWD/collector-config.yaml:/etc/config.yaml \
        -p 4317:4317 -p 4318:4318 \
        docker.io/danielpacak/opentelemetry-collector-ebpf-profiler:latest \
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
