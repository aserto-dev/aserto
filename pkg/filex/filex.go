package filex

import "os"

// DirExist verifies if directory path exists and is a directory.
func DirExists(dirPath string) bool {
	if fi, err := os.Stat(dirPath); err == nil && fi.IsDir() {
		return true
	}
	return false
}

// FileExist verifies if file path exists and is a file (and not a directory).
func FileExists(filePath string) bool {
	fi, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !fi.IsDir()
}
