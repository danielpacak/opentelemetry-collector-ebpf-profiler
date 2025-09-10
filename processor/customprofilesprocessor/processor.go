package customprofilesprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componentstatus"
	"go.opentelemetry.io/collector/pdata/pprofile"
	"go.uber.org/zap"
)

type customprocessor struct {
	logger            *zap.Logger
	config            component.Config
	options           []option
	telemetrySettings component.TelemetrySettings

	foo string
}

func (kp *customprocessor) Start(_ context.Context, host component.Host) error {
	kp.logger.Info("Starting custom profiles processor")
	allOptions := append(createProcessorOpts(kp.config), kp.options...)

	for _, opt := range allOptions {
		if err := opt(kp); err != nil {
			kp.logger.Error("Could not apply option", zap.Error(err))
			componentstatus.ReportStatus(host, componentstatus.NewFatalErrorEvent(err))
			return err
		}
	}

	return nil
}

func (kp *customprocessor) Shutdown(context.Context) error {
	kp.logger.Info("Shutting down custom profiles processor")
	return nil
}

func (kp *customprocessor) processProfiles(_ context.Context, pd pprofile.Profiles) (pprofile.Profiles, error) {
	rp := pd.ResourceProfiles()

	for i := 0; i < rp.Len(); i++ {
		kp.logger.Info("Adding custom resource attribute",
			zap.Any("attributes", rp.At(i).Resource().Attributes().AsRaw()))
		rp.At(i).Resource().Attributes().PutStr("foo", kp.foo)
	}

	return pd, nil
}
