# profilingextensionreceiver

## Troubleshooting

``` console
$ go generate
In file included from /home/dpacak/go/src/github.com/danielpacak/opentelemetry-collector-ebpf-profiler/receiver/profilingextensionreceiver/extension.c:3:
In file included from /usr/include/linux/bpf.h:11:
/usr/include/linux/types.h:5:10: fatal error: 'asm/types.h' file not found
#include <asm/types.h>
         ^~~~~~~~~~~~~
1 error generated.
Error: compile: exit status 1
gen.go:3: running "go": exit status 1
```

```
sudo ln -s /usr/include/x86_64-linux-gnu/asm /usr/include/asm
```
