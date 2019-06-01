package cfg

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
