package customprofilesexporter

import (
	"context"
	"fmt"
	"slices"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pprofile"
	"go.uber.org/zap"
)

type customexporter struct {
	logger *zap.Logger
	config *Config
}

func (e *customexporter) Start(_ context.Context, _ component.Host) error {
	e.logger.Info("Starting custom profiles exporter...", zap.Any("config", e.config))
	return nil
}

func (e *customexporter) ConsumeProfiles(_ context.Context, pd pprofile.Profiles) error {
	mappingTable := pd.Dictionary().MappingTable()
	locationTable := pd.Dictionary().LocationTable()
	attributeTable := pd.Dictionary().AttributeTable()
	functionTable := pd.Dictionary().FunctionTable()
	stringTable := pd.Dictionary().StringTable()

	rps := pd.ResourceProfiles()
	for i := 0; i < rps.Len(); i++ {
		rp := rps.At(i)

		if e.config.IgnoreProfilesWithoutContainerID {
			containerID, ok := rp.Resource().Attributes().Get("container.id")
			if !ok || containerID.AsString() == "" {
				fmt.Println("--------------- New Resource Profile --------------")
				fmt.Println("              SKIPPED (no container.id)")
				fmt.Printf("-------------- End Resource Profile ---------------\n\n")
				continue
			}
		}

		fmt.Println("--------------- New Resource Profile --------------")
		if e.config.ExportResourceAttributes {
			if rp.Resource().Attributes().Len() > 0 {
				rp.Resource().Attributes().Range(func(k string, v pcommon.Value) bool {
					fmt.Printf("  %s: %s (%s)\n", k, v.AsString(), v.Type().String())
					return true
				})
			}
		}

		sps := rp.ScopeProfiles()
		for j := 0; j < sps.Len(); j++ {
			pcs := sps.At(j).Profiles()
			for k := 0; k < pcs.Len(); k++ {
				profile := pcs.At(k)

				fmt.Println("------------------- New Profile -------------------")
				fmt.Printf("  ProfileID: %x\n", [16]byte(profile.ProfileID()))
				fmt.Printf("  Dropped attributes count: %d\n", profile.DroppedAttributesCount())
				sampleType := "samples"
				for n := 0; n < profile.SampleType().Len(); n++ {
					sampleType = stringTable.At(int(profile.SampleType().At(n).TypeStrindex()))
					fmt.Printf("  SampleType: %s\n", sampleType)
				}
				profileAttrs := profile.AttributeIndices()
				if profileAttrs.Len() > 0 {
					for n := 0; n < profileAttrs.Len(); n++ {
						attr := attributeTable.At(int(profileAttrs.At(n)))
						fmt.Printf("  %s: %s (%s)\n", attr.Key(), attr.Value().AsString(), attr.Value().Type().String())
					}
					fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
				}

				samples := profile.Sample()

				for l := 0; l < samples.Len(); l++ {
					sample := samples.At(l)

					fmt.Println("------------------- New Sample --------------------")
					if e.config.ExportSampleAttributes {
						sampleAttrs := sample.AttributeIndices()
						for n := 0; n < sampleAttrs.Len(); n++ {
							attr := attributeTable.At(int(sampleAttrs.At(n)))
							fmt.Printf("  %s: %s (%s)\n", attr.Key(), attr.Value().AsString(), attr.Value().Type().String())
						}
						fmt.Println("---------------------------------------------------")
					}

					profileLocationsIndices := profile.LocationIndices()

					if e.config.ExportStackFrames {
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

							if len(e.config.ExportStackFrameTypes) > 0 &&
								!slices.Contains(e.config.ExportStackFrameTypes, unwindType) {
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
								e.consumeFunction(unwindType, functionName, fileName, line.Line(), line.Column())
							}
						}
					}

					fmt.Println("------------------- End Sample --------------------")
				}
				fmt.Println("------------------- End Profile -------------------")
			}
		}

		fmt.Printf("-------------- End Resource Profile ---------------\n\n")
	}
	return nil
}

func (e *customexporter) consumeFunction(frameType, functionName, fileName string, line, column int64) {
	if e.config.ExportFunctionFile {
		fmt.Printf("Instrumentation: %s, Function: %s, File: %s, Line: %d, Column: %d\n",
			frameType, functionName, fileName, line, column)
	} else {
		fmt.Printf("Instrumentation: %s, Function: %s\n",
			frameType, functionName)
	}
}

func (e *customexporter) Close(_ context.Context) error {
	e.logger.Info("Closing custom profiles exporter...")
	return nil
}
