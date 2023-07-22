/*
Copyright © 2023 Mark Johnson
*/
package file

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/shelterbelt/majic-cli/majic/helpers/core"
	"github.com/spf13/pflag"
)

type FileProcessor interface {
	// The processing cycle defined here feels a little
	// more ad-hoc than ideal, but is largely meeting
	// the need so far.
	Initialize(flags *pflag.FlagSet)
	ShouldProcessFile(fileName string) bool
	UseGeneratedFileNames() bool
	PreprocessNewTargetFile(file *os.File)
	TargetFileName() string
	ProcessLine(input string) string
	Reset()
}

func ProcessPath(inputIdentifier string, processor FileProcessor) {
	outputDirPath := core.Config().GetD(core.ConfigKeyOutputDir, core.DefaultOutputDir).(string)
	info, err := os.Stat(inputIdentifier)
	if err != nil && !os.IsExist(err) {
		core.Output().NormalOutput("Item " + inputIdentifier + " does not exist.")
	} else {
		core.HandleError(err)
		if info.IsDir() {
			ProcessDirectory(inputIdentifier, outputDirPath, processor)
		} else {
			ProcessFile(inputIdentifier, outputDirPath, processor)
		}
	}
}

func ProcessDirectory(sourceDirPath string, outputDirPath string, processor FileProcessor) {
	sourceDir, err := os.Open(sourceDirPath)
	core.HandleError(err)
	defer sourceDir.Close()
	files, err := sourceDir.Readdirnames(-1)
	core.HandleError(err)
	for i := 0; i < len(files); i++ {
		core.Output().DetailedOutput(files[i])
		sourceFilePath := filepath.Join(sourceDirPath, files[i])
		info, err := os.Stat(sourceFilePath)
		if err != nil && !errors.Is(err, os.ErrExist) {
			core.Output().NormalOutput("File " + sourceFilePath + " does not exist.")
		} else {
			core.HandleError(err)
			if info.IsDir() {
				ProcessDirectory(sourceFilePath, outputDirPath, processor)
			} else {
				if processor.ShouldProcessFile(files[i]) {
					ProcessFile(sourceFilePath, outputDirPath, processor)
				} else {
					core.Output().NormalOutput("Skipping: " + files[i])
				}
			}
		}
	}
}

func ProcessFile(filePath string, outputDirPath string, processor FileProcessor) {
	file, err := os.Open(filePath)
	core.HandleError(err)
	defer file.Close()
	contentsScanner := bufio.NewScanner(file)
	var targetFile *os.File

	targetFileName := filepath.Base(filePath)
	if processor.UseGeneratedFileNames() {
		targetFileName = processor.TargetFileName()
	}

	processedLine := ""
	processor.Reset()
	for contentsScanner.Scan() {
		processedLine = contentsScanner.Text()

		processedLine = processor.ProcessLine(processedLine)

		if targetFile == nil {

			if len(targetFileName) > 0 {
				targetFile, err = CreateTargetFile(filepath.Join(outputDirPath, targetFileName))
				core.HandleError(err)
				processor.PreprocessNewTargetFile(targetFile)
				defer targetFile.Close()
			}
		}

		core.Output().VerboseOutput(processedLine)
		targetFile.WriteString(processedLine)
	}
}

func CreateTargetFile(targetFilePath string) (*os.File, error) {
	info, _ := os.Stat(targetFilePath)
	desiredFilePath := targetFilePath
	for i := 0; info != nil; i++ {
		targetFilePath = convertStringToUniqueFileName(desiredFilePath, i)
		info, _ = os.Stat(targetFilePath)
	}

	file, err := os.Create(targetFilePath)
	return file, err
}

func CreateTargetCopyOfInputFile(outputDirPath string, sourceFilePath string) (*os.File, error) {
	outputFileName := filepath.Base(sourceFilePath)
	outputFilePath := filepath.Join(outputDirPath, outputFileName)

	var file *os.File
	err := copyFile(sourceFilePath, outputFilePath, 2048)

	if err != nil {
		core.Output().DetailedOutput("Source file not copied. Creating from scratch.")
		file, err = CreateTargetFile(outputFilePath)
	} else {
		core.Output().DetailedOutput("Opening output copy of source page: " + outputFileName)
		file, err = os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	return file, err
}

func ConvertStringToFileName(name string) string {
	return convertStringToUniqueFileName(name, 0)
}

func convertStringToUniqueFileName(name string, index int) string {
	// TODO:  Need to add proper sanitization for file names here.
	if index > 0 {
		extension := filepath.Ext(name)
		name = strings.TrimSuffix(name, extension)
		name = fmt.Sprintf("%s %d", name, index)
		name = name + extension
	}
	return name
}

func CreateOutputDir() string {
	outputPath := core.Config().GetD(core.ConfigKeyOutputDir, core.DefaultOutputDir).(string)
	err := os.Mkdir(outputPath, 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		core.HandleError(err)
	}
	return outputPath
}

// From: https://github.com/mactsouk/opensource.com/blob/master/cp3.go
func copyFile(src, dst string, BUFFERSIZE int64) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("file %s already exists", dst)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if err != nil {
		panic(err)
	}

	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}
