package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"

	"tfw.io/Go/fsindex/fsindex/config"
)

// Configuration variables

var (
	configuration config.Configuration
	mCli          cli.App

	xCounter  int32
	fCounter  int32
	xpCounter *int32
)

func main() {
	initializeCli()
}

func initializeCli() {
	mCli.Name = filepath.Base(os.Args[0])
	mCli.Authors = []cli.Author{cli.Author{Name: "tfw; et alia" /*, Email: "tfwroble@gmail.com"}, cli.Author{Name: "Et al."*/}}
	mCli.Version = "v0.0.0"
	mCli.Copyright = "tfwio.github.com/go-fsindex\n\n   This is free, open-source software.\n   disclaimer: use at own risk."
	mCli.Action = func(*cli.Context) { initialize() }
	mCli.Commands = []cli.Command{cli.Command{
		Name:        "run",
		Action:      func(*cli.Context) { initialize() },
		Usage:       "Runs the server.",
		Description: "Default operation.",
		Aliases:     []string{"go"},
		Flags:       []cli.Flag{},
	}}
	mCli.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "tls",
			Destination: &config.UseTLS,
			Usage:       "Wether or not to use TLS.\n\tNote: if set, overrides (JSON) conf settings.",
		},
		cli.StringFlag{
			Name:        "conf",
			Usage:       "Points to a custom configuration file.",
			Value:       config.DefaultConfigFile,
			Destination: &config.DefaultConfigFile,
		},
	}
	mCli.Run(os.Args)
}

func initialize() {

	serv, tempPath := "v", `multi-media\\public`

	configuration.InitializeDefaults(tempPath, serv)
	configuration.FromJSON() // loads (or creates conf.json and terminates application)

	gin.SetMode(gin.ReleaseMode)
	mGin := gin.Default()

	configuration.GinConfig(mGin)

	if configuration.DoTLS() {
		println("- TLS on")
		if err := mGin.RunTLS(configuration.Server.Port, configuration.Server.Crt, configuration.Server.Key); err != nil {
			panic(fmt.Sprintf("router error: %s\n", err))
		}
	} else {
		println("- TLS off")
		if err := mGin.Run(configuration.Server.Port); err != nil {
			panic(fmt.Sprintf("router error: %s\n", err))
		}
	}
}
