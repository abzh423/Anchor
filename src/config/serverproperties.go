package config

import (
	"github.com/anchormc/anchor/src/enum"
)

// ServerProperties is the server.properties file which contains necessary
// key-value pair values for the server.
// TODO finish all standard properties
type ServerProperties struct {
	AllowFlight           bool
	AllowNether           bool
	BroadcastConsoleToOps bool
	BroadcastRconToOps    bool
	Difficulty            enum.Difficulty
}
