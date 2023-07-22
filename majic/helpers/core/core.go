/*
Copyright © 2023 Mark Johnson
*/
package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/magiconair/properties"
	"github.com/spf13/pflag"
)

type Processor interface {
	Initialize(flags *pflag.FlagSet)
}

const AppConfigHome string = ".majic"
const AppConfigFile string = "clirc"
const ConfigKeyPluginsDir string = "plugins_dir"
const ConfigKeyInputDir string = "input_dir"
const ConfigKeyOutputDir string = "output_dir"
const DefaultPluginsDir string = "${HOME}/.majic/plugins"
const DefaultInputDir string = "${HOME}/.majic/input"
const DefaultOutputDir string = "${HOME}/.majic/output"
const FlagKeyDetailedOutput string = "detailed"
const FlagKeyVerboseOutput string = "verbose"
const DefaultDetailedOutput bool = false
const DefaultVerboseOutput bool = false

type AppOutput struct {
	// Currently using Cobra for all CLI flag processing.
	//
	// An unfortunate side-effect of this is that insight
	// into CLI flags isn't available until the command
	// stack is executed. The primary implication being
	// that all application initialization logic prior to
	// the execution of the command stack, is unaware of
	// logging-level flags specified via the command line.
	//
	// There's nothing preventing evaluating CLI args
	// directly or changing where the initialization
	// logic is initiated, but either approach would
	// need a decent amount of time to implement.
	//
	// The current solution for this is to support
	// specifying the desired logging-level via the
	// configuration file.  The co-mingled use of specifying
	// logging-level settings via both the configuration
	// file and command line flags is evaluated in an
	// additive sense, meaning that a true value specified
	// in the configuation file cannot be unset via a flag
	// and vice-versa.
	detailed bool
	verbose  bool
}

var output *AppOutput

func Output() *AppOutput {
	if output == nil {
		output = &AppOutput{DefaultDetailedOutput, DefaultVerboseOutput}
	}
	return output
}

func (output *AppOutput) NormalOutput(message string) {
	output.print(message)
}

func (output *AppOutput) DetailedOutput(message string) {
	if output.detailed || output.verbose {
		output.print(message)
	}
}

func (output *AppOutput) VerboseOutput(message string) {
	if output.verbose {
		output.print(message)
	}
}

func (output *AppOutput) print(message string) {
	if strings.HasSuffix(message, "\n") {
		fmt.Print(message)
	} else {
		fmt.Println(message)
	}
}

type AppConfig struct {
	configFileSettings *properties.Properties
}

var appConfig *AppConfig

func Config() *AppConfig {
	if appConfig == nil {
		appConfig = loadConfig()
	}
	return appConfig
}

func (config *AppConfig) Get(key string) interface{} {
	return config.GetD(key, "")
}

func (config *AppConfig) GetD(key string, defaultVal string) interface{} {
	value, found := config.configFileSettings.Get(key)
	if !found {
		Output().NormalOutput("Could not retrieve configuration value for key: " + key)
		value = defaultVal
	}
	return value
}

func (config *AppConfig) Set(key string, value interface{}) bool {
	result := true
	err := config.configFileSettings.SetValue(key, value)
	if err != nil {
		Output().NormalOutput("Could not set configuration value \"" + value.(string) + "\" for key: " + key)
		result = false
	}
	return result
}

func (config *AppConfig) ApplyFlags(flags *pflag.FlagSet) {
	detailed, err := flags.GetBool(FlagKeyDetailedOutput)
	if err != nil {
		detailed = false
	}
	err = nil
	verbose, err := flags.GetBool(FlagKeyVerboseOutput)
	if err != nil {
		verbose = false
	}
	Output().detailed = Output().detailed || detailed
	Output().verbose = Output().verbose || verbose
	config.Set(FlagKeyDetailedOutput, Output().detailed)
	config.Set(FlagKeyVerboseOutput, Output().verbose)
}

// Should this be public?
func loadConfig() *AppConfig {
	appConfig = &AppConfig{loadConfigFile()}
	detailed, err := strconv.ParseBool(appConfig.Get(FlagKeyDetailedOutput).(string))
	if err != nil {
		detailed = false
	}
	verbose, err := strconv.ParseBool(appConfig.Get(FlagKeyVerboseOutput).(string))
	if err != nil {
		verbose = false
	}
	Output().detailed = detailed
	Output().verbose = verbose
	return appConfig
}

func loadConfigFile() *properties.Properties {
	userHomeDirPath, err := os.UserHomeDir()
	HandleError(err)
	appConfigFilePath := filepath.Join(userHomeDirPath, AppConfigHome, AppConfigFile)
	Output().DetailedOutput("Configuration file: " + appConfigFilePath)
	_, err = os.Stat(appConfigFilePath)
	if err != nil && !errors.Is(err, os.ErrExist) {
		Output().DetailedOutput("Generating default properties...")
		generateDefaultConfigFile()
	}
	configFileSettings, err := properties.LoadFile(appConfigFilePath, properties.UTF8)
	HandleError(err)
	return configFileSettings
}

func generateDefaultConfigFile() {
	props := properties.NewProperties()
	props.SetValue(ConfigKeyOutputDir, DefaultOutputDir)
	props.SetValue(ConfigKeyInputDir, DefaultInputDir)
	props.SetValue(ConfigKeyPluginsDir, DefaultPluginsDir)
	userHomeDirPath, err := os.UserHomeDir()
	HandleError(err)
	appConfigFilePath := filepath.Join(userHomeDirPath, AppConfigHome)
	err = os.Mkdir(appConfigFilePath, 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		HandleError(err)
	}
	appConfigFilePath = filepath.Join(appConfigFilePath, AppConfigFile)
	appConfigFile, err := os.OpenFile(appConfigFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil && !errors.Is(err, os.ErrExist) {
		HandleError(err)
	}
	defer appConfigFile.Close()
	_, err = props.Write(appConfigFile, properties.UTF8)
	HandleError(err)
}

func HandleError(e error) {
	if e != nil {
		Output().NormalOutput("Terminating due to error")
		panic(e)
	}
}
