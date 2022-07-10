package cmd

import (
	"chatty/config"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	cfg *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "chatty",
	Short: "chatty - a simple terminal based chat",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(1)
	},
}

func Execute() {
	var err error

	cfg, err = config.Load()
	if err != nil {
		log.Err(err).Msg("failed to load config")
	}

	log.Debug().Interface("config", cfg).Msg("config loaded")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: '%s'", err)
		os.Exit(1)
	}
}
