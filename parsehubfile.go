package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

// NodeClass ...
type NodeClass int

// NodeClass ...
const (
	UnknownClass NodeClass = iota
	InputClass
	ProcessClass
	ForkClass
	MergeClass
	PipeClass
	OutputClass
)

// NodeClassToString ...
func NodeClassToString(class NodeClass) string {
	switch class {
	case InputClass:
		return "input"
	case ProcessClass:
		return "process"
	case ForkClass:
		return "fork"
	case MergeClass:
		return "merge"
	case PipeClass:
		return "pipe"
	case OutputClass:
		return "output"
	default:
		return "unknown"
	}
}

// Node ...
type Node struct {
	Class      NodeClass
	ID         string
	ScanMethod ScanMethod
	Process    ProcessCommand
	OutKeys    []string
}

// ProcessCommand ...
type ProcessCommand struct {
	Command   string
	Arguments []string
}

// NodeCollection ...
type NodeCollection struct {
	index     map[string]Node
	indexKeys []string
}

// ScanMethod ...
type ScanMethod int

// ScanMethods ...
const (
	ScanMessages ScanMethod = iota
	ScanRawBytes
)

type connection struct {
	fromCmd        string
	fromScanMethod ScanMethod
	toCmd          string
}

// NewNodeCollection ...
func NewNodeCollection(filePath string) NodeCollection {
	index := parseHubFile(filePath)
	return NodeCollection{
		index: index,
	}
}

// IDs ...
func (c NodeCollection) IDs() []string {
	if c.indexKeys == nil {
		c.indexKeys = make([]string, len(c.index))
		i := 0
		for k := range c.index {
			c.indexKeys[i] = k
			i++
		}
	}
	return c.indexKeys
}

// Node ...
func (c NodeCollection) Node(id string) (Node, bool) {
	n, ok := c.index[id]
	return n, ok
}

// Outputs ...
func (c NodeCollection) Outputs(n Node) []Node {
	nodes := make([]Node, len(n.OutKeys))
	for i, key := range n.OutKeys {
		nodes[i] = c.index[key]
	}
	return nodes
}

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

func mapToCmdIDs(cmds ...string) []string {
	cmdIDs := make([]string, len(cmds))
	for i, c := range cmds {
		cmdIDs[i] = makeCmdID(c)
	}
	return cmdIDs
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

type uniqueProcess struct {
	name       string
	scanMethod ScanMethod
}

func findUniqueProcesses(connections []connection) []uniqueProcess {
	processes := make([]uniqueProcess, 0)
	contained := make(map[string]bool)
	addProcess := func(cmd string, hasScanMethod bool, scanMethod ScanMethod) {
		_, ok := contained[cmd]
		if !ok || hasScanMethod {
			proc := uniqueProcess{
				name:       cmd,
				scanMethod: scanMethod,
			}
			processes = append(processes, proc)
			contained[cmd] = true
		}
	}
	for _, conn := range connections {
		if len(conn.fromCmd) > 0 {
			addProcess(conn.fromCmd, true, conn.fromScanMethod)
		}
		if len(conn.toCmd) > 0 {
			addProcess(conn.toCmd, false, 0)
		}
	}
	return processes
}

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

func connectionsFrom(cmd string, connections []connection) []string {
	toCmds := make([]string, 0)
	for _, conn := range connections {
		if conn.fromCmd == cmd {
			if len(conn.toCmd) > 0 { // TODO: PROBABLY UNNECESSARY
				toCmds = append(toCmds, conn.toCmd)
			}
		}
	}
	return toCmds
}
func connectionsTo(cmd string, connections []connection) ([]string, []ScanMethod) {
	fromCmds := make([]string, 0)
	fromScanMethods := make([]ScanMethod, 0)
	for _, conn := range connections {
		if conn.toCmd == cmd {
			if len(conn.fromCmd) > 0 { // TODO: PROBABLY UNNECESSARY
				fromCmds = append(fromCmds, conn.fromCmd)
				fromScanMethods = append(fromScanMethods, conn.fromScanMethod)
			}
		}
	}
	return fromCmds, fromScanMethods
}

type forkConnection struct {
	fromCmd string
	toCmds  []string
}

func findForkConnections(connections []connection) []forkConnection {
	doneForFromCmd := make(map[string]bool)
	forkConnections := make([]forkConnection, 0)
	for _, conn := range connections {
		_, done := doneForFromCmd[conn.fromCmd]
		if !done {
			toCmds := connectionsFrom(conn.fromCmd, connections)
			if len(toCmds) > 1 {
				forkConn := forkConnection{
					fromCmd: conn.fromCmd,
					toCmds:  toCmds,
				}
				forkConnections = append(forkConnections, forkConn)
				doneForFromCmd[conn.fromCmd] = true
			}
			// TODO: TEST
			// else {
			// 	doneForFromCmd[fromCmd] = true
			// }
		}
	}
	return forkConnections
}

type mergeConnection struct {
	fromCmds []string
	toCmd    string
}

func findMergeConnections(connections []connection) []mergeConnection {
	doneForToCmd := make(map[string]bool)
	mergeConnections := make([]mergeConnection, 0)
	for _, conn := range connections {
		_, done := doneForToCmd[conn.toCmd]
		if !done {
			fromCmds, _ := connectionsTo(conn.toCmd, connections)
			if len(fromCmds) > 1 {
				mergeConn := mergeConnection{
					fromCmds: fromCmds,
					toCmd:    conn.toCmd,
				}
				mergeConnections = append(mergeConnections, mergeConn)
				doneForToCmd[conn.toCmd] = true
			}
			// TODO: TEST
			// else {
			// 	doneForToCmd[toCmd] = true
			// }
		}
	}
	return mergeConnections
}

type pipeConnection struct {
	fromCmd string
	toCmd   string
}

func findPipeConnections(connections []connection) []pipeConnection {
	doneForFromCmd := make(map[string]bool)
	pipeConnections := make([]pipeConnection, 0)
	for _, conn := range connections {
		_, done := doneForFromCmd[conn.fromCmd]
		if !done {
			toCmds := connectionsFrom(conn.fromCmd, connections)
			fromCmds, _ := connectionsTo(conn.toCmd, connections)
			if len(toCmds) == 1 && len(fromCmds) == 1 {
				pipeConn := pipeConnection{
					fromCmd: fromCmds[0],
					toCmd:   toCmds[0],
				}
				pipeConnections = append(pipeConnections, pipeConn)
				doneForFromCmd[conn.fromCmd] = true
			}
			// TODO: TEST
			// else {
			// 	doneForFromCmd[fromCmd] = true
			// }
		}
	}
	return pipeConnections
}

type inputOutputConnection struct {
	inputTo         string
	inputScanMethod ScanMethod
	outputFrom      string
}

func findInputOutputConnection(connections []connection) []inputOutputConnection {
	var inputTo, outputFrom string
	var inputScanMethod ScanMethod
	for _, conn := range connections {
		fromCmd := strings.TrimSpace(conn.fromCmd)
		toCmd := strings.TrimSpace(conn.toCmd)
		if len(fromCmd) == 0 && len(toCmd) > 0 {
			inputTo = toCmd
			inputScanMethod = conn.fromScanMethod
		}
		if len(fromCmd) > 0 && len(toCmd) == 0 {
			outputFrom = fromCmd
		}
	}
	if len(inputTo) == 0 && len(outputFrom) == 0 {
		return []inputOutputConnection{}
	}
	return []inputOutputConnection{
		inputOutputConnection{
			inputTo:         inputTo,
			inputScanMethod: inputScanMethod,
			outputFrom:      outputFrom,
		},
	}
}
