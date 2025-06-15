package resolver

// ConfigResolver defines the interface for resolving user-specific configuration.
type ConfigResolver interface {
	ResolveStringToString(input string, groups []string) (string, error)
	ResolveStringToStruct(input string, groups []string, out any) error
	ResolveStructToString(cfg *Config, groups []string) (string, error)
	ResolveStructToStruct(cfg *Config, groups []string, out any) error
}

type ConfigResolverError struct{ Err error }

func (e ConfigResolverError) Error() string { return e.Err.Error() }
func (e ConfigResolverError) Unwrap() error { return e.Err }
