package cfg

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Matcher Config", func() {
	It("Parses a valid configuration", func() {
		cfg := MatcherConfig{Host: "my-service.mydomain.com"}
		validatedCfg, err := cfg.Validate()

		Expect(err).NotTo(HaveOccurred())
		Expect(validatedCfg.Host).To(Equal("my-service.mydomain.com"))
	})

	It("Does not validated a configuration with an empty host", func() {
		cfg := MatcherConfig{Host: ""}
		_, err := cfg.Validate()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring(
			"Matcher Host cannot be empty",
		)))
	})
})

var _ = Describe("Service Config", func() {
	It("Parses a valid configuration", func() {
		cfg := ServiceConfig{
			UpstreamURI: "my-service.local",
			Matchers: []MatcherConfig{
				MatcherConfig{Host: "my-service.mydomain.com"},
				MatcherConfig{Host: "my-svc.mydomain.com"},
			},
			Headers: map[string]string{
				"Authorization": "Basic my-basic-auth",
			},
		}

		validatedCfg, err := cfg.Validate("my-service")

		Expect(err).NotTo(HaveOccurred())
		Expect(validatedCfg.Identifier).To(Equal("my-service"))
		Expect(validatedCfg.Headers).To(Equal(cfg.Headers))

		Expect(validatedCfg.UpstreamURI.String()).To(Equal(
			"https://my-service.local",
		))

		Expect(validatedCfg.Matchers).To(HaveLen(2))
		Expect(validatedCfg.Matchers[0].Host).To(Equal("my-service.mydomain.com"))
		Expect(validatedCfg.Matchers[1].Host).To(Equal("my-svc.mydomain.com"))
	})

	It("Parses a valid configuration even if matchers are omitted", func() {
		cfg := ServiceConfig{
			UpstreamURI: "my-service.local",
			Headers: map[string]string{
				"Authorization": "Basic my-basic-auth",
			},
		}

		validatedCfg, err := cfg.Validate("my-service")

		Expect(err).NotTo(HaveOccurred())
		Expect(validatedCfg.Identifier).To(Equal("my-service"))
		Expect(validatedCfg.Headers).To(Equal(cfg.Headers))

		Expect(validatedCfg.UpstreamURI.String()).To(Equal(
			"https://my-service.local",
		))

		Expect(validatedCfg.Matchers).To(HaveLen(0))
	})

	It("Parses a valid configuration even if headers are omitted", func() {
		cfg := ServiceConfig{
			UpstreamURI: "my-service.local",
			Matchers: []MatcherConfig{
				MatcherConfig{Host: "my-service.mydomain.com"},
				MatcherConfig{Host: "my-svc.mydomain.com"},
			},
		}

		validatedCfg, err := cfg.Validate("my-service")

		Expect(err).NotTo(HaveOccurred())
		Expect(validatedCfg.Identifier).To(Equal("my-service"))

		Expect(validatedCfg.UpstreamURI.String()).To(Equal(
			"https://my-service.local",
		))

		Expect(validatedCfg.Matchers).To(HaveLen(2))
		Expect(validatedCfg.Matchers[0].Host).To(Equal("my-service.mydomain.com"))
		Expect(validatedCfg.Matchers[1].Host).To(Equal("my-svc.mydomain.com"))
	})

	It("Does not validate a configuration without a name", func() {
		cfg := ServiceConfig{
			UpstreamURI: "my-service.local",
			Matchers: []MatcherConfig{
				MatcherConfig{Host: "my-service.mydomain.com"},
				MatcherConfig{Host: "my-svc.mydomain.com"},
			},
			Headers: map[string]string{
				"Authorization": "Basic my-basic-auth",
			},
		}

		_, err := cfg.Validate("")

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring(
			"Service Identifier cannot be empty",
		)))
	})

	It("Does not validate a configuration with an invalid upstream uri", func() {
		cfg := ServiceConfig{
			UpstreamURI: "!",
			Matchers: []MatcherConfig{
				MatcherConfig{Host: "my-service.mydomain.com"},
				MatcherConfig{Host: "my-svc.mydomain.com"},
			},
			Headers: map[string]string{
				"Authorization": "Basic my-basic-auth",
			},
		}

		_, err := cfg.Validate("my-service")

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring(
			"Service Upstream URI must be a valid URI",
		)))
	})

	It("Does not validate a configuration with an invalid matcher", func() {
		cfg := ServiceConfig{
			UpstreamURI: "my-service.local",
			Matchers: []MatcherConfig{
				MatcherConfig{Host: "my-service.mydomain.com"},
				MatcherConfig{},
			},
			Headers: map[string]string{
				"Authorization": "Basic my-basic-auth",
			},
		}

		_, err := cfg.Validate("my-service")

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring(
			"Service Matcher 1 was not valid",
		)))
	})
})
