package internal_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/alphagov/iap/internal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Internal packages", func() {
	Context("json", func() {
		var (
			w   *httptest.ResponseRecorder
			l   *bytes.Buffer
			ctx internal.Context
		)

		BeforeEach(func() {
			w = httptest.NewRecorder()
			l = bytes.NewBuffer([]byte{})

			logger := internal.SetupLogger(true, false)
			logger.SetOutput(l)
			ctx = internal.Context{
				Logger: logger,
			}
		})

		It("should compose json response", func() {
			internal.JSONResponse(ctx, w, http.StatusOK, map[string]bool{"test": true})
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(l.String()).NotTo(ContainSubstring("failed to parse json blob"))
			Expect(w.Body.String()).To(MatchJSON(`{"test": true}`))
		})

		It("should fail compose json response", func() {
			internal.JSONResponse(ctx, w, http.StatusOK, make(chan int))
			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			Expect(l.String()).To(ContainSubstring("failed to parse json blob"))
		})
	})

	Context("logger", func() {
		var b *bytes.Buffer

		BeforeEach(func() {
			b = bytes.NewBuffer([]byte{})
		})

		It("should setup logger correctly", func() {
			l := internal.SetupLogger(true, false)
			l.SetOutput(b)
			l.WithField("test", true).Infoln("info line")
			l.WithField("test", true).Debugln("debug line")

			out := b.String()
			Expect(out).To(ContainSubstring("info line"))
			Expect(out).NotTo(ContainSubstring("debug line"))
		})

		It("should setup logger correctly with debug enabled", func() {
			l := internal.SetupLogger(true, true)
			l.SetOutput(b)
			l.WithField("test", true).Infoln("info line")
			l.WithField("test", true).Debugln("debug line")

			out := b.String()
			Expect(out).To(ContainSubstring("info line"))
			Expect(out).To(ContainSubstring("debug line"))
		})
	})

	Context("redis", func() {
		var (
			mr *miniredis.Miniredis
		)

		BeforeSuite(func() {
			var err error
			mr, err = miniredis.Run()
			Expect(err).NotTo(HaveOccurred())
		})

		AfterSuite(func() {
			mr.Close()
		})

		It("should fail to setup redis connection", func() {
			r := internal.SetupRedis("localhost.test:73198246")
			Expect(r.Ping().Err()).To(HaveOccurred())
		})

		It("should setup redis connection correctly", func() {
			r := internal.SetupRedis(mr.Addr())
			Expect(r.Ping().Err()).NotTo(HaveOccurred())
			Expect(r.Set("test", "1234567890", time.Hour).Err()).NotTo(HaveOccurred())

			v, err := r.Get("test").Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(v).To(Equal("1234567890"))
		})
	})
})
