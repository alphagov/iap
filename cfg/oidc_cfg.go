package cfg

import (
	"fmt"
	"net/url"

	"github.com/goware/urlx"
)

// Example configuration file
// ---
// oidc:
//   redirect_uri: https://iap.mydomain.com/oidc/callback
//
//   auth_uri: https://accounts.google.com/o/oauth2/v2/auth
//   token_uri: https://www.googleapis.com/oauth2/v4/token
//
//   scopes: [openid, email]
//   identifier_claim: email
//
//   client_id: foo-0000-1111.apps.googleusercontent.com
//   client_secret: abcd00001111

// OIDCConfig represents an unvalidated OIDC configuration
type OIDCConfig struct {
	RedirectURI string `json:"redirect_uri"`

	AuthURI  string `json:"auth_uri"`
	TokenURI string `json:"token_uri"`

	Scopes          []string `json:"scopes"`
	IdentifierClaim string   `json:"identifier_claim"`

	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// ValidatedOIDCConfig represents a validated OIDC configuration
type ValidatedOIDCConfig struct {
	RedirectURI url.URL

	AuthURI  url.URL
	TokenURI url.URL

	Scopes          []string
	IdentifierClaim string

	ClientID     string
	ClientSecret string
}

// Validate does validation of OIDCConfig
func (c *OIDCConfig) Validate() (ValidatedOIDCConfig, error) {
	cfg := ValidatedOIDCConfig{}

	redirectURI, err := urlx.Parse(c.RedirectURI)
	if err != nil {
		return cfg, fmt.Errorf("OIDC RedirectURI must be a valid URI: %s", err)
	}

	authURI, err := urlx.Parse(c.AuthURI)
	if err != nil {
		return cfg, fmt.Errorf("OIDC AuthURI must be a valid URI: %s", err)
	}

	tokenURI, err := urlx.Parse(c.TokenURI)
	if err != nil {
		return cfg, fmt.Errorf("OIDC TokenURI must be a valid URI: %s", err)
	}

	scopes := c.Scopes
	if len(scopes) == 0 {
		scopes = []string{"openid", "email"}
	}

	identifierClaim := c.IdentifierClaim
	if identifierClaim == "" {
		identifierClaim = "email"
	}

	if c.ClientID == "" {
		return cfg, fmt.Errorf("OIDC ClientID must be present")
	}

	if c.ClientSecret == "" {
		return cfg, fmt.Errorf("OIDC ClientSecret must be present")
	}

	return ValidatedOIDCConfig{
		RedirectURI: *redirectURI,

		AuthURI:  *authURI,
		TokenURI: *tokenURI,

		Scopes:          scopes,
		IdentifierClaim: identifierClaim,

		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
	}, nil
}
