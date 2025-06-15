package resolver

import "testing"

func TestApplyRules(t *testing.T) {
	c := Config{
		DefaultProperties: map[string]interface{}{
			"foo": map[string]interface{}{"bar": 1},
		},
		OverrideRules: []OverrideRule{
			{
				UserIsInAllGroups: []string{"a", "b"},
				Override:          map[string]interface{}{"foo.bar": 2},
			},
		},
	}
	groups := []string{"a", "b"}
	ApplyRules(groups, &c)
	foo := c.DefaultProperties["foo"].(map[string]interface{})
	if foo["bar"] != 2 {
		t.Fatalf("expected override to apply")
	}
}

func TestRuleAppliesCustomExpression(t *testing.T) {
	r := OverrideRule{CustomExpression: "#user.contains('a') and #user.contains('b')"}
	if !ruleApplies([]string{"a", "b"}, r) {
		t.Fatal("expected expression to evaluate to true")
	}
	if ruleApplies([]string{"a"}, r) {
		t.Fatal("expected expression to evaluate to false")
	}
}

func TestConvertExpression(t *testing.T) {
	in := "#user.contains('a') or #user.contains(\"b\") and #user.contains('c')"
	expected := "'a' in user || 'b' in user && 'c' in user"
	if out := convertExpression(in); out != expected {
		t.Fatalf("unexpected conversion: %s", out)
	}
}

func TestOverrideProperty(t *testing.T) {
	m := map[string]interface{}{}
	overrideProperty(m, "a.b.c", 1)
	a := m["a"].(map[string]interface{})
	b := a["b"].(map[string]interface{})
	if b["c"] != 1 {
		t.Fatalf("expected nested property to be set")
	}
}
