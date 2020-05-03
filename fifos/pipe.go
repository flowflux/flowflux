package fifos

import (
	"flowflux/flowscan"
	"flowflux/printer"
	"fmt"
	"io"
	"log"
)

// StartPipe ...
func StartPipe(name string) {
	usrWrFilepath := fmt.Sprintf("%s.wr", name)
	usrRdFilepath := fmt.Sprintf("%s.rd", name)

	usrWrFile, err := createOpenFifo(usrWrFilepath)
	if err != nil {
		log.Fatal("Error making/opening named pipe for writing: ", err)
	}

	usrRdFile, err := createOpenFifo(usrRdFilepath)
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
	scanner := flowscan.NewLengthPrefix(usrWrFile)

	for scanner.Scan() {
		msg := scanner.Message()

		printer.LogLn(fmt.Sprintf("PIPING: %s", msg))

		_, err := usrRdFile.Write(scanner.PrefixedMessage())
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}
