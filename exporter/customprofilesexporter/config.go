package customprofilesexporter

type Config struct {
	ExportResourceAttributes bool     `mapstructure:"export_resource_attributes"`
	ExportProfileAttributes  bool     `mapstructure:"export_profile_attributes"`
	ExportSampleAttributes   bool     `mapstructure:"export_sample_attributes"`
	ExportStackFrames        bool     `mapstructure:"export_stack_frames"`
	ExportStackFrameTypes    []string `mapstructure:"export_stack_frame_types"`
	ExportFunctionFile       bool     `mapstructure:"export_function_file"`
	// TODO Change this param to ExportProfileWithoutContainerID
	IgnoreProfilesWithoutContainerID bool     `mapstructure:"ignore_profiles_without_container_id"`
	ExportSampleTypes                []string `mapstructure:"export_sample_types"`
}
