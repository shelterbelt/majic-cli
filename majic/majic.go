/*
Copyright © 2023 Mark Johnson
*/
package main

import (
	"github.com/shelterbelt/majic-cli/majic/cmd"

	"github.com/shelterbelt/majic-cli/majic/helpers/core"
	"github.com/shelterbelt/majic-cli/majic/helpers/file"
	"github.com/shelterbelt/majic-cli/majic/helpers/plugin"
)

func main() {
	pluginsPath := core.Config().GetD(core.ConfigKeyPluginsDir, core.DefaultPluginsDir).(string)
	inputDirPath := core.Config().GetD(core.ConfigKeyInputDir, core.DefaultInputDir).(string)
	outputDirPath := file.CreateOutputDir()
	core.Output().DetailedOutput("Plugins directory: " + pluginsPath)
	core.Output().DetailedOutput("Input directory: " + inputDirPath)
	core.Output().DetailedOutput("Output directory: " + outputDirPath)

	plugin.LoadPlugins(cmd.GetRootCommand())

	cmd.Execute()
}
