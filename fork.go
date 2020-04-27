package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func startFork(usrWrName string, usrRdNames []string) {
	usrWrFilepath := fmt.Sprintf("%s.wr", usrWrName)

	usrWrFile, err := createOpenFile(usrWrFilepath)
	if err != nil {
		log.Fatal("Error making/opening named pipe for writing: ", err)
	}

	usrRdFilepaths := make([]string, len(usrRdNames))
	usrRdFiles := make([]*os.File, len(usrRdNames))
	for i, usrRdName := range usrRdNames {
		usrRdFilepath := fmt.Sprintf("%s.rd", usrRdName)
		usrRdFilepaths[i] = usrRdFilepath

		usrRdFile, err := createOpenFile(usrRdFilepath)
		if err != nil {
			log.Fatal("Error making/opening named pipe for reading: ", err)
		}
		usrRdFiles[i] = usrRdFile
	}

	usrRdWriters := make([]io.Writer, len(usrRdFiles))
	for i, usrRdFiles := range usrRdFiles {
		usrRdWriters[i] = usrRdFiles
	}

	log.Printf("Running fork from %s to %s\n", usrWrFilepath, strings.Join(usrRdFilepaths, ", "))
	err = runFork(usrWrFile, usrRdWriters)
	if err != nil {
		log.Fatalln("Error operating fork loop:", err)
	}
}

func runFork(usrWrFile io.Reader, usrRdFiles []io.Writer) error {
	scanner := NewHeavyDutyScanner(usrWrFile, MsgDelimiter)
	scanner.Decode = DecodeBase64Message

	for scanner.Scan() {
		msg, err := scanner.DecodedMessage()
		if err != nil {
			return err
		}
		printLogLn(fmt.Sprintf("FORKING: %s", msg))

		for _, usrRdFile := range usrRdFiles {
			_, err = usrRdFile.Write(scanner.DelimitedMessage())
			if err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}
