package enum

type Gamemode int

const (
	GamemodeSurvival Gamemode = iota
	GamemodeCreative
	GamemodeAdventure
	GamemodeSpectator
)
