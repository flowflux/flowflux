package nodecollection

import (
	"fmt"
	"strings"
)

// IDs

func makeInputID(toCmd string) string {
	return inputToKey(toCmd)
}
func makeOutputID(fromCmd string) string {
	return outputFromKey(fromCmd)
}

func makeForkID(fromCmd string, toCmds []string) string {
	return fmt.Sprintf(
		"fork:%s->%s",
		underscoreWhitespace(fromCmd),
		underscoreWhitespace(join(toCmds, ",")),
	)
}

func makeMergeID(fromCmds []string, toCmd string) string {
	return fmt.Sprintf(
		"merge:%s->%s",
		underscoreWhitespace(join(fromCmds, ",")),
		underscoreWhitespace(toCmd),
	)
}

func makePipeID(fromCmd string, toCmd string) string {
	return fmt.Sprintf(
		"pipe:%s->%s",
		underscoreWhitespace(fromCmd),
		underscoreWhitespace(toCmd),
	)
}

func makeCmdID(cmd string) string {
	return fmt.Sprintf(
		"cmd:%s",
		underscoreWhitespace(cmd),
	)
}

func join(elements []string, sep string) string {
	return strings.Join(elements, sep)
}
func underscoreWhitespace(value string) string {
	return strings.ReplaceAll(value, " ", "_")
}

func mapToCmdIDs(cmds ...string) []string {
	cmdIDs := make([]string, len(cmds))
	for i, c := range cmds {
		cmdIDs[i] = makeCmdID(c)
	}
	return cmdIDs
}

// Keys

func pipeFromKey(fromCmd string) string {
	return makeFromKey("pipe", fromCmd)
}
func pipeToKey(toCmd string) string {
	return makeToKey("pipe", toCmd)
}

func forkFromKey(fromCmd string) string {
	return makeFromKey("fork", fromCmd)
}
func forkToKey(toCmd string) string {
	return makeToKey("fork", toCmd)
}

func mergeFromKey(fromCmd string) string {
	return makeFromKey("merge", fromCmd)
}
func mergeToKey(toCmd string) string {
	return makeToKey("merge", toCmd)
}

func outputFromKey(fromCmd string) string {
	return makeFromKey("output", fromCmd)
}
func inputToKey(toCmd string) string {
	return makeToKey("input", toCmd)
}

func makeFromKey(kind string, fromCmd string) string {
	return fmt.Sprintf("%s:%s->", kind, fromCmd)
}
func makeToKey(kind string, toCmd string) string {
	return fmt.Sprintf("%s:->%s", kind, toCmd)
}
