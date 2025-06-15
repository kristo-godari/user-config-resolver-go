package json

type TestDto struct {
	Property1 int `json:"property1"`
	Property2 struct {
		Property21 bool `json:"property2-1"`
	} `json:"property2"`
	Property3 struct {
		Property31 struct {
			Property311 bool `json:"property3-1-1"`
		} `json:"property3-1"`
	} `json:"property3"`
}
