package main

import (
	"os"

	"github.com/alphagov/iap/cmd"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	Version = "0.0.1"
)

func main() {
	run(os.Args[1:], os.Exit)
}

func run(args []string, exit func(int)) {
	app := kingpin.New(
		"iap",
		"GDS - Identity Aware Proxy",
	)

	app.Writer(os.Stdout)
	app.Version(Version)
	app.Terminate(exit)

	cmd.ConfigureGlobals(app)
	cmd.ConfigureWebCommand(app)
	cmd.ConfigureProxyCommand(app)
	cmd.ConfigureSocksCommand(app)

	kingpin.MustParse(app.Parse(args))
}
