package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func startCleanup(dirName string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory: ", err)
	}
	cleanupDir := path.Join(cwd, dirName)
	files, err := ioutil.ReadDir(cleanupDir)
	if err != nil {
		log.Fatal("Error reading directory: ", cleanupDir, ": ", err)
	}
	err = runCleanup(cleanupDir, files)
	if err != nil {
		log.Fatal("Error running cleanup on directory: ", cleanupDir, ": ", err)
	}
}

func runCleanup(basePath string, files []os.FileInfo) error {
	cleanedUpCount := 0
	for _, fi := range files {
		ext := path.Ext(fi.Name())
		if ext == ".wr" || ext == ".rd" {
			if fi.Mode()&os.ModeNamedPipe != 0 {
				filePath := path.Join(basePath, fi.Name())
				err := os.Remove(filePath)
				if err != nil {
					return err
				}
				fmt.Println("Removed:", filePath)
				cleanedUpCount++
			}
		}
	}
	if cleanedUpCount == 0 {
		fmt.Println("Nothing to cleanup")
	} else {
		fmt.Printf("%v files cleaned up\n", cleanedUpCount)
	}
	return nil
}
