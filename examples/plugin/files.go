/*
Copyright © 2023 Mark Johnson
*/
package main

import (
	"github.com/shelterbelt/majic-cli/examples/plugin/myplugin/helpers/myhelper"

	"github.com/shelterbelt/majic-cli/majic/helpers/core"
	"github.com/shelterbelt/majic-cli/majic/helpers/file"
	"github.com/spf13/cobra"
)

// FilesCmd represents the files command
var FilesCmd = &cobra.Command{
	Use:     "files",
	Aliases: []string{},
	Short:   "change \"the\" to \"THE\" in text files",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		core.Config().ApplyFlags(cmd.Flags())
		appConfig := core.Config()
		appConfig.ApplyFlags(cmd.Flags())

		myProcessor := new(myhelper.MyFileProcessor)
		myProcessor.Initialize(cmd.Flags())

		if len(args) >= 1 {
			path := args[0]
			file.ProcessPath(path, myProcessor)
		} else {
			core.Output().NormalOutput("No file or directory to scan specified.")
		}
	},
}

func init() {
	core.Output().VerboseOutput("Initializing command files.go")
	// Normally, this is where the command is registered with the root command
	// instance.  For majic plugins, the command will be registered with
	// the CLI via the plugin's Register method.

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genindexCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// GenIndexCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
