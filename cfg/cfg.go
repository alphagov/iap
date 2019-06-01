package cfg

import (
	"fmt"

	"github.com/ghodss/yaml"
)

// Example configuration file ---
// oidc: <oidc config>
//
// roles: [role1, role2]
//
// services:
//   my-service: <service config>
//
// users:
//   user-identifier-1: <user config>

// Config represents an unvalidated configuration
type Config struct {
	OIDCConfig OIDCConfig               `json:"oidc"`
	Roles      []string                 `json:"roles"`
	Services   map[string]ServiceConfig `json:"services"`
	Users      map[string]UserConfig    `json:"users"`
}

// ValidatedConfig represents a validated configuration
type ValidatedConfig struct {
	OIDCConfig ValidatedOIDCConfig
	Roles      []string
	Services   map[string]ValidatedServiceConfig
	Users      map[string]ValidatedUserConfig
}

// Validate does validation of MatcherConfig
func (c *Config) Validate() (ValidatedConfig, error) {
	cfg := ValidatedConfig{}

	validatedOIDCConfig, err := c.OIDCConfig.Validate()
	if err != nil {
		return cfg, err
	}

	validatedServices := make(map[string]ValidatedServiceConfig)
	for serviceIdentifier, serviceConfig := range c.Services {
		validatedServiceConfig, err := serviceConfig.Validate(serviceIdentifier)

		if err != nil {
			return cfg, fmt.Errorf(
				"Service %s is not valid %s", serviceIdentifier, err,
			)
		}

		validatedServices[serviceIdentifier] = validatedServiceConfig
	}

	validatedUsers := make(map[string]ValidatedUserConfig)
	for userIdentifier, userConfig := range c.Users {
		validatedUserConfig, err := userConfig.Validate(userIdentifier)

		if err != nil {
			return cfg, fmt.Errorf(
				"User %s is not valid %s", userIdentifier, err,
			)
		}

		validatedUsers[userIdentifier] = validatedUserConfig
	}

	return ValidatedConfig{
		OIDCConfig: validatedOIDCConfig,
		Roles:      c.Roles,
		Services:   validatedServices,
		Users:      validatedUsers,
	}, nil
}

// ParseAndValidateConfig parses and validates configuration from a string
func ParseAndValidateConfig(config string) (ValidatedConfig, error) {
	cfg := Config{}
	err := yaml.Unmarshal([]byte(config), &cfg)

	if err != nil {
		return ValidatedConfig{}, fmt.Errorf("Could not unmarshal config: %s", err)
	}

	validatedCfg, err := cfg.Validate()

	if err != nil {
		return ValidatedConfig{}, fmt.Errorf("Could not validate config: %s", err)
	}

	return validatedCfg, nil
}
