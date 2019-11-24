package main

import (
	"github.com/jaqmol/approx/processor"
)

func main() {
	form, err := processor.NewFormation()
	if err != nil {
		panic(err)
	}
	form.Start()
	form.WaitForCommands()
}
