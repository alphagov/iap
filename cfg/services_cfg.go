package cfg

import (
	"fmt"
	"net/url"

	"github.com/goware/urlx"
)

// Example configuration file
// ---
// services:
//   my-service:
//     upstream_uri: http://my-service.local
//     matchers:
//       - host: my-service.mydomain.com
//       - host: my-svc.mydomain.com
//     headers:
//       Authorization: Basic my-basic-auth-secret
//
//   my-other-service:
//     upstream_uri: http://my-service.local
//     matchers:
//       - host: my-other-service.mydomain.com
//

// MatcherConfig represents an unvalidated Matcher configuration
type MatcherConfig struct {
	Host string
}

// ValidatedMatcherConfig represents a validated Matcher configuration
type ValidatedMatcherConfig struct {
	Host string
}

// ServiceConfig represents an unvalidated Service configuration
type ServiceConfig struct {
	UpstreamURI string            `json:"upstream_uri"`
	Matchers    []MatcherConfig   `json:"matchers"`
	Headers     map[string]string `json:"headers"`
}

// Validate does validation of MatcherConfig
func (c *MatcherConfig) Validate() (ValidatedMatcherConfig, error) {
	cfg := ValidatedMatcherConfig{}

	if c.Host == "" {
		return cfg, fmt.Errorf("Matcher Host cannot be empty")
	}

	return ValidatedMatcherConfig{
		Host: c.Host,
	}, nil
}

// ValidatedServiceConfig represents a validated Service configuration
type ValidatedServiceConfig struct {
	Identifier  string
	UpstreamURI url.URL
	Matchers    []ValidatedMatcherConfig
	Headers     map[string]string
}

// Validate does validation of ServiceConfig
func (c *ServiceConfig) Validate(identifier string) (ValidatedServiceConfig, error) {
	cfg := ValidatedServiceConfig{}

	if identifier == "" {
		return cfg, fmt.Errorf("Service Identifier cannot be empty")
	}

	upstreamURI, err := urlx.ParseWithDefaultScheme(c.UpstreamURI, "https")
	if err != nil {
		return cfg, fmt.Errorf("Service Upstream URI must be a valid URI")
	}

	validatedMatchers := make([]ValidatedMatcherConfig, 0)

	for index, matcher := range c.Matchers {
		validatedMatcher, err := matcher.Validate()
		if err != nil {
			return cfg, fmt.Errorf(
				"Service Matcher %d was not valid: %s", index, err,
			)
		}
		validatedMatchers = append(validatedMatchers, validatedMatcher)
	}

	return ValidatedServiceConfig{
		Identifier:  identifier,
		UpstreamURI: *upstreamURI,
		Matchers:    validatedMatchers,
		Headers:     c.Headers,
	}, nil
}
