package main

import (
	"encoding/json"
	"fmt"
	"os"

	jsonresolver "github.com/example/user-config-resolver-go/pkg/resolver/json"
)

type ShopConfig struct {
	NoOfProducts    int `json:"no-of-products"`
	PriceMultiplier int `json:"price-multiplier"`
}

type ExampleConfig struct {
	ShowNewJoinerBanner bool       `json:"show-new-joiner-banner"`
	ShowAdds            bool       `json:"show-adds"`
	ShowFullLayout      bool       `json:"show-full-layout"`
	ButtonColor         string     `json:"button-color"`
	Shop                ShopConfig `json:"shop"`
}

func main() {
	// Read the configuration file.
	raw, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	// Groups that the current user belongs to.
	groups := []string{"paid-user", "discount"}

	svc := jsonresolver.New()
	var result ExampleConfig
	if err := svc.ResolveConfigFromInto(string(raw), groups, &result); err != nil {
		panic(err)
	}

	// Print the resolved configuration.
	out, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(out))
}
