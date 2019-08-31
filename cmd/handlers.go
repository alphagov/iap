package cmd

import (
	"net/http"

	"github.com/alphagov/iap/internal"
	"github.com/alphagov/iap/pkg/auth"
	"github.com/sirupsen/logrus"
)

type healthcheckResponse struct {
	Redis bool `json:"redis"`
}

type credentialResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func healthcheckHandler(ctx internal.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		healthyRedis := true

		_, err := ctx.Redis.Ping().Result()
		if err != nil {
			ctx.Logger.WithField("redis", GlobalFlags.RedisAddress).Warnln("unable to connect with redis")
			status = http.StatusInternalServerError
			healthyRedis = false
		}

		internal.JSONResponse(ctx, w, status, healthcheckResponse{
			Redis: healthyRedis,
		})
	}
}

func generateSOCKS5Credentials(ctx internal.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := auth.New(ctx.Redis, ctx.Logger)
		username, password, err := client.Generate()
		if err != nil {
			ctx.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("failed to generate socks5 credentials")
			internal.JSONResponse(ctx, w, http.StatusInternalServerError, map[string]string{
				"error": "unable to generate credentials",
			})
			return
		}

		ctx.Logger.WithFields(logrus.Fields{
			"username": username,
		}).Debug("generated new socks5 user")

		internal.JSONResponse(ctx, w, http.StatusOK, credentialResponse{
			Username: username,
			Password: password,
		})
	}
}
