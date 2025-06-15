package json

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/example/user-config-resolver-go/resolver"
)

type JsonConfigResolverService struct{}

func New() *JsonConfigResolverService { return &JsonConfigResolverService{} }

func (s *JsonConfigResolverService) ResolveStringToString(input string, groups []string) (string, error) {
	var out string
	if err := s.ResolveStringToStruct(input, groups, &out); err != nil {
		return "", err
	}
	return out, nil
}

func (s *JsonConfigResolverService) ResolveStringToStruct(input string, groups []string, out any) error {
	var c resolver.Config
	dec := json.NewDecoder(strings.NewReader(input))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&c); err != nil {
		return resolver.ConfigResolverError{Err: err}
	}
	resolved := resolver.ApplyRules(groups, &c)
	data, err := json.Marshal(resolved.DefaultProperties)
	if err != nil {
		return resolver.ConfigResolverError{Err: err}
	}
	if strPtr, ok := out.(*string); ok {
		*strPtr = string(data)
		return nil
	}
	return json.Unmarshal(data, out)
}

func (s *JsonConfigResolverService) ResolveStructToString(cfg *resolver.Config, groups []string) (string, error) {
	var out string
	if err := s.ResolveStructToStruct(cfg, groups, &out); err != nil {
		return "", err
	}
	return out, nil
}

func (s *JsonConfigResolverService) ResolveStructToStruct(cfg *resolver.Config, groups []string, out any) error {
	if cfg == nil {
		return resolver.ConfigResolverError{Err: fmt.Errorf("config to resolve is empty")}
	}
	resolved := resolver.ApplyRules(groups, cfg)
	data, err := json.Marshal(resolved.DefaultProperties)
	if err != nil {
		return resolver.ConfigResolverError{Err: err}
	}
	if strPtr, ok := out.(*string); ok {
		*strPtr = string(data)
		return nil
	}
	return json.Unmarshal(data, out)
}
