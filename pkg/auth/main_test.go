package auth_test

import (
	"fmt"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/alphagov/iap/pkg/auth"
	"github.com/go-redis/redis"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("Auth package", func() {
	var (
		a  *auth.Client
		r  *redis.Client
		mr *miniredis.Miniredis
	)

	BeforeSuite(func() {
		var err error
		mr, err = miniredis.Run()
		Expect(err).NotTo(HaveOccurred())

		Expect(mr.Set(fmt.Sprintf(auth.UserTempPassword, "test"), "my_passw0rd!")).NotTo(HaveOccurred())
		Expect(mr.Set(fmt.Sprintf(auth.UserTempPassword, "jeff"), "jefferson")).NotTo(HaveOccurred())
		mr.SetTTL(fmt.Sprintf(auth.UserTempPassword, "jeff"), time.Second)

		r = redis.NewClient(&redis.Options{
			Addr: mr.Addr(),

			ReadTimeout:  time.Second * 2,
			WriteTimeout: time.Second * 1,
			DialTimeout:  time.Second * 1,
		})

		logrus.SetOutput(GinkgoWriter)
		a = auth.New(r, logrus.New())
	})

	AfterSuite(func() {
		r.Close()
		mr.Close()
	})

	It("should fail to validate user", func() {
		Expect(a.Valid("test", "my_password!")).To(BeFalse())
		Expect(a.Valid("jeff", "my_password!")).To(BeFalse())
		Expect(a.Valid("chuck", "my_password!")).To(BeFalse())
	})

	It("should fail to validate expired user", func() {
		Expect(a.Valid("jeff", "jefferson")).To(BeTrue())
		mr.FastForward(5 * time.Second)
		Expect(a.Valid("jeff", "jefferson")).To(BeFalse())
	})

	It("should successfully validate user", func() {
		Expect(a.Valid("test", "my_passw0rd!")).To(BeTrue())
	})

	It("should successfully generate random username and password", func() {
		u1a, p1a, e1a := a.Generate("qwertyuiop1234567890")
		u1b, p1b, e1b := a.Generate("qwertyuiop1234567890")
		u2, p2, e2 := a.Generate("1234567890qwertyuiop")

		Expect(e1a).NotTo(HaveOccurred())
		Expect(u1a).NotTo(Equal(u2))
		Expect(u1a).To(Equal(u1b))

		Expect(e1b).NotTo(HaveOccurred())
		Expect(p1a).To(Equal(p1b))

		Expect(e2).NotTo(HaveOccurred())
		Expect(p1a).NotTo(Equal(p2))
	})

	It("should successfully authenticate with generated username and password", func() {
		u, p, err := a.Generate("qwertyuiop1234567890")

		Expect(err).NotTo(HaveOccurred())
		Expect(a.Valid(u, p)).To(BeTrue())
	})
})
