// which is a command line tool to find the path of a command
package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		writeError("Missing argument.")
	} else if len(os.Args) > 2 {
		writeError("Too many arguments.")
	}

	found := false

	arg := os.Args[1]

	pathEnv := os.Getenv("PATH")
	pathExt := os.Getenv("PATHEXT")

	if pathEnv == "" {
		writeError("Unable to find PATH environment variable")
	}

	if pathExt == "" {
		writeError("Unable to find PATHEXT environment variable")
	}

	paths := splitEnv(pathEnv)
	exts := splitEnv(pathExt)

	for _, dirPath := range paths {
		if _, err := os.Stat(dirPath); errors.Is(err, fs.ErrNotExist) {
			continue
		}

		for _, ext := range exts {
			fullPath := filepath.Join(dirPath, arg+strings.ToLower(ext))

			if _, err := os.Stat(fullPath); errors.Is(err, fs.ErrNotExist) {
				continue
			}

			if found {
				fullPath = "\n" + fullPath
			}

			os.Stdout.WriteString(fullPath)
			found = true
		}
	}

	exitCode := 1
	if found {
		exitCode = 0
	}

	os.Exit(exitCode)
}

func writeError(writeErr string) {
	os.Stderr.WriteString(writeErr)
	os.Exit(1)
}

func splitEnv(env string) []string {
	return strings.Split(env, ";")
}
