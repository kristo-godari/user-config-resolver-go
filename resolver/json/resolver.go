package json

import (
	"encoding/json"
	"strings"

	"github.com/example/user-config-resolver-go/resolver"
)

type JsonConfigResolverService struct{}

func New() *JsonConfigResolverService { return &JsonConfigResolverService{} }

func (s *JsonConfigResolverService) ResolveConfigFrom(cfg string, groups []string) (string, error) {
	var out string
	if err := s.ResolveConfigFromInto(cfg, groups, &out); err != nil {
		return "", err
	}
	return out, nil
}

func (s *JsonConfigResolverService) ResolveConfigFromInto(cfg string, groups []string, target any) error {
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
