package service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service Roles", func() {
	It("A service with no roles is accessible to a user with no roles", func() {
		s := Service{
			Roles: []string{},
		}
		roles := []string{}

		Expect(s.IsAccessible(roles)).To(Equal(true))
	})

	It("A service with no roles is accessible to a user with roles", func() {
		s := Service{
			Roles: []string{},
		}
		roles := []string{"superuser", "readonly"}

		Expect(s.IsAccessible(roles)).To(Equal(true))
	})

	It("A service with roles is not accessible to a user with no roles", func() {
		s := Service{
			Roles: []string{"a-role"},
		}
		roles := []string{}

		Expect(s.IsAccessible(roles)).To(Equal(false))
	})

	It("A service with roles is not accessible to a user without role", func() {
		s := Service{
			Roles: []string{"a-role"},
		}
		roles := []string{"not-a-role"}

		Expect(s.IsAccessible(roles)).To(Equal(false))
	})

	It("A service with roles is accessible to a user with role", func() {
		s := Service{
			Roles: []string{"a-role"},
		}
		roles := []string{"a-role"}

		Expect(s.IsAccessible(roles)).To(Equal(true))
	})

	It("A service with roles is accessible to a user with many roles", func() {
		s := Service{
			Roles: []string{"a-role"},
		}
		roles := []string{"a-role", "another-role"}

		Expect(s.IsAccessible(roles)).To(Equal(true))
	})
})
