package cfg

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Config", func() {
	It("Parses a valid configuration without defaults", func() {
		cfg := UserConfig{Roles: []string{"role1", "role2"}}
		user, err := cfg.Validate("user@mydomain.com")

		Expect(err).NotTo(HaveOccurred())
		Expect(user.Roles).To(Equal(cfg.Roles))
		Expect(user.Identifier).To(Equal("user@mydomain.com"))
	})

	It("Does not validate a configuration without an identifier", func() {
		cfg := UserConfig{Roles: []string{"role1", "role2"}}
		_, err := cfg.Validate("")

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring(
			"User Identifier must not be empty",
		)))
	})
})
