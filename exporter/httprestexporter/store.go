package httprestexporter

type Function struct {
	Language string
	Name     string
	File     string
}

type Process struct {
	PID            int
	ExecutableName string
	ExecutablePath string

	functions map[Function]bool
}

func NewProcess() *Process {
	return &Process{
		functions: make(map[Function]bool),
	}
}

type Container struct {
	ContainerID string

	NamespaceName string
	PodName       string
	ContainerName string

	Processes map[int]*Process
}

func NewContainer(containerID string) *Container {
	return &Container{
		ContainerID: containerID,
		Processes:   make(map[int]*Process),
	}
}

func (c *Container) UpsertProcess(process *Process) *Process {
	ref, ok := c.Processes[process.PID]
	if ok {
		return ref
	}
	c.Processes[process.PID] = process
	return process
}

func (c *Process) UpsertFunction(function Function) {
	_, ok := c.functions[function]
	if ok {
		return
	}
	c.functions[function] = true
}

type store struct {
	containers map[string]*Container
}

func newStore() *store {
	return &store{
		containers: make(map[string]*Container),
	}
}

func (s *store) UpsertContainer(container *Container) *Container {
	ref, ok := s.containers[container.ContainerID]
	if ok {
		return ref
	}
	s.containers[container.ContainerID] = container
	return container
}
