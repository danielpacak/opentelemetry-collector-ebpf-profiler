package profilingextensionreceiver

import (
	"context"
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/xconsumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/xreceiver"
)

var (
	typeStr          = component.MustNewType("profilingextension")
	errInvalidConfig = errors.New("invalid config")
)

func NewFactory() receiver.Factory {
	return xreceiver.NewFactory(
		typeStr,
		defaultConfig,
		xreceiver.WithProfiles(createProfilesFunc, component.StabilityLevelAlpha))
}

func defaultConfig() component.Config {
	return &Config{
		AttachKernelSymbols: []string{
			"copy_process",
		},
	}
}

func createProfilesFunc(ctx context.Context, settings receiver.Settings, config component.Config, nextConsumer xconsumer.Profiles) (xreceiver.Profiles, error) {
	cfg, ok := config.(*Config)
	if !ok {
		return nil, errInvalidConfig
	}

	return newExtensionReceiver(settings.Logger, cfg, nextConsumer), nil
}
