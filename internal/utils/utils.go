package utils

import (
	"path/filepath"
	"runtime"
	"strconv"
)

// Returns source filename and name of the caller function. Used for logging purposes
func CallerFilename() string {
	_, filename, line, _ := runtime.Caller(1)
	_, name := filepath.Split(filename)
	return "[" + name + ":" + strconv.Itoa(line) + "]"
}

// Are you fucking kidding me
func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
