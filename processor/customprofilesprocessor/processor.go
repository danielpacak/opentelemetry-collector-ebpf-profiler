package customprofilesprocessor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pprofile"
	"go.uber.org/zap"
)

type customprocessor struct {
	cfg               component.Config
	options           []option
	logger            *zap.Logger
	telemetrySettings component.TelemetrySettings
}

func (kp *customprocessor) Start(_ context.Context, _ component.Host) error {
	return nil
}

func (kp *customprocessor) Shutdown(context.Context) error {
	return nil
}

func (kp *customprocessor) processProfiles(_ context.Context, pd pprofile.Profiles) (pprofile.Profiles, error) {
	rp := pd.ResourceProfiles()

	for i := 0; i < rp.Len(); i++ {
		fmt.Printf(">>> Custom processing of profiles with resource attributes: %v\n", rp.At(i).Resource().Attributes().AsRaw())
	}

	return pd, nil
}
