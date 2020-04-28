package nodecollection

import "strings"

type connection struct {
	fromCmd        string
	fromScanMethod ScanMethod
	toCmd          string
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
