package profilingextensionreceiver

type Config struct {
	// AttachKernelSymbols is an array of kernel symbols to attach the collect_st program to.
	// See /proc/kallsyms for available symbols. For example, wake_up_new_task or copy_process.
	AttachKernelSymbols []string `mapstructure:"attach_kernel_symbols"`
}
