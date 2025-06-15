package json

import (
	"encoding/json"
	"strings"

	"github.com/example/user-config-resolver-go/resolver"
)

type JsonConfigResolverService struct{}

func New() *JsonConfigResolverService { return &JsonConfigResolverService{} }

// Resolve returns the resolved configuration as a JSON string.
func (s *JsonConfigResolverService) Resolve(cfg string, groups []string) (string, error) {
	var out string
	if err := s.ResolveInto(cfg, groups, &out); err != nil {
		return "", err
	}
	return out, nil
}

// ResolveInto unmarshals the resolved configuration into target.
func (s *JsonConfigResolverService) ResolveInto(cfg string, groups []string, target any) error {
	var c resolver.Config
	dec := json.NewDecoder(strings.NewReader(cfg))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&c); err != nil {
		return resolver.ConfigResolverError{Err: err}
	}
	resolver.ApplyRules(groups, &c)
	data, err := json.Marshal(c.DefaultProperties)
	if err != nil {
		return resolver.ConfigResolverError{Err: err}
	}
	if strPtr, ok := target.(*string); ok {
		*strPtr = string(data)
		return nil
	}
	return json.Unmarshal(data, target)
}
