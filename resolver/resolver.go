package resolver

// ConfigResolver resolves configuration for a set of user groups.
//
// Implementations parse the provided configuration rules, apply the group based
// overrides and return the resulting configuration either as a JSON string or by
// unmarshalling it into the supplied target value.
type ConfigResolver interface {
	// Resolve parses the configuration and returns the resolved JSON.
	Resolve(config string, userGroups []string) (string, error)

	// ResolveInto parses the configuration and populates target with the
	// resolved values. Target must be a pointer to the destination struct or
	// another value that can be unmarshalled by the implementation.
	ResolveInto(config string, userGroups []string, target any) error
}

type ConfigResolverError struct{ Err error }

func (e ConfigResolverError) Error() string { return e.Err.Error() }
func (e ConfigResolverError) Unwrap() error { return e.Err }
