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
		xexporter.WithProfiles(createProfiles, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return &Config{}
}

func createProfiles(ctx context.Context, set exporter.Settings, config component.Config) (xexporter.Profiles, error) {
	customExporter, err := newHTTPRestExporter(set, config)
	if err != nil {
		return nil, err
	}

	return xexporterhelper.NewProfiles(ctx, set, config,
		customExporter.ConsumeProfiles,
		exporterhelper.WithStart(customExporter.Start),
		exporterhelper.WithShutdown(customExporter.Close),
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
	)
}
