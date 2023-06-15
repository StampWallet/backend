package utils

import (
	"math/rand"
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

func RandomSlice[T any](size int, values []T) []T {
	result := make([]T, size)
	for i := range result {
		result[i] = values[rand.Intn(len(values))]
	}
	return result
}

func Map[U any, V any](us []U, fc func(U) V) []V {
	vs := make([]V, len(us))
	for i := range vs {
		vs[i] = fc(us[i])
	}
	return vs
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
