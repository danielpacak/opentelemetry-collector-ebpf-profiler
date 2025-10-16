package httprestexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/exporter/exporterhelper/xexporterhelper"
	"go.opentelemetry.io/collector/exporter/xexporter"
)

var (
	strType = component.MustNewType("httprest")
)

func NewFactory() exporter.Factory {
	return xexporter.NewFactory(
		strType,
		createDefaultConfig,
		xexporter.WithProfiles(createProfilesFunc, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		Address: ":7799",
		ExportSampleTypes: []string{
			"samples",
			"events",
		},
	}
}

func createProfilesFunc(ctx context.Context, set exporter.Settings, config component.Config) (xexporter.Profiles, error) {
	restExporter, err := newHTTPRestExporter(set, config)
	if err != nil {
		return nil, err
	}

	return xexporterhelper.NewProfiles(ctx, set, config,
		restExporter.ConsumeProfiles,
		exporterhelper.WithStart(restExporter.Start),
		exporterhelper.WithShutdown(restExporter.Close),
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
	)
}
