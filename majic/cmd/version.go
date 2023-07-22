/*
Copyright © 2023 Mark Johnson
*/
package cmd

import (
	"github.com/shelterbelt/majic-cli/majic/helpers/core"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"ver"},
	Short:   "version informataion for the application",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		appConfig := core.Config()
		appConfig.ApplyFlags(cmd.Flags())
		core.Output().NormalOutput("majic CLI 1.0.0")
		core.Output().NormalOutput("Plugins directory: " + appConfig.Get(core.ConfigKeyPluginsDir).(string))
		core.Output().NormalOutput("Input directory: " + appConfig.Get(core.ConfigKeyInputDir).(string))
		core.Output().NormalOutput("Output directory: " + appConfig.Get(core.ConfigKeyOutputDir).(string))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
