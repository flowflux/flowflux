package nodecollection

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
