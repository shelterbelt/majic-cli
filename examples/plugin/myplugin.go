/*
Copyright © 2023 Mark Johnson
*/
package main

import (
	"github.com/shelterbelt/majic-cli/majic/helpers/core"
)

type myplugin struct {
}

func (plugin *myplugin) Register() []string {
	core.Output().DetailedOutput("Registering majic CLI Sample plugin.")
	return []string{"SayHiCmd", "FilesCmd"}
}

var Plugin myplugin
