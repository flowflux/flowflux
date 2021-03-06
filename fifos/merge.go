package fifos

import (
	"flowflux/flowscan"
	"flowflux/printer"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// StartMerge ...
func StartMerge(usrWrNames []string, usrRdName string) {
	usrWrFilepaths := make([]string, len(usrWrNames))
	usrWrFiles := make([]*os.File, len(usrWrNames))
	for i, usrWrName := range usrWrNames {
		usrWrFilepath := fmt.Sprintf("%s.wr", usrWrName)
		usrWrFilepaths[i] = usrWrFilepath

		usrWrFile, err := createOpenFifo(usrWrFilepath)
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

	usrRdFile, err := createOpenFifo(usrRdFilepath)
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
			printer.LogLn(logMsg)
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
	scanner := flowscan.NewLengthPrefix(usrWrFile)

	for scanner.Scan() {
		msg := scanner.Message() // TODO: Test if copying is required
		logMsg := fmt.Sprintf("MERGING: %s", msg)
		logging <- logMsg

		merger <- scanner.PrefixedMessage()
	}

	errors <- scanner.Err()
}
