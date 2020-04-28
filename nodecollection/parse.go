package nodecollection

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

func parseHubFile(filePath string) map[string]Node {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading hub file: ", filePath, ": ", err)
	}

	connections := parseConnections(fileContent)
	mainIndex := make(map[string]Node)
	connectionIndex := make(map[string]Node)

	// 1st analyze connection

	for _, forkConn := range findForkConnections(connections) {
		n := Node{
			Class:   ForkClass,
			ID:      makeForkID(forkConn.fromCmd, forkConn.toCmds),
			OutKeys: mapToCmdIDs(forkConn.toCmds...),
		}
		mainIndex[n.ID] = n
		connectionIndex[forkFromKey(forkConn.fromCmd)] = n
		for _, toCmd := range forkConn.toCmds {
			connectionIndex[forkToKey(toCmd)] = n
		}
	}

	for _, mergeConn := range findMergeConnections(connections) {
		n := Node{
			Class:   MergeClass,
			ID:      makeMergeID(mergeConn.fromCmds, mergeConn.toCmd),
			OutKeys: mapToCmdIDs(mergeConn.toCmd),
		}
		mainIndex[n.ID] = n
		connectionIndex[mergeToKey(mergeConn.toCmd)] = n
		for _, fromCmd := range mergeConn.fromCmds {
			connectionIndex[mergeFromKey(fromCmd)] = n
		}
	}

	for _, pipeConn := range findPipeConnections(connections) {
		n := Node{
			Class:   PipeClass,
			ID:      makePipeID(pipeConn.fromCmd, pipeConn.toCmd),
			OutKeys: mapToCmdIDs(pipeConn.toCmd),
		}
		mainIndex[n.ID] = n
		connectionIndex[pipeFromKey(pipeConn.fromCmd)] = n
		connectionIndex[pipeToKey(pipeConn.toCmd)] = n
	}

	for _, inOutConn := range findInputOutputConnection(connections) {
		if len(inOutConn.inputTo) > 0 {
			in := Node{
				Class:      InputClass,
				ID:         makeInputID(inOutConn.inputTo),
				ScanMethod: inOutConn.inputScanMethod,
				OutKeys:    mapToCmdIDs(inOutConn.inputTo),
			}
			mainIndex[in.ID] = in
			connectionIndex[inputToKey(inOutConn.inputTo)] = in
		}
		if len(inOutConn.outputFrom) > 0 {
			out := Node{
				Class: OutputClass,
				ID:    makeOutputID(inOutConn.outputFrom),
			}
			mainIndex[out.ID] = out
			connectionIndex[outputFromKey(inOutConn.outputFrom)] = out
		}
	}

	// 2nd setup process nodes

	for _, uniProc := range findUniqueProcesses(connections) {
		n := Node{
			Class:      ProcessClass,
			ID:         makeCmdID(uniProc.name),
			ScanMethod: uniProc.scanMethod,
			Process:    parseProcessCommand(uniProc.name),
		}

		next, ok := connectionIndex[pipeFromKey(uniProc.name)]
		if ok {
			n.OutKeys = append(n.OutKeys, next.ID)
		}

		next, ok = connectionIndex[forkFromKey(uniProc.name)]
		if ok {
			n.OutKeys = append(n.OutKeys, next.ID)
		}

		next, ok = connectionIndex[mergeFromKey(uniProc.name)]
		if ok {
			n.OutKeys = append(n.OutKeys, next.ID)
		}

		next, ok = connectionIndex[outputFromKey(uniProc.name)]
		if ok {
			n.OutKeys = append(n.OutKeys, next.ID)
		}

		mainIndex[n.ID] = n
	}

	return mainIndex
}

func parseProcessCommand(command string) ProcessCommand {
	rawSplit := strings.Split(command, " ")
	comps := make([]string, 0)
	for _, comp := range rawSplit {
		trimmedComp := strings.TrimSpace(comp)
		if len(trimmedComp) > 0 {
			comps = append(comps, trimmedComp)
		}
	}
	return ProcessCommand{
		Command:   comps[0],
		Arguments: comps[1:],
	}
}

func parseConnections(fileContent []byte) []connection {
	connections := make([]connection, 0)
	contained := make(map[string]bool)
	addTriple := func(fromCmd string, scanMethodStr string, toCmd string) {
		key := fmt.Sprintf("%s%s%s", fromCmd, scanMethodStr, toCmd)
		_, ok := contained[key]
		if !ok {
			var scanMethod ScanMethod
			if scanMethodStr == "->" {
				scanMethod = ScanMessages
			} else if scanMethodStr == "*->" {
				scanMethod = ScanRawBytes
			}
			conn := connection{
				fromCmd:        fromCmd,
				fromScanMethod: scanMethod,
				toCmd:          toCmd,
			}
			connections = append(connections, conn)
			contained[key] = true
		}
	}
	for _, line := range bytes.Split(fileContent, []byte{'\n'}) {
		tokens := tokenizeLine(line)
		iterMax := len(tokens) - 2
		for i := 0; i < iterMax; i += 2 {
			addTriple(
				string(tokens[i]),
				string(tokens[i+1]),
				string(tokens[i+2]),
			)
		}
	}
	return connections
}

func tokenizeLine(line []byte) [][]byte {
	tokens := make([][]byte, 0)
	re := regexp.MustCompile(`\*\->|\->`)
	indexes := re.FindAllIndex(line, -1)
	if indexes == nil {
		return tokens
	}
	lastEnd := 0
	for _, idx := range indexes {
		start, end := idx[0], idx[1]
		fromCmd := bytes.TrimSpace(line[lastEnd:start])
		finding := line[start:end]
		lastEnd = end
		tokens = append(tokens, fromCmd, finding)
	}
	lastCmd := bytes.TrimSpace(line[lastEnd:])
	tokens = append(tokens, lastCmd)
	return tokens
}
