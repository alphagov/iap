package iap_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Iap Suite")
}