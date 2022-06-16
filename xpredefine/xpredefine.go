package xpredefine

var l_outDebugInfo = false
var l_debugMode = false

func Debug() bool {
	return l_debugMode
}

func SetDebug(debugMode bool) {
	l_debugMode = debugMode
}

func OutDebugInfo() bool {
	return l_outDebugInfo
}

func SetOutDebugInfo(outDebugInfo bool) {
	l_outDebugInfo = outDebugInfo
}
