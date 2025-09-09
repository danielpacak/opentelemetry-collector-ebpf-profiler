package customprofilesprocessor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componentstatus"
	"go.opentelemetry.io/collector/pdata/pprofile"
	"go.uber.org/zap"
)

type customprocessor struct {
	cfg               component.Config
	options           []option
	logger            *zap.Logger
	telemetrySettings component.TelemetrySettings

	foo string
}

func (kp *customprocessor) Start(_ context.Context, host component.Host) error {
	allOptions := append(createProcessorOpts(kp.cfg), kp.options...)

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
	return nil
}

func (kp *customprocessor) processProfiles(_ context.Context, pd pprofile.Profiles) (pprofile.Profiles, error) {
	rp := pd.ResourceProfiles()

	for i := 0; i < rp.Len(); i++ {
		fmt.Printf(">>> Custom processing of profiles with resource attributes: %v\n", rp.At(i).Resource().Attributes().AsRaw())
		rp.At(i).Resource().Attributes().PutStr("foo", kp.foo)
	}

	return pd, nil
}
