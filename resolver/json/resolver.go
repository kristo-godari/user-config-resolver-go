package json

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/example/user-config-resolver-go/resolver"
)

type JsonConfigResolverService struct {
	configToResolve string
}

func New() *JsonConfigResolverService { return &JsonConfigResolverService{} }

func (s *JsonConfigResolverService) SetConfigToResolve(config string) { s.configToResolve = config }

func (s *JsonConfigResolverService) ResolveConfig(groups []string) (string, error) {
	if s.configToResolve == "" {
		return "", resolver.ConfigResolverError{Err: fmt.Errorf("config to resolve is empty")}
	}
	return s.ResolveConfigFrom(s.configToResolve, groups)
}

func (s *JsonConfigResolverService) ResolveConfigInto(groups []string, target any) error {
	if s.configToResolve == "" {
		return resolver.ConfigResolverError{Err: fmt.Errorf("config to resolve is empty")}
	}
	return s.ResolveConfigFromInto(s.configToResolve, groups, target)
}

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
