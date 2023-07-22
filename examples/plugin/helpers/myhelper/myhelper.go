/*
Copyright © 2023 Mark Johnson
*/
package myhelper

import (
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type MyFileProcessor struct {
	outputFileName string
}

func (processor *MyFileProcessor) Initialize(flags *pflag.FlagSet) {
}

func (process *MyFileProcessor) PreprocessNewTargetFile(file *os.File) {
}

func (processor *MyFileProcessor) ProcessLine(fileLine string) string {
	return doTheThing(fileLine)
}

func (processor *MyFileProcessor) ShouldProcessFile(fileName string) bool {
	return strings.HasSuffix(fileName, ".md")
}

func (processor *MyFileProcessor) UseGeneratedFileNames() bool {
	return false
}

func (processor *MyFileProcessor) TargetFileName() string {
	return processor.outputFileName
}

func (processor *MyFileProcessor) Reset() {
	processor.outputFileName = ""
}

func doTheThing(contentLine string) string {
	result := strings.ReplaceAll(contentLine, "the", "THE")
	return result + "\n"
}
