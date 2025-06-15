package json

import "testing"

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

func TestResolveInto(t *testing.T) {
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
		if err := svc.ResolveInto(in, tc.groups, &result); err != nil {
			t.Fatal(err)
		}
		if result != expected {
			t.Errorf("unexpected result for %s", tc.in)
		}
	}
}

func TestResolve(t *testing.T) {
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
		out, err := svc.Resolve(in, tc.groups)
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
	if err := svc.ResolveInto(in, groups, &v); err == nil {
		t.Error("expected error")
	}
	if _, err := svc.Resolve(in, groups); err == nil {
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
