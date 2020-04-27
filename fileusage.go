package main

import (
	"os"
	"syscall"
)

func isUsable(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	if info.IsDir() {
		return false
	}
	return info.Mode()&os.ModeNamedPipe != 0
}

func createOpenFile(filePath string) (*os.File, error) {
	if !isUsable(filePath) {
		os.Remove(filePath)
		err := syscall.Mkfifo(filePath, 0666)
		if err != nil {
			return nil, err
		}
	}
	return os.OpenFile(filePath, os.O_RDWR, os.ModeNamedPipe)
}
