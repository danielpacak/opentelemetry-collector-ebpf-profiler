//go:build ignore

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

// A single element map that contains the reference to uprobe__generic program exposed as the interface
// to trigger event-based stack trace collection. This way you can integrate with the opentelemetry-ebpf-profiler
// without modifying its source code.
//
// This map is initialized from user-space by searching for the uprobe__generic program id and
// setting it as value for the element at index 0.
struct {
  __uint(type, BPF_MAP_TYPE_PROG_ARRAY);
  __type(key, __u32);
  __type(value, __u32);
  __uint(max_entries, 1);
} profiler_progs SEC(".maps");

// You can attach this kprobe program to any kernel symbol and trigger stack trace collection on demand.
// Technically, it is similar to off-CPU profiling but usually you'd not attach it somewhere to kernel
// scheduler but someplace that is less frequent, such as process execution.
SEC("kprobe/dummy")
int collect_st(struct pt_regs *ctx)
{
  bpf_tail_call(ctx, &profiler_progs, 0);
  return 0;
}

char __license[] SEC("license") = "Dual MIT/GPL";
