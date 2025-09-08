package customprofilesprocessor

// option represents a configuration option that can be passes.
// to the k8s-tagger
type option func(*customprocessor) error
