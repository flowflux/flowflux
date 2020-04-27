package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	argsLen := len(args)

	if argsLen == 0 {
		printHeader()
		printHelp()
		return
	} else if argsLen == 1 {
		coll := NewNodeCollection(args[0])
		err := validateNodeCollection(coll)
		if err != nil {
			log.Fatal(err)
		}
		runNodeCollection(coll)
		return
	}

	switch args[0] {
	case "pipe":
		if argsLen < 2 {
			fmt.Println("Not enough arguments for pipe:")
			printPipeHelp()
		} else {
			startPipe(args[1])
		}
	case "fork":
		if argsLen < 4 {
			fmt.Println("Not enough arguments for fork:")
			printForkHelp()
		} else {
			startFork(args[1], args[2:])
		}
	case "merge":
		if argsLen < 4 {
			fmt.Println("Not enough arguments for merge:")
			printMergeHelp()
		} else {
			lastIdx := len(args) - 1
			startMerge(args[1:lastIdx], args[lastIdx])
		}
	case "input":
		if argsLen < 2 {
			fmt.Println("Not enough arguments for input:")
			printInputHelp()
		} else {
			startInput(args[1])
		}
	case "cleanup":
		if argsLen < 2 {
			startCleanup(".")
		} else {
			startCleanup(args[1])
		}
	default:
		printHeader()
		printHelp()
	}
}

func printHeader() {
	fmt.Println("approx/hub -- utility to build actor-based systems by composing command line processes")
	fmt.Println("Usage with <hub.flow> file:")
	fmt.Println("  ./hub <path-to-file>")
}

func printHelp() {
	fmt.Println("Usage for named pipe management:")
	printPipeHelp()
	printForkHelp()
	printMergeHelp()
	printInputHelp()
	printCleanupHelp()
}
func printPipeHelp() {
	fmt.Println("  pipe <name>")
	fmt.Println("    Pipe message stream from <name>.wr to <name>.rd")
}
func printForkHelp() {
	fmt.Println("  fork <wr-name> <rd-name-1> <rd-name-2> <...>")
	fmt.Println("    Fork message stream from wr-fifo into all provided rd-fifos")
}
func printMergeHelp() {
	fmt.Println("  merge <wr-name-1> <wr-name-2> <...< <rd-name>")
	fmt.Println("    Merge message stream from all provided wr-fifos into rd-fifo")
}
func printInputHelp() {
	fmt.Println("  input <name>")
	fmt.Println("    Input JSON messages to stream them to <name>")
}
func printCleanupHelp() {
	fmt.Println("  cleanup <directory>")
	fmt.Println("    Cleanup directory from fifos (wr & rd)")
}
