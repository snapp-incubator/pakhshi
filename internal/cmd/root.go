package cmd

import (
	"os"

	"github.com/1995parham/pakhshi/internal/cmd/publish"
	"github.com/1995parham/pakhshi/internal/cmd/subscribe"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// ExitFailure status code.
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// nolint: exhaustivestruct
	root := &cobra.Command{
		Use: "pakhshi",
	}

	brokers := new([]string)
	root.PersistentFlags().StringSliceVarP(brokers, "brokers", "b",
		[]string{"tcp://127.0.0.1:1883"}, "brokers e.g. tcp://127.0.0.1:1883")

	subscribe.Register(root, brokers)
	publish.Register(root, brokers)

	if err := root.Execute(); err != nil {
		pterm.Error.Printf("failed to execute root command: %s", err.Error())
		os.Exit(ExitFailure)
	}
}
