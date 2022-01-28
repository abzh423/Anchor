package data

import "github.com/anchormc/anchor/src/api/protocol"

type StatusResponse struct {
	Version     StatusResponseVersion `json:"version"`
	Players     StatusResponsePlayers `json:"players"`
	Description protocol.Chat         `json:"description"`
	Favicon     *string               `json:"favicon,omitempty"`
}

type StatusResponseVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type StatusResponsePlayers struct {
	Online int                          `json:"online"`
	Max    int                          `json:"max"`
	Sample []StatusResponseSamplePlayer `json:"sample"`
}

type StatusResponseSamplePlayer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
