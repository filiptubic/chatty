package cmd

import (
	"chatty/config"
	"chatty/server"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start chatty server",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return err
		}

		s := server.New(cfg)
		return s.Start()
	},
}
