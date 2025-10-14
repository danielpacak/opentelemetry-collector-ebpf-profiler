package customprofilesreceiver

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/xconsumer"
	"go.opentelemetry.io/collector/pdata/pprofile"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.uber.org/zap"
)

func NewController(logger *zap.Logger, config *Config, nextConsumer xconsumer.Profiles) *Controller {
	return &Controller{
		logger:       logger,
		config:       config,
		nextConsumer: nextConsumer,
	}
}

type Controller struct {
	logger       *zap.Logger
	config       *Config
	nextConsumer xconsumer.Profiles
}

func (c *Controller) Start(ctx context.Context, _ component.Host) error {
	c.logger.Info("Starting custom profiles receiver")
	go func() {
		tick := time.NewTicker(c.config.ReportInterval)
		defer tick.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				c.logger.Info(">>> GENERATING REPORT <<<")
				profiles, err := c.GenerateProfiles()
				if err != nil {
					// TODO log error
					c.logger.Error(err.Error())
				}
				err = c.nextConsumer.ConsumeProfiles(ctx, profiles)
				if err != nil {
					// TODO log error
					c.logger.Error(err.Error())
				}
			}
		}
	}()

	return nil
}

func (c *Controller) GenerateProfiles() (pprofile.Profiles, error) {
	profiles := pprofile.NewProfiles()
	rp := profiles.ResourceProfiles().AppendEmpty()
	rp.Resource().Attributes().PutStr(string(semconv.ContainerIDKey),
		"abc123...")
	rp.SetSchemaUrl(semconv.SchemaURL)

	sp := rp.ScopeProfiles().AppendEmpty()
	sp.SetSchemaUrl(semconv.SchemaURL)

	return profiles, nil
}

func (c *Controller) Shutdown(_ context.Context) error {
	c.logger.Info("Shutting down custom profiles receiver")
	return nil
}
