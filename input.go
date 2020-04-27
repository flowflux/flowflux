package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func startInput(usrRdFilepath string) {
	usrRdFile, err := createOpenFile(usrRdFilepath)
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
		text, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		text = strings.Replace(text, "\n", "", -1)

		msg := []byte(text)
		msgB64 := make([]byte, base64.StdEncoding.EncodedLen(len(msg)))

		base64.StdEncoding.Encode(msgB64, msg)
		msgB64 = append(msgB64, MsgDelimiter...)
		_, err = usrRdFile.Write(msgB64)
		if err != nil {
			return err
		}
	}
}
