package cmd

import (
	"fmt"
	"log"

	"github.com/alphagov/iap/internal"
	"github.com/alphagov/iap/pkg/auth"
	socks5 "github.com/armon/go-socks5"
	"github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// SocksCommandInput is a configuration only to be used by this particular command.
type SocksCommandInput struct {
	Host     string
	Port     uint16
	Protocol string
}

// ConfigureSocksCommand should fill in the above input struct with some usable values.
// It is also responsible for creating the context to be used within the app itself.
func ConfigureSocksCommand(app *kingpin.Application) {
	input := SocksCommandInput{}
	cmd := app.Command("socks5", "Run a SOCKS5 proxy server.")

	cmd.Flag("host", "Host the SOCKS5 proxy will be available under.").
		Short('H').
		Default("127.0.0.1").
		OverrideDefaultFromEnvar("HOST").
		StringVar(&input.Host)

	cmd.Flag("port", "Port the SOCKS5 proxy will be available under.").
		Short('p').
		Default("1080").
		OverrideDefaultFromEnvar("PORT").
		Uint16Var(&input.Port)

	cmd.Flag("protocol", "Protocol the SOCKS5 proxy will be communicating with.").
		Short('P').
		Default("tcp").
		OverrideDefaultFromEnvar("PROTOCOL").
		StringVar(&input.Protocol)

	cmd.Action(func(c *kingpin.ParseContext) error {
		ctx := internal.Context{
			Logger: internal.SetupLogger(GlobalFlags.Debug),
			Redis:  internal.SetupRedis(GlobalFlags.RedisAddress),
		}
		return SocksCommand(ctx, input)
	})
}

// SocksCommand is the main brain behind this commands. It will start the SOCKS5 Proxy server
// and hang tight accepting, rejecting and working with requests.
func SocksCommand(ctx internal.Context, cfg SocksCommandInput) error {
	client := auth.New(ctx.Redis, ctx.Logger)

	w := ctx.Logger.Writer()
	defer w.Close()

	srv, err := socks5.New(&socks5.Config{
		Logger:      log.New(w, "", 0),
		Credentials: client,
	})
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	ctx.Logger.WithFields(logrus.Fields{
		"address":  addr,
		"protocol": cfg.Protocol,
	}).Info("starting socks5 server")

	return srv.ListenAndServe(cfg.Protocol, addr)
}
