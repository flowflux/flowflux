package config

// EvntEnd ...
const EvntEnd = "#APRX:EVNT:END#"

// EvntEndBytes ...
var EvntEndBytes []byte

func init() {
	EvntEndBytes = []byte(EvntEnd)
}
