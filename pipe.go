package main

import (
	"fmt"
	"io"
	"log"
)

func startPipe(name string) {
	usrWrFilepath := fmt.Sprintf("%s.wr", name)
	usrRdFilepath := fmt.Sprintf("%s.rd", name)

	usrWrFile, err := createOpenFile(usrWrFilepath)
	if err != nil {
		log.Fatal("Error making/opening named pipe for writing: ", err)
	}

	usrRdFile, err := createOpenFile(usrRdFilepath)
	if err != nil {
		log.Fatal("Error making/opening named pipe for reading: ", err)
	}

	log.Printf("Running pipe from %s to %s\n", usrWrFilepath, usrRdFilepath)
	err = runPipe(usrWrFile, usrRdFile)
	if err != nil {
		log.Fatalln("Error operating pipe loop:", err)
	}
}

func runPipe(usrWrFile io.Reader, usrRdFile io.Writer) error {
	scanner := NewHeavyDutyScanner(usrWrFile, MsgDelimiter)
	scanner.Decode = DecodeBase64Message

	for scanner.Scan() {
		msg, err := scanner.DecodedMessage()
		if err != nil {
			return err
		}

		printLogLn(fmt.Sprintf("PIPING: %s", msg))

		_, err = usrRdFile.Write(scanner.DelimitedMessage())
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}
