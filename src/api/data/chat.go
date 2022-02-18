package data

type Chat struct {
	Text          string `json:"text" yaml:"text"`
	Color         string `json:"color,omitempty" yaml:"color,omitempty"`
	Bold          bool   `json:"bold,omitempty" yaml:"bold,omitempty"`
	Italic        bool   `json:"italic,omitempty" yaml:"italic,omitempty"`
	Underline     bool   `json:"underlined,omitempty" yaml:"underlined,omitempty"`
	Strikethrough bool   `json:"strikethrough,omitempty" yaml:"strikethrough,omitempty"`
	Obfuscated    bool   `json:"obfuscated,omitempty" yaml:"obfuscated,omitempty"`
	Extra         []Chat `json:"extra,omitempty" yaml:"extra,omitempty"`
}
