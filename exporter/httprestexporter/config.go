package httprestexporter

type Config struct {
	Address           string   `mapstructure:"address"`
	ExportSampleTypes []string `mapstructure:"export_sample_types"`
}
