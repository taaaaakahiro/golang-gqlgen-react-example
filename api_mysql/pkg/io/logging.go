package io

import (
	"io"
	"log"
	"os"
)

// opens a file for logging. Will return a fatal error if the file cannot be opened.
func openLogFile(path string) *os.File {
	var flag int

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		flag = os.O_CREATE | os.O_WRONLY
	} else {
		flag = os.O_APPEND | os.O_WRONLY
	}

	f, err := os.OpenFile(path, flag, 0600)

	if err != nil {
		log.Fatal(err)
	}

	return f
}

// CreateLogWriters creates appropriate writers for logs based on the logging section of the application configuration.
// Will log a fatal error if the logging output is invalid, or log files cannot be opened.
func CreateLogWriters(output, accessFile, errorFile string) (access io.Writer, error io.Writer) {
	switch output {
	case "std":
		access = io.MultiWriter(os.Stdout)
		error = io.MultiWriter(os.Stdout)
	case "file":
		access = io.MultiWriter(openLogFile(accessFile))
		error = io.MultiWriter(openLogFile(errorFile))
	case "both":
		access = io.MultiWriter(openLogFile(accessFile), os.Stdout)
		error = io.MultiWriter(openLogFile(errorFile), os.Stdout)
	default:
		log.Fatal("logging must be std, file, or both")
	}
	return
}
