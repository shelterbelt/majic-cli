/*
Copyright © 2023 Mark Johnson
*/
package cmd

import (
	"os"

	"github.com/shelterbelt/majic-cli/majic/helpers/core"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Permenantly delete temporary and generated output files created by the cli",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		appConfig := core.Config()
		appConfig.ApplyFlags(cmd.Flags())
		outputPath := appConfig.GetD(core.ConfigKeyOutputDir, core.DefaultOutputDir).(string)
		core.Output().NormalOutput("Deleting contents of " + outputPath)
		os.RemoveAll(outputPath)
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
