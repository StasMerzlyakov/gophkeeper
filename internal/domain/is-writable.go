//go:build !windows
// +build !windows

package domain

import (
	"fmt"
	"os"
	"syscall"
)

func IsWritable(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {

		return false, fmt.Errorf("%w - path doesn't exist", err)
	}

	err = nil
	if !info.IsDir() {
		return false, fmt.Errorf("path isn't a directory")
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return false, fmt.Errorf("write permission bit is not set on this file for user")
	}

	var stat syscall.Stat_t
	if err = syscall.Stat(path, &stat); err != nil {
		return false, fmt.Errorf("unable to get stat")
	}

	err = nil
	if uint32(os.Geteuid()) != stat.Uid {
		return false, fmt.Errorf("user doesn't have permission to write to this directory")
	}

	return true, nil
}
