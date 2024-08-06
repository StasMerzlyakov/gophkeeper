package domain

import (
	"fmt"
	"os"
)

func IsWritable(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("%w - path doesn't exists", err)
	}

	err = nil
	if !info.IsDir() {
		return false, fmt.Errorf("%w - path is't a directory", err)
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return false, fmt.Errorf("write permission bit is not set on this file for user")
	}
	return true, nil
}
