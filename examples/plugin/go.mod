module github.com/shelterbelt/majic-cli/examples/plugin/myplugin

go 1.24

toolchain go1.24.3

replace github.com/shelterbelt/majic-cli/majic => ../../majic

require (
	github.com/shelterbelt/majic-cli/majic v1.0.1
	github.com/spf13/cobra v1.9.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/magiconair/properties v1.8.10 // indirect
	github.com/spf13/pflag v1.0.6
)
