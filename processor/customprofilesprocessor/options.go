package customprofilesprocessor

// option represents a configuration option that can be passes.
// to the customprocessor.
type option func(*customprocessor) error

func withFoo(foo string) option {
	return func(c *customprocessor) error {
		c.foo = foo
		return nil
	}
}
