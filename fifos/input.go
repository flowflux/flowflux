package fifos

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// StartInput ...
func StartInput(usrRdFilepath string) {
	usrRdFile, err := createOpenFifo(usrRdFilepath)
	if err != nil {
		log.Fatal("Error making/opening named pipe for reading: ", err)
	}

	log.Printf("Dispatching input to %s\n", usrRdFilepath)
	err = runInput(os.Stdin, usrRdFile)
	if err != nil {
		log.Fatalln("Error operating pipe loop:", err)
	}
}

func runInput(usrWrFile io.Reader, usrRdFile io.Writer) error {
	reader := bufio.NewReader(usrWrFile)
	fmt.Println("Newline concludes message")

	for {
		fmt.Print("JSON -> ")
		msgStr, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		msgStr = strings.Replace(msgStr, "\n", "", -1)

		msgBytes := []byte(msgStr)
		msgLen := len(msgBytes)
		msgLenStr := fmt.Sprintf("%020d", msgLen)
		envelope := append([]byte(msgLenStr), msgBytes...)

		_, err = usrRdFile.Write(envelope)
		if err != nil {
			return err
		}
	}
}
