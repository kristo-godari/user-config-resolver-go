package json

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func readFile(name string) (string, error) {
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func readFileInto(name string, target any) error {
	s, err := readFile(name)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(s), target)
}
