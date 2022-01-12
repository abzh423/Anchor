package types

import (
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Host       string `yaml:"host"`
	Port       uint16 `yaml:"port"`
	MOTD       string `yaml:"motd"`
	MaxPlayers int    `yaml:"max_players"`
	Seed       string `yaml:"seed"`
}

func NewConfiguration() *Configuration {
	seed := make([]byte, 16)

	if _, err := rand.Read(seed); err != nil {
		log.Fatal(err)
	}

	return &Configuration{
		Host:       "0.0.0.0",
		Port:       25565,
		MOTD:       "A Minecraft Server",
		MaxPlayers: 20,
		Seed:       hex.EncodeToString(seed),
	}
}

func (c *Configuration) ReadFromFile(path string) error {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}

func (c Configuration) WriteToFile(config string) error {
	data, err := yaml.Marshal(c)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(config, data, 0666)
}
