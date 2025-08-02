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
       BpfVerifierLogLevel: 1
       VerboseMode: true
       SendErrorFrames: false
       OffCPUThreshold: 0

   # processors:
   #   custom_profiles_processor:
   #     foo: "bar"

   exporters:
     debug:
       verbosity: normal
     customprofilesexporter:
       export_resource_attributes: true
       export_profile_attributes: true
       export_sample_attributes: true
       export_stack_frames: true
       export_stack_frame_types:
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
   ------------------- New Resource -----------------
     container.id: e5db738b79589d83a5cf6aec3eeb14f6e63c5a15cdba0cb58712a440d523277c
   ---------------------------------------------------
   ------------------- New Profile -------------------
     ProfileID: 00000000000000000000000000000000
     Dropped attributes count: 0
     SampleType: samples
   ------------------- New Sample --------------------
     thread.name: kube-apiserver
     process.executable.name: kube-apiserver
     process.executable.path: /usr/local/bin/kube-apiserver
     process.pid: 2853
     thread.id: 2902
   ---------------------------------------------------
   Instrumentation: go, Function: runtime.usleep, File: runtime/sys_linux_amd64.s, Line: 135, Column: 0
   Instrumentation: go, Function: runtime.sysmon, File: runtime/proc.go, Line: 6063, Column: 0
   Instrumentation: go, Function: runtime.sysmon, File: runtime/proc.go, Line: 6063, Column: 0
   Instrumentation: go, Function: runtime.mstart1, File: runtime/proc.go, Line: 1834, Column: 0
   Instrumentation: go, Function: runtime.mstart0, File: runtime/proc.go, Line: 1800, Column: 0
   Instrumentation: go, Function: runtime.mstart, File: runtime/asm_amd64.s, Line: 396, Column: 0
   ------------------- End Sample --------------------
   ------------------- New Sample --------------------
     thread.name: kube-apiserver
     process.executable.name: kube-apiserver
     process.executable.path: /usr/local/bin/kube-apiserver
     process.pid: 2853
     thread.id: 2919
   ---------------------------------------------------
   Instrumentation: kernel, Function: _raw_spin_unlock_irqrestore, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: try_to_wake_up, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: wake_up_q, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: futex_wake, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: do_futex, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: __x64_sys_futex, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: x64_sys_call, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: do_syscall_64, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: entry_SYSCALL_64_after_hwframe, File: , Line: 0, Column: 0
   Instrumentation: go, Function: runtime.futex, File: runtime/sys_linux_amd64.s, Line: 558, Column: 0
   Instrumentation: go, Function: runtime.futexwakeup, File: runtime/os_linux.go, Line: 82, Column: 0
   Instrumentation: go, Function: runtime.notewakeup, File: runtime/lock_futex.go, Line: 156, Column: 0
   Instrumentation: go, Function: runtime.startm, File: runtime/runtime1.go, Line: 612, Column: 0
   Instrumentation: go, Function: runtime.wakep, File: runtime/runtime1.go, Line: 612, Column: 0
   Instrumentation: go, Function: runtime.resetspinning, File: runtime/proc.go, Line: 3862, Column: 0
   Instrumentation: go, Function: runtime.schedule, File: runtime/proc.go, Line: 4039, Column: 0
   Instrumentation: go, Function: runtime.park_m, File: runtime/proc.go, Line: 4104, Column: 0
   Instrumentation: go, Function: runtime.mcall, File: runtime/asm_amd64.s, Line: 463, Column: 0
   ------------------- End Sample --------------------
   ------------------- New Sample --------------------
     thread.name: kube-apiserver
     process.executable.name: kube-apiserver
     process.executable.path: /usr/local/bin/kube-apiserver
     process.pid: 2853
     thread.id: 2904
   ---------------------------------------------------
   Instrumentation: kernel, Function: finish_task_switch.isra.0, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: __schedule, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: schedule, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: schedule_hrtimeout_range_clock, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: schedule_hrtimeout_range, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: ep_poll, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: do_epoll_wait, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: do_epoll_pwait.part.0, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: __x64_sys_epoll_pwait, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: x64_sys_call, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: do_syscall_64, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: entry_SYSCALL_64_after_hwframe, File: , Line: 0, Column: 0
   Instrumentation: go, Function: internal/runtime/syscall.Syscall6, File: internal/runtime/syscall/asm_linux_amd64.s, Line: 36, Column: 0
   Instrumentation: go, Function: internal/runtime/syscall.EpollWait, File: internal/runtime/syscall/syscall_linux.go, Line: 33, Column: 0
   Instrumentation: go, Function: runtime.netpoll, File: runtime/netpoll_epoll.go, Line: 117, Column: 0
   Instrumentation: go, Function: runtime.findRunnable, File: runtime/proc.go, Line: 3581, Column: 0
   Instrumentation: go, Function: runtime.schedule, File: runtime/proc.go, Line: 3996, Column: 0
   Instrumentation: go, Function: runtime.park_m, File: runtime/proc.go, Line: 4104, Column: 0
   Instrumentation: go, Function: runtime.mcall, File: runtime/asm_amd64.s, Line: 463, Column: 0
   ------------------- End Sample --------------------
   ------------------- New Sample --------------------
     thread.name: kube-apiserver
     process.executable.name: kube-apiserver
     process.executable.path: /usr/local/bin/kube-apiserver
     process.pid: 2853
     thread.id: 2902
   ---------------------------------------------------
   Instrumentation: kernel, Function: _raw_spin_unlock_irqrestore, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: hrtimer_start_range_ns, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: do_nanosleep, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: hrtimer_nanosleep, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: __x64_sys_nanosleep, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: x64_sys_call, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: do_syscall_64, File: , Line: 0, Column: 0
   Instrumentation: kernel, Function: entry_SYSCALL_64_after_hwframe, File: , Line: 0, Column: 0
   Instrumentation: go, Function: runtime.usleep, File: runtime/sys_linux_amd64.s, Line: 135, Column: 0
   Instrumentation: go, Function: runtime.sysmon, File: runtime/proc.go, Line: 6063, Column: 0
   Instrumentation: go, Function: runtime.sysmon, File: runtime/proc.go, Line: 6063, Column: 0
   Instrumentation: go, Function: runtime.mstart1, File: runtime/proc.go, Line: 1834, Column: 0
   Instrumentation: go, Function: runtime.mstart0, File: runtime/proc.go, Line: 1800, Column: 0
   Instrumentation: go, Function: runtime.mstart, File: runtime/asm_amd64.s, Line: 396, Column: 0
   ------------------- End Sample --------------------
   ------------------- End Profile -------------------
   ------------------- End Resource ------------------
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
