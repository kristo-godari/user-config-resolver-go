package api

// ConfigResolver defines the interface for resolving user-specific configuration.
type ConfigResolver interface {
	SetConfigToResolve(config string)
	ResolveConfig(userGroups []string) (string, error)
	ResolveConfigInto(userGroups []string, target any) error
	ResolveConfigFrom(config string, userGroups []string) (string, error)
	ResolveConfigFromInto(config string, userGroups []string, target any) error
}

type ConfigResolverError struct{ Err error }

func (e ConfigResolverError) Error() string { return e.Err.Error() }
func (e ConfigResolverError) Unwrap() error { return e.Err }
