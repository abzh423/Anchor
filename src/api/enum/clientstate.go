package enum

type ClientState int

var (
	ClientStateNone   ClientState = 0
	ClientStateStatus ClientState = 1
	ClientStatePlay   ClientState = 2
)
