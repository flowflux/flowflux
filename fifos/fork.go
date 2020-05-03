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

// StartFork ...
func StartFork(usrWrName string, usrRdNames []string) {
	usrWrFilepath := fmt.Sprintf("%s.wr", usrWrName)

	usrWrFile, err := createOpenFifo(usrWrFilepath)
	if err != nil {
		log.Fatal("Error making/opening named pipe for writing: ", err)
	}

	usrRdFilepaths := make([]string, len(usrRdNames))
	usrRdFiles := make([]*os.File, len(usrRdNames))
	for i, usrRdName := range usrRdNames {
		usrRdFilepath := fmt.Sprintf("%s.rd", usrRdName)
		usrRdFilepaths[i] = usrRdFilepath

		usrRdFile, err := createOpenFifo(usrRdFilepath)
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
	scanner := flowscan.NewLengthPrefix(usrWrFile)

	for scanner.Scan() {
		msg := scanner.Message()
		printer.LogLn(fmt.Sprintf("FORKING: %s", msg))

		for _, usrRdFile := range usrRdFiles {
			_, err := usrRdFile.Write(scanner.PrefixedMessage())
			if err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}
