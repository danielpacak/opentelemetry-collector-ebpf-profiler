package httprestexporter

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/pprofile"
	"go.uber.org/zap"
)

type customexporter struct {
	logger       *zap.Logger
	config       *Config
	shutdownFunc func(ctx context.Context) error
	store        *store
}

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

func newHTTPRestExporter(set exporter.Settings, config component.Config) (*customexporter, error) {
	return &customexporter{
		logger:       set.Logger,
		config:       config.(*Config),
		shutdownFunc: func(_ context.Context) error { return nil },
		store:        newStore(),
	}, nil
}

func (e *customexporter) Start(_ context.Context, _ component.Host) error {
	e.logger.Info("Starting custom profiles exporter...", zap.Any("config", e.config))

	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/workloads", e.getWorkloadsHandler)
	router.HandlerFunc(http.MethodGet, "/workload/:id/functions", e.getWorkloadFunctionsHandler)

	srv := &http.Server{Addr: e.config.Address, Handler: router}
	ln, err := net.Listen("tcp", e.config.Address)
	if err != nil {
		return err
	}

	e.shutdownFunc = func(ctx context.Context) error {
		return srv.Shutdown(ctx)
	}
	go func() {
		_ = srv.Serve(ln)
	}()

	return nil
}

func (e *customexporter) ConsumeProfiles(_ context.Context, pd pprofile.Profiles) error {
	mappingTable := pd.Dictionary().MappingTable()
	locationTable := pd.Dictionary().LocationTable()
	attributeTable := pd.Dictionary().AttributeTable()
	functionTable := pd.Dictionary().FunctionTable()
	stringTable := pd.Dictionary().StringTable()
	stackTable := pd.Dictionary().StackTable()

	rps := pd.ResourceProfiles()
	for i := 0; i < rps.Len(); i++ {
		rp := rps.At(i)
		rAttributes := rp.Resource().Attributes()

		containerID, ok := rAttributes.Get("container.id")
		if !ok || containerID.AsString() == "" {
			continue
		}
		container := NewContainer(containerID.AsString())

		if podName, ok := rAttributes.Get("k8s.pod.name"); ok {
			container.PodName = podName.AsString()
		}

		if namespaceName, ok := rAttributes.Get("k8s.namespace.name"); ok {
			container.NamespaceName = namespaceName.AsString()
		}

		if containerName, ok := rAttributes.Get("k8s.container.name"); ok {
			container.ContainerName = containerName.AsString()
		}

		container = e.store.UpsertContainer(container)

		sps := rp.ScopeProfiles()
		for j := 0; j < sps.Len(); j++ {
			pcs := sps.At(j).Profiles()
			for k := 0; k < pcs.Len(); k++ {
				profile := pcs.At(k)

				samples := profile.Sample()

				for l := 0; l < samples.Len(); l++ {
					sample := samples.At(l)

					process := NewProcess()

					sampleAttrs := sample.AttributeIndices()
					for n := 0; n < sampleAttrs.Len(); n++ {
						attr := attributeTable.At(int(sampleAttrs.At(n)))
						attrKey := stringTable.At(int(attr.KeyStrindex()))
						if "process.pid" == attrKey {
							process.PID = int(attr.Value().Int())
						}
						if "process.executable.name" == attrKey {
							process.ExecutableName = attr.Value().AsString()
						}
						if "process.executable.path" == attrKey {
							process.ExecutablePath = attr.Value().AsString()
						}
					}

					process = container.UpsertProcess(process)

					stack := stackTable.At(int(sample.StackIndex()))
					profileLocationsIndices := stack.LocationIndices()

					for m := 0; m < profileLocationsIndices.Len(); m++ {
						location := locationTable.At(int(profileLocationsIndices.At(m)))
						locationAttrs := location.AttributeIndices()

						unwindType := "unknown"
						for la := 0; la < locationAttrs.Len(); la++ {
							attr := attributeTable.At(int(locationAttrs.At(la)))
							if stringTable.At(int(attr.KeyStrindex())) == "profile.frame.type" {
								unwindType = attr.Value().AsString()
								break
							}
						}

						locationLine := location.Line()
						if locationLine.Len() == 0 {
							filename := "<unknown>"
							if location.MappingIndex() > 0 {
								mapping := mappingTable.At(int(location.MappingIndex()))
								filename = stringTable.At(int(mapping.FilenameStrindex()))
							}
							fmt.Printf("Instrumentation: %s: Function: %#04x, File: %s\n", unwindType, location.Address(), filename)
						}

						for n := 0; n < locationLine.Len(); n++ {
							line := locationLine.At(n)
							function := functionTable.At(int(line.FunctionIndex()))
							functionName := stringTable.At(int(function.NameStrindex()))
							fileName := stringTable.At(int(function.FilenameStrindex()))
							process.UpsertFunction(Function{
								Language: unwindType,
								Name:     functionName,
								File:     fileName,
							})
						}

					}
				}
			}
		}
	}

	return nil
}

func (e *customexporter) Close(ctx context.Context) error {
	e.logger.Info("Closing custom profiles exporter...")
	return e.shutdownFunc(ctx)
}

// HTTP REST API Handlers
func (e *customexporter) getWorkloadsHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(e.store.containers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (e *customexporter) getWorkloadFunctionsHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	containerID := params.ByName("id")
	if "" == containerID {
		http.NotFound(w, r)
		return
	}

	container, ok := e.store.containers[containerID]
	if !ok {
		http.NotFound(w, r)
		return
	}

	var functions []Function

	for _, process := range container.Processes {
		for function := range process.functions {
			functions = append(functions, function)
		}
	}

	js, err := json.Marshal(functions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
