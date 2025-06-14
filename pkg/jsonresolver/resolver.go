package jsonresolver

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/example/user-config-resolver-go/pkg/resolver"
)

type JsonConfigResolverService struct {
	configToResolve string
}

func New() *JsonConfigResolverService { return &JsonConfigResolverService{} }

func (s *JsonConfigResolverService) SetConfigToResolve(config string) { s.configToResolve = config }

func (s *JsonConfigResolverService) ResolveConfig(groups []string) (string, error) {
	if s.configToResolve == "" {
		return "", resolver.ConfigResolverError{fmt.Errorf("config to resolve is empty")}
	}
	return s.ResolveConfigFrom(s.configToResolve, groups)
}

func (s *JsonConfigResolverService) ResolveConfigInto(groups []string, target any) error {
	if s.configToResolve == "" {
		return resolver.ConfigResolverError{fmt.Errorf("config to resolve is empty")}
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
	var c Config
	dec := json.NewDecoder(strings.NewReader(cfg))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&c); err != nil {
		return resolver.ConfigResolverError{err}
	}
	applyRules(groups, &c)
	data, err := json.Marshal(c.DefaultProperties)
	if err != nil {
		return resolver.ConfigResolverError{err}
	}
	if strPtr, ok := target.(*string); ok {
		*strPtr = string(data)
		return nil
	}
	return json.Unmarshal(data, target)
}

func applyRules(groups []string, c *Config) {
	for _, r := range c.OverrideRules {
		if ruleApplies(groups, r) {
			overrideProperties(c.DefaultProperties, r.Override)
		}
	}
}

func ruleApplies(groups []string, r OverrideRule) bool {
	if len(r.UserIsInAllGroups) > 0 && containsAll(groups, r.UserIsInAllGroups) {
		return true
	}
	if len(r.UserIsInAnyGroup) > 0 && containsAny(groups, r.UserIsInAnyGroup) {
		return true
	}
	if len(r.UserIsInNoneOfTheGroups) > 0 && !containsAny(groups, r.UserIsInNoneOfTheGroups) {
		return true
	}
	if r.CustomExpression != "" {
		exprStr := convertExpression(r.CustomExpression)
		expr, err := govaluate.NewEvaluableExpression(exprStr)
		if err == nil {
			params := map[string]any{"user": toAnySlice(groups)}
			if res, err := expr.Evaluate(params); err == nil {
				if b, ok := res.(bool); ok {
					return b
				}
			}
		}
	}
	return false
}

func containsAll(userGroups, required []string) bool {
	set := make(map[string]struct{}, len(userGroups))
	for _, g := range userGroups {
		set[g] = struct{}{}
	}
	for _, g := range required {
		if _, ok := set[g]; !ok {
			return false
		}
	}
	return true
}

func containsAny(userGroups, any []string) bool {
	set := make(map[string]struct{}, len(userGroups))
	for _, g := range userGroups {
		set[g] = struct{}{}
	}
	for _, g := range any {
		if _, ok := set[g]; ok {
			return true
		}
	}
	return false
}

func overrideProperties(base, overrides map[string]interface{}) {
	for k, v := range overrides {
		overrideProperty(base, k, v)
	}
}

func overrideProperty(node map[string]interface{}, path string, value interface{}) {
	parts := strings.Split(path, ".")
	m := node
	for i, p := range parts {
		if i == len(parts)-1 {
			m[p] = value
			return
		}
		next, ok := m[p].(map[string]interface{})
		if !ok {
			next = map[string]interface{}{}
			m[p] = next
		}
		m = next
	}
}

func convertExpression(expr string) string {
	re := regexp.MustCompile(`#user\.contains\(['"]([^'"]+)['"]\)`)
	expr = re.ReplaceAllString(expr, "'$1' in user")
	expr = strings.ReplaceAll(expr, " or ", " || ")
	expr = strings.ReplaceAll(expr, " and ", " && ")
	return expr
}

func toAnySlice(in []string) []interface{} {
	out := make([]interface{}, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}
