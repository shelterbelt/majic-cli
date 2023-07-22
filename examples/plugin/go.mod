module github.com/shelterbelt/majic-cli/examples/plugin/myplugin

go 1.20

replace github.com/shelterbelt/majic-cli/majic => ../../majic

require (
	github.com/shelterbelt/majic-cli/majic v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.7.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/spf13/pflag v1.0.5
)
