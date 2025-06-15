# User Config Resolver Library (Go)

User Config Resolver is a lightweight feature flag and configuration override library for Go. It helps you customise application behaviour by resolving settings for specific user groups or segments. Use it for targeted feature toggles, A/B testing, dynamic configuration or other personalisation scenarios.

This Go library mirrors the features of the [Java version](https://github.com/kristo-godari/user-config-resolver-java). It resolves user configuration based on user groups and custom expressions. The resolver supports pluggable input formats and ships with a JSON implementation out of the box.

The resolver can drive feature flags and configuration overrides without touching your code. Load a JSON file with override rules, provide the groups for a given user and receive the final settings as structs or raw JSON. Embedding the library into Go microservices or command line tools is straightforward.

## Use Case
You may want to adapt application behaviour for different user groups. Instead of encoding logic in code, you can keep it in configuration so that updating it does not require a redeploy.

## Features
- **User Group-Based Overrides:** override configuration properties when a user is in specific groups.
- **Custom Expression Support:** use simple expressions to create complex conditions.
- **Feature Flags & A/B Testing:** toggle features or run experiments without redeploying your application.
- **Resolve to Structs or JSON:** resolve configuration into Go structs or as JSON strings.

## Installation
```bash
go get github.com/example/user-config-resolver-go/resolver/json
```

### Extending With Custom Formats
Implement the `resolver.ConfigResolver` interface in your own package and provide the parsing logic for the desired format (for example XML). Place the implementation under `resolver/<format>` so it can be imported as `github.com/example/user-config-resolver-go/resolver/<format>`.

## Define Configuration Override Rules
Create a JSON configuration that describes your default properties and override rules.

```json
{
  "override-rules": [
    {
      "user-is-in-all-groups": ["paid-user", "premium-user"],
      "override": {"show-adds": false}
    },
    {
      "user-is-in-any-group": ["new-joiner"],
      "override": {
        "show-new-joiner-banner": true,
        "show-full-layout": false
      }
    },
    {
      "user-is-none-of-the-groups": ["button-blue"],
      "override": {"button-color": "gray"}
    },
    {
      "custom-expression": "#user.contains('discount') or #user.contains('black-friday')",
      "override": {
        "shop.no-of-products": 20,
        "shop.price-multiplier": 0
      }
    }
  ],
  "default-properties": {
    "show-new-joiner-banner": false,
    "show-adds": true,
    "show-full-layout": true,
    "button-color": "blue",
    "shop": {
      "no-of-products": 10,
      "price-multiplier": 2
    }
  }
}
```

Override rules are applied from top to bottom.

## Resolving Configuration
```go
import resjson "github.com/example/user-config-resolver-go/resolver/json"

svc := resjson.New()
var result MyConfigStruct
err := svc.ResolveConfigFromInto(configString, groups, &result)
```
Alternatively you can store the configuration using `SetConfigToResolve` and then call `ResolveConfig` or `ResolveConfigInto`.

## Example
Run the example to see the resolver in action:

```bash
go run ./example
```

`example/main.go` reads `example/config.json`, resolves it for a user in the `paid-user` and `discount` groups and prints the resulting JSON configuration.

## License
This library is released under the MIT License.

### Keywords
feature flags, configuration override, dynamic configuration, user segments, Go library, A/B testing


