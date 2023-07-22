/*
Copyright © 2023 Mark Johnson
*/
package plugin

import (
	"errors"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/shelterbelt/majic-cli/majic/helpers/core"
	"github.com/spf13/cobra"
)

type MajicPlugin interface {
	Register() []string
}

func LoadPlugins(root *cobra.Command) {
	appConfig := core.Config()
	pluginsDirPath := appConfig.GetD(core.ConfigKeyPluginsDir, core.DefaultPluginsDir).(string)
	pluginsDir, err := os.Open(pluginsDirPath)
	if (err == nil) || (err != nil && !errors.Is(err, os.ErrNotExist)) {
		core.HandleError(err)
		defer pluginsDir.Close()
		files, err := pluginsDir.Readdirnames(-1)
		core.HandleError(err)
		for i := 0; i < len(files); i++ {
			pluginPath := filepath.Join(pluginsDirPath, files[i])
			info, err := os.Stat(pluginPath)
			if err != nil && !errors.Is(err, os.ErrExist) {
				core.Output().NormalOutput("File " + pluginPath + " does not exist.")
			} else {
				core.HandleError(err)
				if !info.IsDir() && strings.HasSuffix(files[i], ".so") {
					plugin, err := plugin.Open(pluginPath)
					core.HandleError(err)
					instanceSym, err := plugin.Lookup("Plugin")
					core.HandleError(err)
					var instance MajicPlugin
					instance, found := instanceSym.(MajicPlugin)
					if found {
						commands := instance.Register()
						registerCommands(plugin, root, commands)
					} else {
						core.Output().NormalOutput("Failed to load plugin: " + pluginPath)
					}

				} else {
					core.Output().VerboseOutput("Skipping non-plugin file: " + files[i])
				}
			}
		}
	}
}

func registerCommands(plugin *plugin.Plugin, root *cobra.Command, commands []string) {
	for j := 0; len(commands) > j; j++ {
		cmdVarSym, err := plugin.Lookup(commands[j])
		if err == nil {
			initSym, err := plugin.Lookup("Init" + commands[j])
			if err == nil {
				initSym.(func())()
			}
			root.AddCommand(*cmdVarSym.(**cobra.Command))
		} else {
			core.Output().NormalOutput("Warning: " + err.Error())
		}
	}
}
