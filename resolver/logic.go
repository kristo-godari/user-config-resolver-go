package resolver

import (
	"regexp"
	"strings"

	"github.com/Knetic/govaluate"
)

// ApplyRules applies override rules to the provided configuration based on user groups.
func ApplyRules(groups []string, c *Config) {
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
