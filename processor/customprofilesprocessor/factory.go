package customprofilesprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/xconsumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper/xprocessorhelper"
	"go.opentelemetry.io/collector/processor/xprocessor"
)

var (
	Type                 = component.MustNewType("customprofilesprocessor")
	consumerCapabilities = consumer.Capabilities{MutatesData: true}
)

const (
	ProfilesStability = component.StabilityLevelDevelopment
)

// NewFactory returns a new factory for the custom profiles processor.
func NewFactory() processor.Factory {
	return xprocessor.NewFactory(
		Type,
		createDefaultConfig,
		xprocessor.WithProfiles(createProfilesProcessor, ProfilesStability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		Foo: "bar",
	}
}

func createProfilesProcessor(
	ctx context.Context,
	params processor.Settings,
	cfg component.Config,
	nextProfilesConsumer xconsumer.Profiles,
) (xprocessor.Profiles, error) {
	return createProfilesProcessorWithOptions(ctx, params, cfg, nextProfilesConsumer)
}

func createProfilesProcessorWithOptions(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	nextProfilesConsumer xconsumer.Profiles,
	options ...option,
) (xprocessor.Profiles, error) {
	kp := createCustomProcessor(set, cfg, options...)

	return xprocessorhelper.NewProfiles(
		ctx,
		set,
		cfg,
		nextProfilesConsumer,
		kp.processProfiles,
		xprocessorhelper.WithCapabilities(consumerCapabilities),
		xprocessorhelper.WithStart(kp.Start),
		xprocessorhelper.WithShutdown(kp.Shutdown),
	)
}

func createCustomProcessor(
	params processor.Settings,
	cfg component.Config,
	options ...option,
) *customprocessor {
	kp := &customprocessor{
		logger:            params.Logger,
		cfg:               cfg,
		options:           options,
		telemetrySettings: params.TelemetrySettings,
	}

	return kp
}

func createProcessorOpts(cfg component.Config) []option {
	oCfg := cfg.(*Config)
	var opts []option

	opts = append(opts,
		withFoo(oCfg.Foo))

	return opts
}
