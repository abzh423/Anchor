package src

type Configuration struct {
	Host       string `yaml:"host"`
	Port       uint16 `yaml:"port"`
	MOTD       string `yaml:"motd"`
	MaxPlayers int32  `yaml:"max_players"`
}

func (c *Configuration) LoadDefaults() {
	c.Host = "127.0.0.1"
	c.Port = 25565
	c.MOTD = "A Minecraft Server"
	c.MaxPlayers = 20
}
