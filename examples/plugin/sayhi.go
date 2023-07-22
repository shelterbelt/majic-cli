/*
Copyright © 2023 Mark Johnson
*/
package main

import (
	"github.com/shelterbelt/majic-cli/majic/helpers/core"
	"github.com/spf13/cobra"
)

// SayHiCmd represents the sayhi command
var SayHiCmd = &cobra.Command{
	Use:     "sayhi",
	Aliases: []string{"hi"},
	Short:   "sample CLI command that says 'hi'",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		core.Config().ApplyFlags(cmd.Flags())

		name := "World"
		if len(args) >= 1 {
			name = args[0]
		}

		excited, _ := cmd.Flags().GetBool("excited")
		if excited {
			core.Output().NormalOutput("Hi " + name + "!!!!!")
		} else {
			core.Output().NormalOutput("Hello " + name + ".")
		}
	},
}

func init() {
	core.Output().VerboseOutput("Initializing command sayhi.go")
	// Normally, this is where the command is registered with the root command
	// instance.  For majic plugins, the command will be registered with
	// the CLI via the plugin's Register method.

	// Here you will define your flags and configuration settings.
	SayHiCmd.Flags().BoolP("excited", "e", false, "say it like you mean it!")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// formatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// formatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
