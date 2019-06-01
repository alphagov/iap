package cfg

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OIDC Config", func() {
	It("Parses a valid configuration without defaults", func() {
		cfg := OIDCConfig{
			RedirectURI: "https://iap.mydomain.com/oidc/callback",

			AuthURI:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURI: "https://www.googleapis.com/oauth2/v4/token",

			Scopes:          []string{"openid", "email"},
			IdentifierClaim: "email",

			ClientID:     "my-client-id",
			ClientSecret: "my-client-secret",
		}

		validatedCfg, err := cfg.Validate()

		Expect(err).NotTo(HaveOccurred())

		Expect(validatedCfg.RedirectURI.String()).To(
			Equal("https://iap.mydomain.com/oidc/callback"),
		)

		Expect(validatedCfg.AuthURI.String()).To(
			Equal("https://accounts.google.com/o/oauth2/v2/auth"),
		)
		Expect(validatedCfg.TokenURI.String()).To(
			Equal("https://www.googleapis.com/oauth2/v4/token"),
		)

		Expect(validatedCfg.Scopes).To(ConsistOf("openid", "email"))
		Expect(validatedCfg.IdentifierClaim).To(Equal("email"))

		Expect(validatedCfg.ClientID).To(Equal("my-client-id"))
		Expect(validatedCfg.ClientSecret).To(ContainSubstring("my-client-secret"))
	})

	It("Parses a valid configuration and provides defaults", func() {
		cfg := OIDCConfig{
			RedirectURI: "https://iap.mydomain.com/oidc/callback",

			AuthURI:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURI: "https://www.googleapis.com/oauth2/v4/token",

			ClientID:     "my-client-id",
			ClientSecret: "my-client-secret",
		}

		validatedCfg, err := cfg.Validate()

		Expect(err).NotTo(HaveOccurred())
		Expect(validatedCfg.Scopes).To(ConsistOf("openid", "email"))
	})

	It("Does not validate a configuration with an invalid redirect uri", func() {
		cfg := OIDCConfig{
			RedirectURI: "!",

			AuthURI:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURI: "https://www.googleapis.com/oauth2/v4/token",

			ClientID:     "my-client-id",
			ClientSecret: "my-client-secret",
		}

		_, err := cfg.Validate()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(
			ContainSubstring("OIDC RedirectURI must be a valid URI"),
		))
	})

	It("Does not validate a configuration with an invalid auth uri", func() {
		cfg := OIDCConfig{
			RedirectURI: "https://iap.mydomain.com/oidc/callback",

			AuthURI:  "!",
			TokenURI: "https://www.googleapis.com/oauth2/v4/token",

			ClientID:     "my-client-id",
			ClientSecret: "my-client-secret",
		}

		_, err := cfg.Validate()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(
			ContainSubstring("OIDC AuthURI must be a valid URI"),
		))
	})

	It("Does not validate a configuration with an invalid token uri", func() {
		cfg := OIDCConfig{
			RedirectURI: "https://iap.mydomain.com/oidc/callback",

			AuthURI:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURI: "!",

			ClientID:     "my-client-id",
			ClientSecret: "my-client-secret",
		}

		_, err := cfg.Validate()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(
			ContainSubstring("OIDC TokenURI must be a valid URI"),
		))
	})

	It("Does not validate a configuration without a client id", func() {
		cfg := OIDCConfig{
			RedirectURI: "https://iap.mydomain.com/oidc/callback",

			AuthURI:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURI: "https://www.googleapis.com/oauth2/v4/token",

			ClientSecret: "my-client-secret",
		}

		_, err := cfg.Validate()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(
			ContainSubstring("OIDC ClientID must be present"),
		))
	})

	It("Does not validate a configuration without a client id", func() {
		cfg := OIDCConfig{
			RedirectURI: "https://iap.mydomain.com/oidc/callback",

			AuthURI:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURI: "https://www.googleapis.com/oauth2/v4/token",

			ClientID: "my-client-id",
		}

		_, err := cfg.Validate()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(
			ContainSubstring("OIDC ClientSecret must be present"),
		))
	})
})
