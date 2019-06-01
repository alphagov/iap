package cfg

import (
	"fmt"

	"github.com/alphagov/iap/user"
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

// Validate does validation of UserConfig
func (c *UserConfig) Validate(identifier string) (user.User, error) {
	cfg := user.User{
		Identifier: identifier,
		Roles:      c.Roles,
	}

	if cfg.Identifier == "" {
		return cfg, fmt.Errorf("User Identifier must not be empty")
	}

	return cfg, nil
}
