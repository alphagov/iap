package cfg

import (
	"fmt"
)

// Example configuration file
// ---
// users:
//   value-of-oidc-claim:
//     roles:
//       - role1
//		   - role2
//   fname.lname@mydomain.com:
//     roles:
//       - superuser
//		   - readonlyuser

// UserConfig represents an unvalidated User configuration
type UserConfig struct {
	Roles []string `json:"roles"`
}

// ValidatedUserConfig represents a validated User configuration
type ValidatedUserConfig struct {
	Identifier string
	Roles      []string
}

// Validate does validation of UserConfig
func (c *UserConfig) Validate(identifier string) (ValidatedUserConfig, error) {
	cfg := ValidatedUserConfig{
		Identifier: identifier,
		Roles:      c.Roles,
	}

	if cfg.Identifier == "" {
		return cfg, fmt.Errorf("User Identifier must not be empty")
	}

	return cfg, nil
}
