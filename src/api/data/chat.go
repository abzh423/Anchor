package data

type Chat struct {
	Text  string `json:"text"`
	Bold  bool   `json:"bold,omitempty"`
	Extra []Chat `json:"extra,omitempty"`
}
