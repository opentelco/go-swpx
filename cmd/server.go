package cmd

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"git.liero.se/opentelco/go-swpx/config"
	"git.liero.se/opentelco/go-swpx/core"
	"git.liero.se/opentelco/go-swpx/core/api"
	"git.liero.se/opentelco/go-swpx/database"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	Root.AddCommand(serverCmd)
	serverCmd.AddCommand(StartCmd)
	StartCmd.Flags().StringP("config", "c", "config.hcl", "the config file to use")
	if err := StartCmd.MarkFlagFilename("config", "hcl", "hcl"); err != nil {
		panic(err)
	}

}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server commands for swpx",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
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

		mongoClient, err := database.New(*appConfig.MongoServer, logger.Named("mongodb"))
		if err != nil {
			logger.Warn("could not establish mongodb connection", "error", err)
			logger.Info("no mongo connection established", "cache_enabled", false)
		}

		c, err := core.New(&appConfig, mongoClient, logger)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		if err := c.Start(); err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		commander, err := core.NewCommander(c, logger)
		if err != nil {
			cmd.PrintErr("could not create commander service: ", err)
			os.Exit(1)
		}

		// start API endpoint and add the queue
		// the queue is initated in the core and n workers takes request from it

		// HTTP
		server := api.NewServer(c, logger)

		go func() {
			listner, err := net.Listen("tcp", appConfig.HttpAddr)
			if err != nil {
				cmd.PrintErr(err)
				os.Exit(1)
			}

			err = server.Serve(listner)
			if err != nil {
				cmd.PrintErr(err)
				os.Exit(1)
			}

		}()

		// GRPC Core
		logger.Info("starting grpc listner", "address", appConfig.GrpcAddr)
		lis, err := net.Listen("tcp", appConfig.GrpcAddr)
		if err != nil {
			cmd.PrintErrf("failed to listen: %s", err)
			os.Exit(1)
		}
		grpcServer := grpc.NewServer()
		api.NewGrpc(c, grpcServer, logger)
		api.NewCommanderGrpc(commander, grpcServer)

		go func() {
			logger.Debug("serve grpc server", "address", appConfig.GrpcAddr)
			err = grpcServer.Serve(lis)
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
