package utils

import (
	"path/filepath"
	"runtime"
	"strconv"
)

func CallerFilename() string {
	_, filename, line, _ := runtime.Caller(1)
	_, name := filepath.Split(filename)
	return "[" + name + ":" + strconv.Itoa(line) + "]"
}
