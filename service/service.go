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
	Identifier  string
	UpstreamURI url.URL
	Matchers    []Matcher
	Headers     map[string]string
}
