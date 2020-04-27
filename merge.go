package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func startMerge(usrWrNames []string, usrRdName string) {
	usrWrFilepaths := make([]string, len(usrWrNames))
	usrWrFiles := make([]*os.File, len(usrWrNames))
	for i, usrWrName := range usrWrNames {
		usrWrFilepath := fmt.Sprintf("%s.wr", usrWrName)
		usrWrFilepaths[i] = usrWrFilepath

		usrWrFile, err := createOpenFile(usrWrFilepath)
		if err != nil {
			log.Fatal("Error making/opening named pipe for writing: ", err)
		}
		usrWrFiles[i] = usrWrFile
	}

	usrWrReaders := make([]io.Reader, len(usrWrFiles))
	for i, usrWrFiles := range usrWrFiles {
		usrWrReaders[i] = usrWrFiles
	}

	usrRdFilepath := fmt.Sprintf("%s.rd", usrRdName)

	usrRdFile, err := createOpenFile(usrRdFilepath)
	if err != nil {
		log.Fatal("Error making/opening named pipe for reading: ", err)
	}

	log.Printf("Running merge from %s to %s\n", strings.Join(usrWrFilepaths, ", "), usrRdFilepath)
	runMerge(usrWrReaders, usrRdFile)
}

func runMerge(usrWrFiles []io.Reader, usrRdFile io.Writer) {
	merger := make(chan []byte)
	logging := make(chan string)
	errors := make(chan error)
	quit := make(chan error)
	quitCount := 0

	for _, usrWrFile := range usrWrFiles {
		go scanReader(usrWrFile, merger, logging, errors, quit)
	}

	for {
		select {
		case msgWithDelim := <-merger:
			_, err := usrRdFile.Write(msgWithDelim)
			if err != nil {
				log.Fatalln("Error dispatching merge message:", err)
			}
		case logMsg := <-logging:
			printLogLn(logMsg)
		case err := <-errors:
			log.Println(err)
		case err := <-quit:
			if err != nil {
				log.Fatalln("Error operating merge loop:", err)
			}
			quitCount++
			if quitCount == len(usrWrFiles) {
				return
			}
		}
	}
}

func scanReader(
	usrWrFile io.Reader,
	merger chan<- []byte,
	logging chan<- string,
	errors chan<- error,
	quit chan<- error,
) {
	scanner := NewHeavyDutyScanner(usrWrFile, MsgDelimiter)
	scanner.Decode = DecodeBase64Message

	for scanner.Scan() {
		msg, err := scanner.DecodedMessage()
		if err != nil {
			errors <- err
		}
		logMsg := fmt.Sprintf("MERGING: %s", msg)
		logging <- logMsg

		merger <- scanner.DelimitedMessage()
	}

	errors <- scanner.Err()
}
