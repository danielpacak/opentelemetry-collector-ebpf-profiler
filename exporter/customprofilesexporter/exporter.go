package customprofilesexporter

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pprofile"
	"go.uber.org/zap"
)

type customProfilesExporterConfig struct {
	Foo                    string   `mapstructure:"foo,omitempty"`
	ExportSampleAttributes bool     `mapstructure:"export_sample_attributes"`
	ExportUnwindTypes      []string `mapstructure:"export_unwind_types"`
}

type customProfilesExporter struct {
	logger *zap.Logger
	config *customProfilesExporterConfig
}

func (e *customProfilesExporter) Start(_ context.Context, _ component.Host) error {
	e.logger.Info("Starting custom profiles exporter...")
	return nil
}

func (e *customProfilesExporter) ConsumeProfiles(_ context.Context, pd pprofile.Profiles) error {
	rps := pd.ResourceProfiles()
	for i := 0; i < rps.Len(); i++ {
		rp := rps.At(i)

		sps := rp.ScopeProfiles()
		for j := 0; j < sps.Len(); j++ {
			pcs := sps.At(j).Profiles()
			for k := 0; k < pcs.Len(); k++ {
				profile := pcs.At(k)
				profileLocationsIndices := profile.LocationIndices()

				locations := pd.ProfilesDictionary().LocationTable()
				attributesTable := pd.ProfilesDictionary().AttributeTable()
				functions := pd.ProfilesDictionary().FunctionTable()
				samples := profile.Sample()
				stringTable := pd.ProfilesDictionary().StringTable()

				// print the type of the sample
				sampleType := "samples"
				for n := 0; n < profile.SampleType().Len(); n++ {
					sampleType = stringTable.At(int(profile.SampleType().At(n).TypeStrindex()))
					fmt.Println("SampleType: ", sampleType)
				}
				fmt.Println("------------------- New Profile -------------------")
				fmt.Println("Dropped attributes count", strconv.FormatUint(uint64(profile.DroppedAttributesCount()), 10))

				for l := 0; l < samples.Len(); l++ {
					sample := samples.At(l)

					fmt.Println("------------------- New Sample -------------------")
					if e.config.ExportSampleAttributes {
						sampleAttrs := sample.AttributeIndices()
						for n := 0; n < sampleAttrs.Len(); n++ {
							attr := attributesTable.At(int(sampleAttrs.At(n)))
							fmt.Printf("  %s: %s\n", attr.Key(), attr.Value().AsString())
						}
						fmt.Println("---------------------------------------------------")
					}

					for m := sample.LocationsStartIndex(); m < sample.LocationsStartIndex()+sample.LocationsLength(); m++ {
						location := locations.At(int(profileLocationsIndices.At(int(m))))
						locationAttrs := location.AttributeIndices()

						unwindType := "unknown"
						for la := 0; la < locationAttrs.Len(); la++ {
							attr := attributesTable.At(int(locationAttrs.At(la)))
							if attr.Key() == "profile.frame.type" {
								unwindType = attr.Value().AsString()
								break
							}
						}

						if len(e.config.ExportUnwindTypes) > 0 &&
							!slices.Contains(e.config.ExportUnwindTypes, unwindType) {
							continue
						}

						locationLine := location.Line()

						if locationLine.Len() == 0 {
							fmt.Printf("??? Instrumentation: %s ???\n", unwindType)
						}

						for n := 0; n < locationLine.Len(); n++ {

							line := locationLine.At(n)

							function := functions.At(int(line.FunctionIndex()))
							lineNumber := line.Line()
							functionName := stringTable.At(int(function.NameStrindex()))
							fileName := stringTable.At(int(function.FilenameStrindex()))
							fmt.Printf("Instrumentation: %s, Function: %s, File: %s, Line: %d\n",
								unwindType, functionName, fileName, lineNumber)
						}
					}
					fmt.Println("------------------- End New Sample -------------------")
				}
				fmt.Println("------------------- End of Profile -------------------")
			}
		}
	}
	return nil
}

func (e *customProfilesExporter) Close(_ context.Context) error {
	e.logger.Info("Closing custom profiles exporter...")
	return nil
}
