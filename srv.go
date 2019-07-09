package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"

	"github.com/tfwio/sekhm/fsindex/config"
	"github.com/tfwio/sekhm/util"
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
	mCli.Commands = []cli.Command{
		cli.Command{
			Name:        "run",
			Action:      func(*cli.Context) { initialize() },
			Usage:       "Runs the server.",
			Description: "Default operation.",
			Aliases:     []string{"go"},
			Flags:       []cli.Flag{},
		},
		cli.Command{
			Name:        "make-conf",
			Description: "Generate configuration file: <[file-path].json>.",
			Usage:       fmt.Sprintf("%s make-conf <[file-path].json>", filepath.Base(os.Args[0])),
			Flags:       []cli.Flag{},
			Action: func(clictx *cli.Context) {
				if clictx.NArg() == 0 {
					fmt.Println("- supply a file-name to generate.\nI.E. \"conf.json\"")
					os.Exit(0)
				}
				fmt.Printf("- found %s\n", util.Abs(clictx.Args().First()))
				thearg := clictx.Args().First()
				input := util.Abs(clictx.Args().First())
				if util.FileExists(input) {
					fmt.Printf("- please delete the file (%s) before calling this command\n", thearg)
					os.Exit(0)
				}
				configuration.InitializeDefaults(defaultConfPath, defaultConfTarget)
				configuration.ToJSON(input)
			},
		}}
	mCli.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "tls",
			Destination: &config.UseTLS,
			Usage:       "Wether or not to use TLS.\n\tNote: if set, overrides (JSON) conf settings.",
		},
		cli.StringFlag{
			Name:        "host",
			Destination: &config.UseHost,
			Usage:       "UseHost is identifies the host to use to over-ride JSON config.",
		},
		cli.IntFlag{
			Name:        "port",
			Destination: &config.UsePORT,
			Usage:       "UseHost is identifies the host to use to over-ride JSON config.",
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

	configuration.InitializeDefaults(defaultConfPath, defaultConfTarget)
	configuration.FromJSON(config.DefaultConfigFile) // loads (or creates conf.json and terminates application)
	configuration.TLS = configuration.DoTLS()
	if config.UseHost != "" {
		configuration.Server.Host = config.UseHost
	}
	if config.UsePORT != -1 {
		configuration.Server.Port = fmt.Sprintf(":%d", config.UsePORT)
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	configuration.GinConfig(router)

	if configuration.TLS {
		println("- TLS on")
		if err := router.RunTLS(configuration.Server.Port, configuration.Server.Crt, configuration.Server.Key); err != nil {
			panic(fmt.Sprintf("router error: %s\n", err))
		}
	} else {
		println("- TLS off")
		if err := router.Run(configuration.Server.Port); err != nil {
			panic(fmt.Sprintf("router error: %s\n", err))
		}
	}
}

const (
	defaultConfPath   = "multi-media\\public"
	defaultConfTarget = "v"
)
