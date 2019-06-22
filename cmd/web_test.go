package cmd

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/alphagov/iap/internal"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("Web server", func() {
	var (
		mr  *miniredis.Miniredis
		ctx internal.Context
	)

	BeforeSuite(func() {
		var err error
		mr, err = miniredis.Run()
		Expect(err).NotTo(HaveOccurred())

		logger := logrus.New()
		logger.SetOutput(GinkgoWriter)

		ctx = internal.Context{
			Logger: logger,
			Redis: redis.NewClient(&redis.Options{
				Addr: mr.Addr(),

				ReadTimeout:  time.Second * 2,
				WriteTimeout: time.Second * 1,
				DialTimeout:  time.Second * 1,
			}),
		}
	})

	AfterSuite(func() {
		ctx.Redis.Close()
		mr.Close()
	})

	It("should pass the healthcheck", func() {
		req, err := http.NewRequest("GET", "/healthcheck", nil)
		Expect(err).NotTo(HaveOccurred())

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(healthcheckHandler(ctx))

		handler.ServeHTTP(rr, req)

		Expect(rr.Code).To(Equal(http.StatusOK))
		Expect(rr.Body.String()).To(MatchJSON(`{"redis": true}`))
	})

	It("should fail the healthcheck due to lack of redis conectivity", func() {
		req, err := http.NewRequest("GET", "/healthcheck", nil)
		Expect(err).NotTo(HaveOccurred())

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(healthcheckHandler(internal.Context{
			Logger: ctx.Logger,
			Redis: redis.NewClient(&redis.Options{
				Addr:        "0.0.0.0:56789",
				DialTimeout: time.Second * 1,
			}),
		}))

		handler.ServeHTTP(rr, req)

		Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		Expect(rr.Body.String()).To(MatchJSON(`{"redis": false}`))
	})

	It("should generate a new set of credentials for user", func() {
		req, err := http.NewRequest("GET", "/socks5/generate", nil)
		Expect(err).NotTo(HaveOccurred())

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(generateSOCKS5Credentials(ctx))

		handler.ServeHTTP(rr, req)

		Expect(rr.Code).To(Equal(http.StatusOK))
		Expect(rr.Body.String()).To(ContainSubstring(`"username"`))
		Expect(rr.Body.String()).To(ContainSubstring(`"password"`))
	})
})
