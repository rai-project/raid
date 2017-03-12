package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/rai-project/cmd"
	"github.com/rai-project/config"
	"github.com/rai-project/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	isColor    bool
	isVerbose  bool
	isDebug    bool
	inShutdown int32
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:          "raid",
	Short:        "The server is used to accept jobs from the rai queue.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		signalNotify(interrupt)

		server, err := server.New(
			server.Stdout(os.Stdout),
			server.Stderr(os.Stderr),
			server.NumWorkers(1),
			server.JobQueueName("rai"),
		)
		if err != nil {
			return err
		}

		quitting := make(chan struct{})
		go handleInterrupt(interrupt, quitting, server)

		if err := server.Connect(); err != nil {
			if err != nil {
				// If the underlying listening is closed, Serve returns an error
				// complaining about listening on a closed socket. This is expected, so
				// let's ignore the error if we are the ones who explicitly closed the
				// socket.
				select {
				case <-quitting:
					return nil
				default:
					return err
				}
			}
		}
		<-quitting
		println("quit")
		return nil
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initColor)

	RootCmd.AddCommand(cmd.VersionCmd)
	RootCmd.AddCommand(cmd.LicenseCmd)
	RootCmd.AddCommand(cmd.EnvCmd)

	RootCmd.PersistentFlags().StringVarP(&appsecret, "secret", "s", "", "Pass in application secret.")
	RootCmd.PersistentFlags().BoolVarP(&isColor, "color", "c", color.NoColor, "Toggle color output.")
	RootCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Toggle verbose mode.")
	RootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "Toggle debug mode.")

	// mark secret flag hidden
	RootCmd.PersistentFlags().MarkHidden("secret")

	viper.BindPFlag("app.secret", RootCmd.PersistentFlags().Lookup("secret"))
	viper.BindPFlag("app.debug", RootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("app.verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("app.color", RootCmd.PersistentFlags().Lookup("color"))
}

func initConfig() {
	config.Init(
		config.AppName("raid"),
		config.AppSecret(appsecret),
	)
}

func initColor() {
	color.NoColor = !isColor
}
