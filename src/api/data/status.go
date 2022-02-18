package data

type StatusResponse struct {
	Version     StatusVersion `json:"version"`
	Players     StatusPlayers `json:"players"`
	Description Chat          `json:"description"`
	Favicon     string        `json:"favicon,omitempty"`
}

type StatusVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type StatusPlayers struct {
	Online int                  `json:"online"`
	Max    int                  `json:"max"`
	Sample []StatusSamplePlayer `json:"sample,omitempty"`
}

type StatusSamplePlayer struct {
	UUID     string `json:"id"`
	Username string `json:"name"`
}
