package json

import "testing"
import "github.com/example/user-config-resolver-go/resolver"

func testCases() []struct {
	groups  []string
	in, out string
} {
	return []struct {
		groups  []string
		in, out string
	}{
		{[]string{"group-a", "group-b"}, "user-in-all-groups/input.json", "user-in-all-groups/output.json"},
		{[]string{"group-a"}, "user-not-in-all-groups/input.json", "user-not-in-all-groups/output.json"},
		{[]string{"group-d"}, "user-in-any-groups/input.json", "user-in-any-groups/output.json"},
		{[]string{"group-c"}, "user-in-no-groups/input.json", "user-in-no-groups/output.json"},
		{[]string{"group-a", "group-b", "group-c"}, "user-in-different-groups/input.json", "user-in-different-groups/output.json"},
		{[]string{"group-a", "group-b", "group-c"}, "custom-user-groups/input.json", "custom-user-groups/output.json"},
	}
}

func TestResolveStringToStruct(t *testing.T) {
	svc := New()
	for _, tc := range testCases() {
		in, err := readFile(tc.in)
		if err != nil {
			t.Fatal(err)
		}
		var expected TestDto
		if err := readFileInto(tc.out, &expected); err != nil {
			t.Fatal(err)
		}
		var result TestDto
		if err := svc.ResolveStringToStruct(in, tc.groups, &result); err != nil {
			t.Fatal(err)
		}
		if result != expected {
			t.Errorf("unexpected result for %s", tc.in)
		}
	}
}

func TestResolveStringToString(t *testing.T) {
	svc := New()
	for _, tc := range testCases() {
		in, err := readFile(tc.in)
		if err != nil {
			t.Fatal(err)
		}
		expected, err := readFile(tc.out)
		if err != nil {
			t.Fatal(err)
		}
		expected = compact(expected)
		out, err := svc.ResolveStringToString(in, tc.groups)
		if err != nil {
			t.Fatal(err)
		}
		if out != expected {
			t.Errorf("unexpected result for %s", tc.in)
		}
	}
}

func TestResolveStructToStruct(t *testing.T) {
	svc := New()
	for _, tc := range testCases() {
		var cfg resolver.Config
		if err := readFileInto(tc.in, &cfg); err != nil {
			t.Fatal(err)
		}
		var expected TestDto
		if err := readFileInto(tc.out, &expected); err != nil {
			t.Fatal(err)
		}
		var result TestDto
		if err := svc.ResolveStructToStruct(&cfg, tc.groups, &result); err != nil {
			t.Fatal(err)
		}
		if result != expected {
			t.Errorf("unexpected result for %s", tc.in)
		}
	}
}

func TestResolveStructToString(t *testing.T) {
	svc := New()
	for _, tc := range testCases() {
		var cfg resolver.Config
		_ = readFileInto(tc.in, &cfg)
		expected, _ := readFile(tc.out)
		expected = compact(expected)
		out, err := svc.ResolveStructToString(&cfg, tc.groups)
		if err != nil {
			t.Fatal(err)
		}
		if out != expected {
			t.Errorf("unexpected result for %s", tc.in)
		}
	}
}

func TestInvalidInput(t *testing.T) {
	svc := New()
	groups := []string{"group-a", "group-b"}
	in, _ := readFile("invalid-config/input.json")
	var v TestDto
	if err := svc.ResolveStringToStruct(in, groups, &v); err == nil {
		t.Error("expected error")
	}
	if _, err := svc.ResolveStringToString(in, groups); err == nil {
		t.Error("expected error")
	}
	if err := svc.ResolveStructToStruct(nil, groups, &v); err == nil {
		t.Error("expected error")
	}
	if _, err := svc.ResolveStructToString(nil, groups); err == nil {
		t.Error("expected error")
	}
}

func compact(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if r != ' ' && r != '\n' && r != '\t' && r != '\r' {
			out = append(out, r)
		}
	}
	return string(out)
}
