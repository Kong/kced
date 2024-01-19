/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/kong/go-apiops/logbasics"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-apiops",
	Short: "A CLI for testing the Kong go-apiops library",
	Long: `A CLI for testing the Kong go-apiops library.

go-apiops houses an improved APIOps toolset for operating Kong Gateway deployments.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// set the verbosity level of the log output
		verbosity, err := cmd.Flags().GetInt("verbose")
		if err != nil {
			return err
		}
		logbasics.Initialize(log.LstdFlags, verbosity)
		return nil
	},

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().Int("verbose", 0,
		"this value sets the verbosity level of the log output (higher == more verbose)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
