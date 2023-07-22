/*
Copyright © 2023 Mark Johnson
*/
package cmd

import (
	"os"

	"github.com/shelterbelt/majic-cli/majic/helpers/core"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "majic",
	Short: "Because who doesn't need a bespoke CLI to handle day-to-day tasks?",
	Long:  `majic is a general-purpose command-line interface application with a modular architecture to facilitate easily adding new commands.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func GetRootCommand() *cobra.Command {
	return rootCmd
}

func init() {
	// Here you will define your flags and configuration settings.

	rootCmd.PersistentFlags().Bool(core.FlagKeyDetailedOutput, core.DefaultDetailedOutput, "detailed output")
	rootCmd.PersistentFlags().Bool(core.FlagKeyVerboseOutput, core.DefaultVerboseOutput, "verbose output (i.e. everything)")

	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cbrapp.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
