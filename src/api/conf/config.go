package conf

import (
	"io/ioutil"

	"github.com/anchormc/anchor/src/api/data"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Host       string    `yaml:"host"`
	Port       uint16    `yaml:"port"`
	LogLevel   int       `yaml:"log_level"`
	OnlineMode bool      `yaml:"online_mode"`
	MaxPlayers int       `yaml:"max_players"`
	MOTD       data.Chat `yaml:"motd"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Host:       "0.0.0.0",
		Port:       25565,
		LogLevel:   1,
		OnlineMode: true,
		MaxPlayers: 20,
		MOTD: data.Chat{
			Text: "A Minecraft Server",
		},
	}
}

func (c *Configuration) ReadFile(file string) error {
	data, err := ioutil.ReadFile(file)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}

func (c Configuration) Validate() error {
	// TODO config file validation

	return nil
}

func (c Configuration) WriteFile(file string) error {
	data, err := yaml.Marshal(c)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, data, 0777)
}
