package cmd

import (
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// GlobalFlags are going to store some data used by all of the commands.
var GlobalFlags struct {
	Debug             bool
	RedisAddress      string
	StructuredLogging bool
}

// ConfigureGlobals should fill in the above struct with usable data values.
func ConfigureGlobals(app *kingpin.Application) {
	app.Flag("debug", "Verbose mode of running the IAP").
		Short('v').
		OverrideDefaultFromEnvar("DEBUG").
		BoolVar(&GlobalFlags.Debug)

	app.Flag("structured-logging", "Verbose mode of running the IAP").
		Short('J').
		OverrideDefaultFromEnvar("STRUCTURED_LOGGING").
		BoolVar(&GlobalFlags.StructuredLogging)

	app.Flag("redis-address", "The address Redis server is running under.").
		Short('R').
		Default("127.0.0.1:6379").
		OverrideDefaultFromEnvar("REDIS_ADDRESS").
		StringVar(&GlobalFlags.RedisAddress)
}
