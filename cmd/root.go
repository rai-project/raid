package cmd

import (
	"os"
	"runtime/pprof"
	"syscall"

	"github.com/fatih/color"
	"github.com/rai-project/cmd"
	"github.com/rai-project/config"
	"github.com/rai-project/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vrecan/death"
)

var (
	AppSecret  string
	isColor    bool
	isVerbose  bool
	isDebug    bool
	inShutdown int32
)

type prof struct{}

func (c prof) Close() error {
	pprof.StopCPUProfile()
	return nil
}

func serverOptions() []server.Option {
	return []server.Option{
		server.Stdout(os.Stdout),
		server.Stderr(os.Stderr),
		server.NumWorkers(1),
		server.JobQueueName("rai"),
	}
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:          "raid",
	Short:        "The server is used to accept jobs from the rai queue.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		death := death.NewDeath(syscall.SIGINT, syscall.SIGTERM)

		server, err := server.New(serverOptions()...)
		if err != nil {
			return err
		}

		if err := server.Connect(); err != nil {
			return err
		}

		if config.IsDebug || config.IsVerbose {
			death.SetLogger(log)
		}

		death.WaitForDeath(server, prof{})

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
	RootCmd.AddCommand(cmd.GendocCmd)
	RootCmd.AddCommand(cmd.CompletionCmd)

	RootCmd.PersistentFlags().StringVarP(&AppSecret, "secret", "s", "", "Pass in application secret.")
	RootCmd.PersistentFlags().BoolVarP(&isColor, "color", "c", true, "Toggle color output.")
	RootCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Toggle verbose mode.")
	RootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "Toggle debug mode.")

	// mark secret flag hidden
	RootCmd.PersistentFlags().MarkHidden("secret")

	// viper.BindPFlag("app.secret", RootCmd.PersistentFlags().Lookup("secret"))
	viper.BindPFlag("app.debug", RootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("app.verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("app.color", RootCmd.PersistentFlags().Lookup("color"))
}

func initConfig() {
	config.Init(
		config.AppName("raid"),
		config.AppSecret(AppSecret),
		config.ConfigFileBaseName(".rai_config"),
	)
}

func initColor() {
	color.NoColor = !isColor
}
