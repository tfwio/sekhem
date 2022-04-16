// srv.go is the main entry-point into the server application.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"

	"github.com/tfwio/srv/fsindex/config"
	"github.com/tfwio/srv/util"
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

func misterActionFunction(*cli.Context) error {
	initialize(true)
	return nil
}

func initializeCli() {
	// mCli := cli.NewApp()
	mCli.Name = filepath.Base(os.Args[0])
	mCli.Authors = []cli.Author{{Name: "tfw; et alia" /*, Email: "tfwroble@gmail.com"}, cli.Author{Name: "Et al."*/}}
	mCli.Version = "v0.0.0a"
	mCli.Copyright = "github.com/tfwio/srv\n\n   This is free, open-source software.\n   disclaimer: use at own risk."

	mCli.Action = func(*cli.Context) error {
		initialize(true)
		return nil
	}

	mCli.Commands = []cli.Command{
		{
			Name: "run",
			Action: func(*cli.Context) error {
				initialize(true)
				return nil
			},
			Usage:       "Runs the server.",
			Description: "Default operation.",
			Aliases:     []string{"go"},
			Flags:       []cli.Flag{},
		},
		{
			Name:        "make-conf",
			Description: "Generate configuration file: <[file-path].json>.",
			Usage:       fmt.Sprintf("%s make-conf <[file-path].json>", filepath.Base(os.Args[0])),
			Flags:       []cli.Flag{},
			Action: func(clictx *cli.Context) error {
				if clictx.NArg() == 0 {
					fmt.Println("- supply a file-name to generate.\nI.E. \"conf.json\"")
					os.Exit(0)
					return nil
				}
				fmt.Printf("- found %s\n", util.Abs(clictx.Args().First()))
				thearg := clictx.Args().First()
				input := util.Abs(clictx.Args().First())
				if util.FileExists(input) {
					fmt.Printf("- please delete the file (%s) before calling this command\n", thearg)
					os.Exit(0)
					return nil
				}
				configuration.InitializeDefaults(defaultConfPath, defaultConfTarget)
				configuration.ToJSON(input)
				return nil
			},
		}}
	mCli.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "tls",
			Destination: &config.UseTLS,
			Usage:       "Sets TLS on.  This only works if/when tls is set in conf.json to false, and if you have valid tls cert/key files wired into the configuration.",
		}, &cli.StringFlag{
			Name:        "host",
			Destination: &config.UseHost,
			Usage:       "UseHost is identifies the host to use to over-ride JSON config.",
		}, &cli.UintFlag{
			Name:        "port",
			Destination: &config.UsePORT,
			Value:       5500,
			Usage:       "UseHost is identifies the host to use to over-ride JSON config.",
		}, &cli.StringFlag{
			Name:        "conf",
			Usage:       "Points to a custom configuration file.",
			Value:       config.DefaultConfigFile,
			Destination: &config.DefaultConfigFile,
		}}
	mCli.Run(os.Args)
}

// initializeJSONConf loads JSON configuration and
// sets our data source.  No file indexing.
//
// From here we can execute database operations.
func initializeJSONConf() {

	configuration.InitializeDefaults(defaultConfPath, defaultConfTarget)
	configuration.FromJSON(config.DefaultConfigFile) // loads (or creates conf.json and terminates application)
	configuration.TLS = configuration.DoTLS()

	if config.UseHost != "" {
		configuration.Host = config.UseHost
	}
	if config.UsePORT != defaultPort {
		configuration.Port = fmt.Sprintf(":%d", config.UsePORT)
	}

}

// initialize can be called with or without starting the server.
// First, this function loads JSON conf followed by
// building up all the file-indexes and finally running
// the gin-server.
func initialize(andServe bool) {

	initializeJSONConf()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	configuration.GinConfigure(andServe, router)

	if !andServe {
		return
	}

	if configuration.TLS {
		println("- TLS on")
		if err := router.RunTLS(configuration.Port, configuration.Crt, configuration.Key); err != nil {
			panic(fmt.Sprintf("router error: %s\n", err))
		}
	} else {
		println("- TLS off")
		if err := router.Run(configuration.Port); err != nil {
			panic(fmt.Sprintf("router error: %s\n", err))
		}
	}
}

const (
	defaultPort       uint = 5500
	defaultConfPath        = "multi-media\\public"
	defaultConfTarget      = "v"
)
