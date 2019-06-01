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

// IsAccessible returns if service should be accessed by any one of the roles
func (s *Service) IsAccessible(roles []string) bool {
	if len(s.Roles) == 0 {
		return true
	}

	for _, serviceRole := range s.Roles {
		for _, role := range roles {
			if serviceRole == role {
				return true
			}
		}
	}

	return false
}
