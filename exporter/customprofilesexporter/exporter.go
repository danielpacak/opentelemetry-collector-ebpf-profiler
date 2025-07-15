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
	mappingTable := pd.ProfilesDictionary().MappingTable()
	locationTable := pd.ProfilesDictionary().LocationTable()
	attributeTable := pd.ProfilesDictionary().AttributeTable()
	functionTable := pd.ProfilesDictionary().FunctionTable()
	stringTable := pd.ProfilesDictionary().StringTable()

	rps := pd.ResourceProfiles()
	for i := 0; i < rps.Len(); i++ {
		rp := rps.At(i)

		sps := rp.ScopeProfiles()
		for j := 0; j < sps.Len(); j++ {
			pcs := sps.At(j).Profiles()
			for k := 0; k < pcs.Len(); k++ {
				profile := pcs.At(k)
				profileLocationsIndices := profile.LocationIndices()

				samples := profile.Sample()

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
							attr := attributeTable.At(int(sampleAttrs.At(n)))
							fmt.Printf("  %s: %s\n", attr.Key(), attr.Value().AsString())
						}
						fmt.Println("---------------------------------------------------")
					}

					for m := sample.LocationsStartIndex(); m < sample.LocationsStartIndex()+sample.LocationsLength(); m++ {
						location := locationTable.At(int(profileLocationsIndices.At(int(m))))
						locationAttrs := location.AttributeIndices()

						unwindType := "unknown"
						for la := 0; la < locationAttrs.Len(); la++ {
							attr := attributeTable.At(int(locationAttrs.At(la)))
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
							filename := "<unknown>"
							if location.HasMappingIndex() {
								mapping := mappingTable.At(int(location.MappingIndex()))
								filename = stringTable.At(int(mapping.FilenameStrindex()))
							}
							fmt.Printf("Instrumentation: %s: Function: %#04x, File: %s\n", unwindType, location.Address(), filename)
						}

						for n := 0; n < locationLine.Len(); n++ {
							line := locationLine.At(n)
							function := functionTable.At(int(line.FunctionIndex()))
							functionName := stringTable.At(int(function.NameStrindex()))
							fileName := stringTable.At(int(function.FilenameStrindex()))
							fmt.Printf("Instrumentation: %s, Function: %s, File: %s, Line: %d, Column: %d\n",
								unwindType, functionName, fileName, line.Line(), line.Column())
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
