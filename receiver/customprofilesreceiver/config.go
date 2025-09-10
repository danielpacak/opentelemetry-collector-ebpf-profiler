package customprofilesreceiver

import "time"

type Config struct {
	Foo            string        `mapstructure:"foo"`
	ReportInterval time.Duration `mapstructure:"report_interval"`
}
