package customprofilesreceiver

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/xconsumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/xreceiver"
)

var (
	typeStr          = component.MustNewType("customprofilesreceiver")
	errInvalidConfig = errors.New("invalid config")
)

// NewFactory creates a factory for the receiver.
func NewFactory() receiver.Factory {
	return xreceiver.NewFactory(
		typeStr,
		defaultConfig,
		xreceiver.WithProfiles(createProfilesReceiver, component.StabilityLevelAlpha))
}

func createProfilesReceiver(
	_ context.Context,
	settings receiver.Settings,
	baseCfg component.Config,
	nextConsumer xconsumer.Profiles) (xreceiver.Profiles, error) {
	cfg, ok := baseCfg.(*Config)
	if !ok {
		return nil, errInvalidConfig
	}

	return NewController(settings.Logger, cfg, nextConsumer), nil
}

func defaultConfig() component.Config {
	return &Config{
		Foo:            "bar",
		ReportInterval: 5 * time.Second,
	}
}
