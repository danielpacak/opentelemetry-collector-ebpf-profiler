package customprofilesexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/exporter/exporterhelper/xexporterhelper"
	"go.opentelemetry.io/collector/exporter/xexporter"
)

func NewFactory() exporter.Factory {
	return xexporter.NewFactory(
		component.MustNewType("customprofilesexporter"),
		createDefaultConfig,
		xexporter.WithProfiles(createProfiles, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return &customProfilesExporterConfig{
		ExportSampleAttributes: true,
	}
}

func createProfiles(ctx context.Context, set exporter.Settings, config component.Config) (xexporter.Profiles, error) {
	customExporter := &customProfilesExporter{
		logger: set.Logger,
		config: config.(*customProfilesExporterConfig),
	}

	return xexporterhelper.NewProfilesExporter(ctx, set, config,
		customExporter.ConsumeProfiles,
		exporterhelper.WithStart(customExporter.Start),
		exporterhelper.WithShutdown(customExporter.Close),
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
	)
}
