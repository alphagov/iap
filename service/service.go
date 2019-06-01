package service

import (
	"net/url"
)

// Matcher represents a validated Matcher
type Matcher struct {
	Host string
}

// Service represents a validated Service
type Service struct {
	Headers     map[string]string
	Identifier  string
	Matchers    []Matcher
	Roles       []string
	UpstreamURI url.URL
}
