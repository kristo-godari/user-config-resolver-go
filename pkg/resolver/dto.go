package resolver

type Config struct {
	OverrideRules     []OverrideRule         `json:"override-rules"`
	DefaultProperties map[string]interface{} `json:"default-properties"`
}

type OverrideRule struct {
	UserIsInAllGroups       []string               `json:"user-is-in-all-groups"`
	UserIsInAnyGroup        []string               `json:"user-is-in-any-group"`
	UserIsInNoneOfTheGroups []string               `json:"user-is-none-of-the-groups"`
	CustomExpression        string                 `json:"custom-expression"`
	Override                map[string]interface{} `json:"override"`
}
