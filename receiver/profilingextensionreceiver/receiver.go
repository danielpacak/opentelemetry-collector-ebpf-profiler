package profilingextensionreceiver

import (
	"context"
	"errors"
	"fmt"
	"syscall"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/xconsumer"
	"go.uber.org/zap"
)

type extensionReceiver struct {
	logger       *zap.Logger
	config       *Config
	nextConsumer xconsumer.Profiles

	objs  extensionObjects
	links []link.Link
}

func newExtensionReceiver(logger *zap.Logger, config *Config, nextConsumer xconsumer.Profiles) *extensionReceiver {
	return &extensionReceiver{
		logger:       logger,
		config:       config,
		nextConsumer: nextConsumer,
	}
}

// do not produce any new profiles, just load and attach eBPF objects and tail call uprobe__generic program
func (e *extensionReceiver) Start(_ context.Context, _ component.Host) error {
	e.logger.Info("Starting profilingextension receiver")

	// Remove resource limits for kernels <5.11.
	err := rlimit.RemoveMemlock()
	if err != nil {
		return fmt.Errorf("removing memlock: %w", err)
	}

	// Load the compiled eBPF ELF and load it into the kernel.
	err = loadExtensionObjects(&e.objs, nil)
	if err != nil {
		return fmt.Errorf("loading eBPF objects: %w", err)
	}

	uprobeGenericProg, err := e.findUprobeGenericProg()
	if err != nil {
		return err
	}
	uprobeGenericProgInfo, err := uprobeGenericProg.Info()
	if err != nil {
		return err
	}
	uprobeGenericProgID := uint32(0)
	if id, ok := uprobeGenericProgInfo.ID(); ok {
		uprobeGenericProgID = uint32(id)
	}

	e.logger.Info("Updating eBPF prog array map", zap.Uint32("prog.id", uprobeGenericProgID))
	err = e.objs.ProfilerProgs.Update(uint32(0), uprobeGenericProg, ebpf.UpdateAny)
	if err != nil {
		return err
	}

	// Attach collect_st to configured symbols
	for _, symbol := range e.config.AttachKernelSymbols {
		e.logger.Info("Attaching eBPF program to perf event",
			zap.String("program", "collect_st"),
			zap.String("symbol", symbol))
		link, err := link.Kprobe(symbol, e.objs.CollectSt, nil)
		if err != nil {
			return fmt.Errorf("attaching collect_st to perf event: %s: %w", symbol, err)
		}
		e.links = append(e.links, link)
	}

	return nil
}

func (e *extensionReceiver) Shutdown(_ context.Context) error {
	e.logger.Info("Shutting down profilingextension receiver")
	for _, link := range e.links {
		_ = link.Close()
	}
	_ = e.objs.Close()
	return nil
}

var ErrProgramNotFound = errors.New("program not found")

func (e *extensionReceiver) findUprobeGenericProg() (*ebpf.Program, error) {
	var err error
	var id ebpf.ProgramID = 0
	var prog *ebpf.Program
	var progInfo *ebpf.ProgramInfo

	for {
		id, err = ebpf.ProgramGetNextID(id)
		if err != nil {
			if errors.Is(err, syscall.ENOENT) {
				return nil, ErrProgramNotFound
			}
			return nil, err
		}

		prog, err = ebpf.NewProgramFromID(id)
		if err != nil {
			return nil, err
		}

		progInfo, err = prog.Info()
		if err != nil {
			return nil, err
		}

		if progInfo.Type == ebpf.Kprobe && progInfo.Name == "uprobe__generic" {
			return prog, nil
		}
	}
}
