package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alphagov/iap/internal"
	"github.com/alphagov/iap/pkg/auth"
	"github.com/elazarl/goproxy"
	goproxyAuth "github.com/elazarl/goproxy/ext/auth"
	"github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// ProxyCommandInput is a configuration only to be used by this particular command.
type ProxyCommandInput struct {
	Host string
	Port uint16
}

// ConfigureProxyCommand should fill in the above input struct with some usable values.
// It is also responsible for creating the context to be used within the app itself.
func ConfigureProxyCommand(app *kingpin.Application) {
	input := ProxyCommandInput{}
	cmd := app.Command("proxy", "Run a HTTP proxy server.")

	cmd.Flag("host", "Host the HTTP proxy will be available under.").
		Short('H').
		Default("127.0.0.1").
		OverrideDefaultFromEnvar("HOST").
		StringVar(&input.Host)

	cmd.Flag("port", "Port the HTTP proxy will be available under.").
		Short('p').
		Default("8080").
		OverrideDefaultFromEnvar("PORT").
		Uint16Var(&input.Port)

	cmd.Action(func(c *kingpin.ParseContext) error {
		ctx := internal.Context{
			Logger: internal.SetupLogger(GlobalFlags.StructuredLogging, GlobalFlags.Debug),
			Redis:  internal.SetupRedis(GlobalFlags.RedisAddress),
		}
		return ProxyCommand(ctx, input)
	})
}

// ProxyCommand is the main brain behind this commands. It will start the HTTP Proxy server
// and hang tight accepting, rejecting and working with requests.
func ProxyCommand(ctx internal.Context, cfg ProxyCommandInput) error {
	client := auth.New(ctx.Redis, ctx.Logger)

	w := ctx.Logger.Writer()
	defer w.Close()

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = GlobalFlags.Debug
	proxy.Logger = log.New(w, "", 0)

	goproxyAuth.ProxyBasic(proxy, "qwertyuiop", client.Valid)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	ctx.Logger.WithFields(logrus.Fields{
		"address": addr,
	}).Info("starting proxy server")

	return http.ListenAndServe(addr, proxy)
}
