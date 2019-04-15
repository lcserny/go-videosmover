package goutils

import (
	"bufio"
	"os"
)

func ReadFileToLines(file string) []string {
	openFile, err := os.Open(file)
	LogFatal(err)
	defer CloseFile(openFile)

	var lines []string
	scanner := bufio.NewScanner(openFile)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func CloseFile(file *os.File) {
	err := file.Close()
	LogFatal(err)
}
