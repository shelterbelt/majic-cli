# majic-cli

*Because who doesn't need a bespoke CLI to handle day-to-day tasks?*

> My primary motivation for this project is learning [Go](https://go.dev), and I'm really enjoying the experience. :tada:  There's a good chance I'll have better-informed opinions about my implementation choices as I gain more experience with the language, so conventions and patterns in the project are likely to evolve over time.

## Overview

A general-purpose command-line interface application with a modular architecture to facilitate easily adding new commands.

**majic** provides:

- robust command-line interface behavior via [Cobra](https://github.com/spf13/cobra)
- functionality to simplify processing directories of files and their contents 
- a drop-in plugin architecture based on [Go plugins](https://pkg.go.dev/plugin)

## Concepts

**majic** is built with extensibility as its primary objective. The tool is intended to provide a foundation to build small, simple actions that can be quickly deployed and combined to create robust functionality without requiring changes to **majic** itself.

There are three primary elements to **majic**'s code:

- [Helpers](#helpers)
- [Processors](#processors)
- [Plugins](#plugins)

### Helpers

Helpers in **majic** are essentially little more than [Go packages](https://go.dev/doc/effective_go#package-names) that serve as an API facade for the category of functionality provided by the Helper.  The majority of **majic**'s functionality is implemented in [Helpers](./majic/helpers).

**majic** includes a handful of embedded helper packages that provide the base capabilities of the application:

- [`core`](./majic/helpers/core/core.go): configuration, output, and general error handling
- [`file`](./majic/helpers/file/file.go): reading and modifying directories of files and file contents
- [`plugin`](./majic/helpers/plugin/plugin.go): adding functionality via plugins

### Processors

Processors are structure instances that conform to the [`Processor`](./majic/helpers/core/core.go#L18-L20) interface or specialized extensions of the [`Processor`](./majic/helpers/core/core.go#L18-L20) interface.

[*Base Processor interface*](./majic/helpers/core/core.go#L18-L20)
```
type Processor interface {
	Initialize(flags *pflag.FlagSet)
}
```

[*Example of a specialized Processor interface*](./majic/helpers/file/file.go#L19-L30)
```
type FileProcessor interface {
	Initialize(flags *pflag.FlagSet)
	ShouldProcessFile(fileName string) bool
	UseGeneratedFileNames() bool
	PreprocessNewTargetFile(file *os.File)
	TargetFileName() string
	ProcessLine(input string) string
	Reset()
}
```

These interfaces enable inserting different behavior into otherwise functionally similar operations, such as applying different formatting rules to all files in a directory of files.

### Plugins

Plugins are external modules that implement the [`MajicPlugin`](./majic/helpers/plugin/plugin.go#L17-L19) interface and are compiled into a separate binary from the **majic** executable.

```
type MajicPlugin interface {
	Register() []string
}
```

Adding functionality to **majic** is done by adding a structure that implements the [`MajicPlugin`](./majic/helpers/plugin/plugin.go#L17-L19) interface to an application and then exposing a public variable of that `struct` type named `Plugin` in the application's `main` package.

```
type myplugin struct {
}

// Implement MajicPlugin interface
func (plugin *myplugin) Register() []string {
	return []string{"Command1Cmd", "Command2Cmd", "Command3Cmd"}
}

var Plugin myplugin
```

In the `Register` function, return a slice of `string`s listing the names of the Cobra Command variables that define the commands for the plugin.

Build the application as a plugin using: 

```
go build -buildmode=plugin
``` 

and copy the resulting `*.so` binary to **majic's** plugins directory (default: `~/.majic/plugins`).

A sample plugin implementation is included in [examples](./examples)
