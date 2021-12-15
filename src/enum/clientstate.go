package enum

type ClientState int

const (
	ClientStateNone ClientState = iota
	ClientStateStatus
	ClientStatePlay
)
