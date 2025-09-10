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

var (
	strType = component.MustNewType("customprofilesexporter")
)

func NewFactory() exporter.Factory {
	return xexporter.NewFactory(
		strType,
		createDefaultConfig,
		xexporter.WithProfiles(createProfiles, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		ExportResourceAttributes:         true,
		ExportProfileAttributes:          true,
		ExportSampleAttributes:           true,
		ExportStackFrames:                true,
		IgnoreProfilesWithoutContainerID: true,
	}
}

func createProfiles(ctx context.Context, set exporter.Settings, config component.Config) (xexporter.Profiles, error) {
	customExporter := &customexporter{
		logger: set.Logger,
		config: config.(*Config),
	}

	return xexporterhelper.NewProfiles(ctx, set, config,
		customExporter.ConsumeProfiles,
		exporterhelper.WithStart(customExporter.Start),
		exporterhelper.WithShutdown(customExporter.Close),
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
	)
}
