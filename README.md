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
       Tracers: "perl,php,python,hotspot,ruby,v8,dotnet,go"
       SamplesPerSecond: 20

   # processors:
   #   custom_profiles_processor:
   #     foo: "bar"

   exporters:
     debug:
       verbosity: normal
     customprofilesexporter:
       export_sample_attributes: true
       export_unwind_types:
         - native
         - kernel
         - go
         - jvm
         - php
         - cpython

   service:
     pipelines:
       profiles:
         receivers:
           - profiling
   #     processors:
   #       - custom_profiles_processor
         exporters:
           - debug
           - customprofilesexporter
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
   ```
   ------------------- New Sample -------------------
     container.id: /../docker-d4f73fc52bbd83e58259256a0b321925bd1c43d0900904d4ed7c6113d1f9c4a4.scope/init
     thread.name: buildkitd
     process.executable.name: buildkitd
     process.executable.path: /usr/bin/buildkitd
     process.pid: 1921
     thread.id: 2072
   ---------------------------------------------------
   Instrumentation: kernel, Function: link_path_walk.part.0.constprop.0, File: , Line: 0
   Instrumentation: kernel, Function: path_lookupat, File: , Line: 0
   Instrumentation: kernel, Function: filename_lookup, File: , Line: 0
   Instrumentation: kernel, Function: vfs_statx, File: , Line: 0
   Instrumentation: kernel, Function: vfs_fstatat, File: , Line: 0
   Instrumentation: kernel, Function: __do_sys_newfstatat, File: , Line: 0
   Instrumentation: kernel, Function: __x64_sys_newfstatat, File: , Line: 0
   Instrumentation: kernel, Function: x64_sys_call, File: , Line: 0
   Instrumentation: kernel, Function: do_syscall_64, File: , Line: 0
   Instrumentation: kernel, Function: entry_SYSCALL_64_after_hwframe, File: , Line: 0
   Instrumentation: go, Function: internal/runtime/syscall.Syscall6, File: /usr/local/go/src/internal/runtime/syscall/asm_linux_amd64.s, Line: 36
   Instrumentation: go, Function: syscall.RawSyscall6, File: /usr/local/go/src/syscall/syscall_linux.go, Line: 66
   Instrumentation: go, Function: syscall.Syscall6, File: /usr/local/go/src/syscall/syscall_linux.go, Line: 96
   Instrumentation: go, Function: syscall.fstatat, File: /usr/local/go/src/syscall/zsyscall_linux_amd64.go, Line: 1438
   Instrumentation: go, Function: os.lstatNolog, File: /usr/local/go/src/syscall/syscall_linux_amd64.go, Line: 69
   ```

## Example Kubernetes Deployment

```
kubectl apply -f example/kubernetes/node-agent.yaml
```

``` mermaid
flowchart LR
  subgraph nodeA["Node: A"]
    collector-ebpf-profiler-a["Pod: OTel Collector eBPF Profiling Distro"]
  end
  subgraph nodeB["Node: B"]
    collector-ebpf-profiler-b["Pod: OTel Collector eBPF Profiling Distro"]
  end
  subgraph nodeC["Node: C"]
    collector-ebpf-profiler-c["Pod: OTel Collector eBPF Profiling Distro"]
    collector-kubernetes["Pod: OTel Collector Kubernetes Distro"]
  end
  pyroscope["Pyroscope"]
  clickhouse@{ shape: cyl, label: "ClickHouse" }
  otel-collector["OTel Collector Contrib Distro"]
  collector-ebpf-profiler-a e2@--> otel-collector
  collector-ebpf-profiler-b e3@--> otel-collector
  collector-ebpf-profiler-c e4@--> otel-collector
  otel-collector e5@--> pyroscope
  otel-collector e1@--> clickhouse
  collector-kubernetes e6@--> otel-collector

  e1@{ animate: true }
  e2@{ animate: true }
  e3@{ animate: true }
  e4@{ animate: true }
  e5@{ animate: true }
  e6@{ animate: true }
```

``` mermaid
flowchart TD
  subgraph node-a["Node: A"]
    subgraph pod["Pod: collector-ebpf-profiler"]
      subgraph container["Container: profiler"]
        procfs@{ shape: cyl, label: "/proc" }
        cgroupfs@{ shape: cyl, label: "/sys/fs/cgroup" }
        debugfs@{ shape: cyl, label: "/sys/kernel/debug" }
      end
    end
    volumeMount1["/proc"]
    volumeMount2["/sys/fs/cgroup"]
    volumeMount3["/sys/kernel/debug"]

    procfs --> volumeMount1
    cgroupfs --> volumeMount2
    debugfs --> volumeMount3

  end
```

## Example Docker Compose Deployment

```
cd example/docker
docker compose up -d
```

Pyroscope is accessible at http://localhost:4040 and Grafana at http://localhost:3000. Grafana is
provisioned with the Pyroscope datasource so you can either see profiles in Pyroscope web UI or with
Grafana's Pyroscope application.

```
docker compose down
```

---

## Building and Running Locally


1. Install the builder. For linux/amd64 platform you can use the following command:

   ```
   curl --proto '=https' --tlsv1.2 -fL -o ocb \
   https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/cmd%2Fbuilder%2Fv0.131.0/ocb_0.131.0_linux_amd64
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

## Further Reading

1. https://opentelemetry.io/docs/collector/distributions/
2. https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-ebpf-profiler
3. https://github.com/open-telemetry/opentelemetry-ebpf-profiler/issues/521
4. https://opentelemetry.io/docs/collector/custom-collector/
5. https://blog.jaimyn.dev/how-to-build-multi-architecture-docker-images-on-an-m1-mac/
6. https://github.com/grafana/pyroscope/blob/main/examples/grafana-alloy-auto-instrumentation/ebpf-otel/docker/docker-compose.yml
7. https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/k8sattributesprocessor/README.md

[eBPF profiler receiver]: https://github.com/open-telemetry/opentelemetry-ebpf-profiler
