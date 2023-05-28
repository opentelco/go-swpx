package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.liero.se/opentelco/go-swpx/api"
	"git.liero.se/opentelco/go-swpx/config"
	"git.liero.se/opentelco/go-swpx/core"
	"git.liero.se/opentelco/go-swpx/core/worker"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/client"
)

func init() {
	Root.AddCommand(StartCmd)
	StartCmd.Flags().StringP("config", "c", "config.hcl", "the config file to use")
	if err := StartCmd.MarkFlagFilename("config", "hcl", "hcl"); err != nil {
		panic(err)
	}

}

// Start starts Switchpoller daemon/server
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "start the swpx daemon",
	Long:  `switchpoller x. the long description of the application`,
	Run: func(cmd *cobra.Command, args []string) {

		configSrc, err := cmd.Flags().GetString("config")
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		appConfig := config.Configuration{}
		err = config.LoadConfig(configSrc, &appConfig)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		logger = hclog.New(&hclog.LoggerOptions{
			Name:            APP_NAME,
			Level:           hclog.LevelFromString(appConfig.Logger.Level),
			Color:           hclog.AutoColor,
			IncludeLocation: true,
			JSONFormat:      appConfig.Logger.AsJson,
		})

		// setup temporal
		opts := client.Options{
			HostPort:  appConfig.Temporal.Address,
			Namespace: appConfig.Temporal.Namespace,
			Logger:    logger,
		}

		tc, err := client.Dial(opts)
		if err != nil {
			cmd.PrintErr(fmt.Errorf("could not create temporal client: %w", err))
			os.Exit(1)
		}

		c, err := core.New(&appConfig, logger)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		// setup and start the temporal worker
		w := worker.New(tc, appConfig.Temporal.TaskQueue, c, logger)
		if err := w.Start(); err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		if err := c.Start(); err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		// start API endpoint and add the queue
		// the queue is initated in the core and n workers takes request from it

		// HTTP
		server := api.NewServer(c, logger)
		go func() {
			err = server.ListenAndServe(appConfig.HttpAddr)
			if err != nil {
				cmd.PrintErr(err)
				os.Exit(1)
			}
		}()
		// GRPC
		grpcServer := api.NewCoreGrpcServer(c, logger)
		go func() {
			err = grpcServer.ListenAndServe(appConfig.GrpcAddr)
			if err != nil {
				cmd.PrintErr(err)
				os.Exit(1)
			}
		}()

		signalChan := make(chan os.Signal, 1)

		signal.Notify(
			signalChan,
			syscall.SIGHUP,  // kill -SIGHUP XXXX
			syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
			syscall.SIGQUIT, // kill -SIGQUIT XXXX
		)

		<-signalChan
		cmd.Println("os.Interrupt - shutting down...")

		go func() {
			<-signalChan
			cmd.Println("os.Kill - terminating...")
			os.Exit(1)
		}()

		// manually cancel context if not using httpServer.RegisterOnShutdown(cancel)

		defer os.Exit(0)
	},
}
