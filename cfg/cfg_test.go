package cfg

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/lithammer/dedent"
)

var _ = Describe("Config", func() {
	It("Rejects an empty configuration", func() {
		cfg := Config{
			OIDCConfig: OIDCConfig{},
			Roles:      []string{},
			Services:   map[string]ServiceConfig{},
			Users:      map[string]UserConfig{},
		}

		_, err := cfg.Validate()
		Expect(err).To(HaveOccurred())
	})

	It("Parses a valid configuration without defaults", func() {
		cfg := Config{
			OIDCConfig: OIDCConfig{
				RedirectURI: "https://iap.mydomain.com/oidc/callback",

				AuthURI:  "https://accounts.google.com/o/oauth2/v2/auth",
				TokenURI: "https://www.googleapis.com/oauth2/v4/token",

				Scopes:          []string{"openid", "email"},
				IdentifierClaim: "email",

				ClientID:     "my-client-id",
				ClientSecret: "my-client-secret",
			},

			Roles: []string{"superuser", "readonlyuser"},

			Services: map[string]ServiceConfig{
				"my-service": ServiceConfig{
					UpstreamURI: "http://my-service.local",
					Matchers: []MatcherConfig{
						MatcherConfig{Host: "my-service.mydomain.com"},
						MatcherConfig{Host: "my-svc.mydomain.com"},
					},
					Headers: map[string]string{
						"Authorization": "Basic my-basic-auth-secret",
					},
				},
				"my-other-service": ServiceConfig{
					UpstreamURI: "http://my-other-service.local",
					Matchers: []MatcherConfig{
						MatcherConfig{Host: "my-other-service.mydomain.com"},
					},
				},
			},

			Users: map[string]UserConfig{
				"a-user@mydomain.com": UserConfig{
					Roles: []string{"superuser"},
				},
				"another-user@mydomain.com": UserConfig{
					Roles: []string{"readonlyuser"},
				},
			},
		}

		validatedCfg, err := cfg.Validate()
		Expect(err).NotTo(HaveOccurred())

		Expect(validatedCfg.Roles).To(HaveLen(2))
		Expect(validatedCfg.Services).To(HaveLen(2))
		Expect(validatedCfg.Users).To(HaveLen(2))
	})
})

var _ = Describe("Config from String", func() {
	It("Rejects an empty configuration", func() {
		_, err := ParseAndValidateConfig("")
		Expect(err).To(HaveOccurred())
	})

	It("Parses valid YAML configuration without defaults", func() {
		config := dedent.Dedent(`
    oidc:
      redirect_uri: https://iap.mydomain.com/oidc/callback
      auth_uri: https://accounts.google.com/o/oauth2/v2/auth
      token_uri: https://www.googleapis.com/oauth2/v4/token
      scopes: [openid, email]
      identifier_claim: email
      client_id: foo-0000-1111.apps.googleusercontent.com
      client_secret: abcd-0000-1111
    roles:
      - superuser
      - readonlyuser
    services:
      my-service:
        upstream_uri: http://my-service.local
        matchers:
          - host: my-service.mydomain.com
          - host: my-svc.mydomain.com
        headers:
          Authorization: Basic my-basic-auth-secret
      my-other-service:
        upstream_uri: http://my-service.local
        matchers:
          - host: my-other-service.mydomain.com
    users:
      value-of-oidc-claim:
        roles:
          - role1
          - role2
      fname.lname@mydomain.com:
        roles:
          - superuser
          - readonlyuser
		`)

		validatedCfg, err := ParseAndValidateConfig(config)

		Expect(err).NotTo(HaveOccurred())
		Expect(validatedCfg.Roles).To(HaveLen(2))
		Expect(validatedCfg.Services).To(HaveLen(2))
		Expect(validatedCfg.Users).To(HaveLen(2))
	})
})
