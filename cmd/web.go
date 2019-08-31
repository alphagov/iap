package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alphagov/iap/internal"
	"github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	routeDefault        = "/"
	routeHealthcheck    = "/healthcheck"
	routeCredentials    = "/credentials"
	routeGoogleLogin    = "/auth/google/login"
	routeGoogleCallback = "/auth/google/callback"
)

// WebCommandInput is a configuration only to be used by this particular command.
type WebCommandInput struct {
	ExternalURL string
	Insecure    bool
	Port        uint16

	GoogleClientID     string
	GoogleClientSecret string
}

// ConfigureWebCommand should fill in the above input struct with some usable values.
// It is also responsible for creating the context to be used within the app itself.
func ConfigureWebCommand(app *kingpin.Application) {
	input := WebCommandInput{}
	cmd := app.Command("web", "Web frontend for IAP.")

	cmd.Flag("port", "Port the web server will be operating on.").
		Short('p').
		Default("8080").
		OverrideDefaultFromEnvar("PORT").
		Uint16Var(&input.Port)

	cmd.Flag("insecure", "Run the service without TLS.").
		Short('k').
		OverrideDefaultFromEnvar("INSECURE").
		BoolVar(&input.Insecure)

	cmd.Flag("external-url", "The external URL used for Oauth2 dance.").
		Short('u').
		OverrideDefaultFromEnvar("IAP_WEB_ADDRESS").
		StringVar(&input.ExternalURL)

	cmd.Flag("google-client-id", "The Google client ID used for Oauth2 dance.").
		OverrideDefaultFromEnvar("GOOGLE_CLIENT_ID").
		StringVar(&input.GoogleClientID)

	cmd.Flag("google-client-secret", "The Google client secret used for Oauth2 dance.").
		OverrideDefaultFromEnvar("GOOGLE_CLIENT_SECRET").
		StringVar(&input.GoogleClientSecret)

	cmd.Action(func(c *kingpin.ParseContext) error {
		ctx := internal.Context{
			Logger: internal.SetupLogger(GlobalFlags.StructuredLogging, GlobalFlags.Debug),
			Redis:  internal.SetupRedis(GlobalFlags.RedisAddress),
		}
		return WebCommand(ctx, input)
	})
}

// WebCommand is the main brain behind this commands. It will start the web frontend
// and hang tight accepting, rejecting and working with requests.
// It will take the job of authenticating with OIDC and generating bunch of secrets
// for the IAP users.
func WebCommand(ctx internal.Context, cfg WebCommandInput) error {
	mux := http.DefaultServeMux
	mux.HandleFunc(routeHealthcheck, healthcheckHandler(ctx))
	mux.HandleFunc(routeDefault, defaultHandler(ctx))
	mux.HandleFunc(routeCredentials, authenticatedMiddleware(ctx, getCredentialsHanlder(ctx)))
	mux.HandleFunc(routeGoogleLogin, googleLoginHanlder(ctx, cfg))
	mux.HandleFunc(routeGoogleCallback, googleCallbackHanlder(ctx, cfg))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 1 * time.Minute,
		IdleTimeout:  15 * time.Second,
	}

	ctx.Logger.WithFields(logrus.Fields{
		"port":        cfg.Port,
		"externalURL": cfg.ExternalURL,
		"google":      len(cfg.GoogleClientID) > 0,
		"insecure":    cfg.Insecure,
	}).Info("starting web server")

	return srv.ListenAndServe()
}
