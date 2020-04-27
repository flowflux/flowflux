package main

import "os"

func printLogLn(text string) {
	if len(text) > 78 {
		text = text[:78]
	}
	os.Stderr.WriteString(text + "\n")
}

func printErrLn(text string) {
	os.Stderr.WriteString(text + "\n")
}
