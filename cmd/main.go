package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"

	config "git.ctisoftware.vn/back-end/base/config"
	"git.ctisoftware.vn/back-end/base/src/database"
	"git.ctisoftware.vn/back-end/base/src/server"
	"git.ctisoftware.vn/back-end/base/src/utilities"
	ctiLog "git.ctisoftware.vn/back-end/utilities/data/provider/log"
)

var (
	configPrefix string
	configSource string
)

func main() {
	app := cli.NewApp()
	app.Name = "Auth microservice"
	app.Usage = "Auth microservice"
	app.Copyright = "Copyright © 2024 CTI Groups. All Rights Reserved."
	app.Compiled = time.Now()

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "configPrefix",
			Aliases:     []string{"confPrefix"},
			Usage:       "prefix for config",
			Value:       "auth",
			Destination: &configPrefix,
		},
		&cli.StringFlag{
			Name:        "configSource",
			Aliases:     []string{"confSource"},
			Value:       "../config/.env",
			Usage:       "set path to environment file",
			Destination: &configSource,
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:   "serve",
			Usage:  "Start the auth server",
			Action: Serve,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "addr-graph",
					Aliases: []string{"address-graph"},
					Value:   "0.0.0.0:8080",
					Usage:   "address for serve graph",
				},
				&cli.StringFlag{
					Name:    "addr-grpc",
					Aliases: []string{"address-grpc"},
					Value:   "0.0.0.0:8090",
					Usage:   "address for serve grpc",
				},
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		return config.LoadFromEnv(configPrefix, configSource)
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endSignal := make(chan os.Signal, 1)
	signal.Notify(endSignal, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func(ctx context.Context, errChan chan error) {
		err := app.RunContext(ctx, os.Args)
		errChan <- err
	}(ctx, errChan)

	select {
	case sign := <-endSignal:
		log.Println("shutting down. reason:", sign)
		return
	case err := <-errChan:
		if err == nil {
			return
		}
		log.Println("encountered error:", err)
		return
	}
}

func Serve(c *cli.Context) error {
	serviceLogger, err := ctiLog.NewLoggerProvider(ctiLog.VerbosityLevel(config.Get().LogLevel), config.Get().ServiceLogFile)
	if err != nil {
		return err
	}
	defer serviceLogger.GetLogger().Sync()
	utilities.SetLogger(serviceLogger.GetLogger())

	ctx := c.Context
	err = database.ConnectDatabse(ctx)
	if err != nil {
		panic(err)
	}

	go func() {
		err = server.ServeGrpc(ctx, c.String("addr-grpc"))
		if err != nil {
			panic(err)
		}
	}()

	return server.ServeGraph(c.Context, c.String("addr-graph"))
}
