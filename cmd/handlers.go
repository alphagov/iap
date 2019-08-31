package cmd

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/alphagov/iap/internal"
	"github.com/alphagov/iap/pkg/auth"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type healthcheckResponse struct {
	Redis bool `json:"redis"`
}

type credentialResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	oauthGoogleURLAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	stateCookieName   = "iap-oauth-state"
	tokenCookieName   = "iap-access-token"
)

func defaultHandler(ctx internal.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<a href=/auth/google/login>Google</a>"))
	}
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

func getCredentialsHanlder(ctx internal.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uniqueCode := fmt.Sprintf("iap-%d", time.Now().Unix())
		userData, err := obtainUserDataFromGoogle(r)
		if err != nil {
			ctx.Logger.WithFields(logrus.Fields{"error": err, "code": uniqueCode}).Error("couldn't retreive userData")
			internal.JSONResponse(ctx, w, http.StatusInternalServerError, logrus.Fields{"error": "unable to generate credentials", "code": uniqueCode})
			return
		}
		client := auth.New(ctx.Redis, ctx.Logger)
		username, password, err := client.Generate(userData.ID)
		if err != nil {
			ctx.Logger.WithFields(logrus.Fields{"error": err, "code": uniqueCode}).Error("failed to generate socks5 credentials")
			internal.JSONResponse(ctx, w, http.StatusInternalServerError, logrus.Fields{"error": "unable to generate credentials", "code": uniqueCode})
			return
		}

		ctx.Logger.WithFields(logrus.Fields{"email": userData.Email, "username": username}).Debug("generated new socks5 user")
		internal.JSONResponse(ctx, w, http.StatusOK, credentialResponse{Username: username, Password: password})
	}
}

func googleLoginHanlder(ctx internal.Context, cfg WebCommandInput) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uniqueCode := fmt.Sprintf("iap-%d", time.Now().Unix())
		conf, err := googleConfig(cfg)
		if err != nil {
			ctx.Logger.WithFields(logrus.Fields{"error": err, "code": uniqueCode}).Error("unable to compose google configuration")
			internal.JSONResponse(ctx, w, http.StatusInternalServerError, logrus.Fields{"error": "misconfigured authentication", "code": uniqueCode})
			return
		}
		nonce, err := auth.Nonce(32)
		if err != nil {
			ctx.Logger.WithFields(logrus.Fields{"error": err, "code": uniqueCode}).Error("unable to generate nonce")
			internal.JSONResponse(ctx, w, http.StatusInternalServerError, logrus.Fields{"error": "issues with the system", "code": uniqueCode})
			return
		}
		u := conf.AuthCodeURL(nonce)
		http.SetCookie(w, &http.Cookie{Secure: !cfg.Insecure, Name: stateCookieName, Value: nonce, Expires: time.Now().Add(time.Minute), MaxAge: 60})
		ctx.Logger.WithField("dialog", u).Debug("redirecting for auth dialog")
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	}
}

func googleCallbackHanlder(ctx internal.Context, cfg WebCommandInput) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uniqueCode := fmt.Sprintf("iap-%d", time.Now().Unix())
		conf, err := googleConfig(cfg)
		if err != nil {
			ctx.Logger.WithFields(logrus.Fields{"error": err, "code": uniqueCode}).Error("unable to compose google configuration")
			internal.JSONResponse(ctx, w, http.StatusInternalServerError, logrus.Fields{"error": "misconfigured authentication", "code": uniqueCode})
			return
		}
		stateCookie, err := r.Cookie(stateCookieName)
		if err != nil {
			ctx.Logger.WithFields(logrus.Fields{"error": err, "code": uniqueCode}).Error("unable to obtain cookie")
			internal.JSONResponse(ctx, w, http.StatusInternalServerError, logrus.Fields{"error": "couldn't access cookie", "code": uniqueCode})
			return
		}
		if r.FormValue("state") != stateCookie.Value {
			ctx.Logger.WithFields(logrus.Fields{"error": err, "code": uniqueCode}).Error("untrusted activity")
			internal.JSONResponse(ctx, w, http.StatusInternalServerError, logrus.Fields{"error": "messed up the order of things", "code": uniqueCode})
			return
		}
		accessToken, err := obtainAccessTokenFromGoogle(conf, r.FormValue("code"))
		if err != nil {
			ctx.Logger.WithFields(logrus.Fields{"error": err, "code": uniqueCode}).Error("unable to obtain accessToken")
			internal.JSONResponse(ctx, w, http.StatusInternalServerError, logrus.Fields{"error": "problems with authentication", "code": uniqueCode})
			return
		}
		ctx.Logger.Debug("got accessToken")
		http.SetCookie(w, &http.Cookie{Secure: !cfg.Insecure, Name: tokenCookieName, Value: accessToken, Path: "/", Expires: time.Now().Add(time.Hour)})
		http.Redirect(w, r, routeCredentials, http.StatusTemporaryRedirect)
	}
}

func obtainAccessTokenFromGoogle(conf *oauth2.Config, code string) (string, error) {
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func googleConfig(cfg WebCommandInput) (*oauth2.Config, error) {
	u, err := url.Parse(cfg.ExternalURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse url: %s", err)
	}
	u.Path = routeGoogleCallback
	return &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  u.String(),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}, nil
}
