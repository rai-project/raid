// +build disabled

package cmd

import (
	"github.com/coreos/go-systemd/activation"
	"github.com/pkg/errors"
	"github.com/rai-project/server"
	"github.com/spf13/cobra"
)

var DaemonCmd = &cobra.Command{
	Use:          "daemon",
	Short:        "Starts the server in daemon mode.",
	Hidden:       false,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := serverOptions()

		server, err := server.New(opts...)
		if err != nil {
			return err
		}

		listeners, err := activation.Listeners(true)
		if err != nil {
			return errors.Wrap(err, "cannot create a systemd listiner")
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(DaemonCmd)
}
