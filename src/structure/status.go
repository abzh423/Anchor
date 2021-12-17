package structure

type JSONStatusResponse struct {
	Version     StatusVersion `json:"version"`
	Players     StatusPlayers `json:"players"`
	Description string        `json:"description"`
	Favicon     *string       `json:"favicon,omitempty"`
}

type StatusVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type StatusPlayers struct {
	Online int                  `json:"online"`
	Max    int                  `json:"max"`
	Sample []StatusSamplePlayer `json:"sample"`
}

type StatusSamplePlayer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
