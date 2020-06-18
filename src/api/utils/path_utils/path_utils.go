package path_utils

import (
	"os"
	"path/filepath"
)

func EnforceDirectory(path string) string {

	sep := `/`
	if path == "" {
		return sep
	} else if path[len(path) - 1:] == sep {
		return path
	}
	return path + sep
}

func GetCurrentDir() string {
	ex, _ := os.Executable()
	return filepath.Dir(ex)
}