package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alphagov/iap/internal"
	"github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// WebCommandInput is a configuration only to be used by this particular command.
type WebCommandInput struct {
	Port uint16
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
	mux.HandleFunc("/healthcheck", healthcheckHandler(ctx))
	mux.HandleFunc("/socks5/generate", generateSOCKS5Credentials(ctx))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 1 * time.Minute,
		IdleTimeout:  15 * time.Second,
	}

	ctx.Logger.WithFields(logrus.Fields{
		"port": cfg.Port,
	}).Info("starting web server")

	return srv.ListenAndServe()
}
