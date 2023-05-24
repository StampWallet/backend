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
